package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dataSourceNameFormat = "%s:%s@tcp(%s:%s)/%s"
	driverName           = "mysql"
)

// DBConf 创建数据库连接池所需的配置
type DBConf struct {
	Name        string
	Host        string
	Port        string
	UserName    string
	Password    string
	MaxLifeTime int
	MaxOpenConn int
	MaxIdleConn int
}

// DB 对sql.DB进行装饰, 对常用的操作方法进行封装
type DB struct {
	*sql.DB
	stmts   sync.Map
	stmtsmu sync.Mutex
}

// NewDB 返回包装了指定配置创建的DB连接池的DB实例
func NewDB(cf *DBConf) (db *DB, err error) {
	dsn := fmt.Sprintf(dataSourceNameFormat,
		cf.UserName,
		cf.Password,
		cf.Host,
		cf.Port,
		cf.Name,
	)

	odb, err := sql.Open(driverName, dsn)
	if err != nil {
		err = fmt.Errorf("sql open: %w", err)
		return
	}

	if err = odb.Ping(); err != nil {
		err = fmt.Errorf("db ping: %w", err)
		return
	}

	odb.SetConnMaxLifetime(time.Duration(cf.MaxLifeTime) * time.Second)
	odb.SetMaxOpenConns(cf.MaxOpenConn)
	odb.SetMaxIdleConns(cf.MaxIdleConn)

	db = &DB{
		DB:    odb,
		stmts: sync.Map{},
	}

	return
}

// Prepare 缓存预处理语句，避免频繁的预处理调度
func (db *DB) Prepare(qs string) (stmt *sql.Stmt, err error) {
	val, ok := db.stmts.Load(qs)
	if !ok {
		db.stmtsmu.Lock()
		defer db.stmtsmu.Unlock()
		val, ok = db.stmts.Load(qs)
		if !ok {
			stmt, err = db.DB.Prepare(qs)
			if err != nil {
				return
			}
			db.stmts.Store(qs, stmt)
			return
		}
	}
	stmt = val.(*sql.Stmt)
	return
}

// Close 关闭缓存的预处理语句
func (db *DB) Close() (err error) {
	db.stmts.Range(func(key, val interface{}) bool {
		stmt := val.(*sql.Stmt)
		_ = stmt.Close()
		return true
	})
	err = db.DB.Close()
	return
}

// Query 查询多行记录
func (db *DB) Query(qs string, st interface{}, args ...interface{}) (err error) {
	if ok := isstlist(st); !ok {
		err = NewErrInvalidScanTo("non-nil *[]*struct")
		return
	}

	stmt, err := db.Prepare(qs)
	if err != nil {
		err = fmt.Errorf("sql prepare: %w", err)
		return
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		err = fmt.Errorf("sql query: %w", err)
		return
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		err = fmt.Errorf("sql columns: %w", err)
		return
	}

	irt := reflect.TypeOf(st).Elem().Elem().Elem()
	lrv := reflect.ValueOf(st).Elem()
	cmp := getcolmp(irt)

	for rows.Next() {
		ivp := reflect.New(irt)
		sts := getsts(ivp, cols, cmp)
		lrv = reflect.Append(lrv, ivp)
		if err = rows.Scan(sts...); err != nil {
			err = fmt.Errorf("sql scan: %w", err)
			return
		}
	}

	reflect.ValueOf(st).Elem().Set(lrv)

	return
}

// QueryRow 查询单行
func (db *DB) QueryRow(qs string, st interface{}, args ...interface{}) (err error) {
	if ok := isstrecord(st); !ok {
		return NewErrInvalidScanTo("non-nil *struct")
	}

	stmt, err := db.Prepare(qs)
	if err != nil {
		err = fmt.Errorf("sql prepare: %w", err)
		return
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		err = fmt.Errorf("sql query: %w", err)
		return
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("sql columns: %w", err)
	}

	rt := reflect.TypeOf(st).Elem()
	colmp := getcolmp(rt)
	rp := reflect.ValueOf(st)
	sts := getsts(rp, cols, colmp)

	if !rows.Next() {
		err = fmt.Errorf("sql scan: %w", sql.ErrNoRows)
		return
	}

	if err = rows.Scan(sts...); err != nil {
		err = fmt.Errorf("sql scan: %w", err)
		return
	}

	return
}

// QueryRowAndScan 查询单行并将值填充到对应变量上
func (db *DB) QueryRowAndScan(qs string, args []interface{}, st ...interface{}) (err error) {
	stmt, err := db.Prepare(qs)
	if err != nil {
		err = fmt.Errorf("sql prepare: %w", err)
		return
	}

	if err = stmt.QueryRow(args...).Scan(st...); err != nil {
		err = fmt.Errorf("sql query and scan: %w", err)
		return
	}

	return
}

// Exec 执行sql语句
func (db *DB) Exec(qs string, args ...interface{}) (affected, lastID int64, err error) {
	stmt, err := db.Prepare(qs)
	if err != nil {
		err = fmt.Errorf("sql prepare: %w", err)
		return
	}

	rst, err := stmt.Exec(args...)
	if err != nil {
		err = fmt.Errorf("sql exec: %w", err)
		return
	}

	affected, err = rst.RowsAffected()
	if err != nil {
		err = fmt.Errorf("sql rows affected: %w", err)
		return
	}

	lastID, err = rst.LastInsertId()
	if err != nil {
		err = fmt.Errorf("sql last id: %w", err)
		return
	}

	return
}
