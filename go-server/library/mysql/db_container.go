package mysql

import (
	"errors"
	"fmt"
	"go-sever/library/conf"
	"time"
)

type DBContainer struct {
	*conf.Container
}

var ErrGetDBConfFuncIsNil = errors.New("get db conf func is nil")

type GetDBConfFunc func() (*DBConf, error)

func NewDBContainer(getDBConf GetDBConfFunc) (ct *DBContainer, err error) {
	if getDBConf == nil {
		err = ErrGetDBConfFuncIsNil
		return
	}

	getObjConf := func() (icf conf.IConf, err error) {
		icf, err = getDBConf()
		if err != nil {
			err = fmt.Errorf("get db conf: %w", err)
			return
		}
		return
	}

	ict, err := conf.NewContainer(getObjConf, compareDBConf, newDBObj, resetDBObj)
	if err != nil {
		err = fmt.Errorf("new conf container: %w", err)
		return
	}

	ct = &DBContainer{
		Container: ict,
	}

	return
}

func compareDBConf(iocf, incf conf.IConf) (rst conf.CompareObjConfRst, err error) {
	ocf, ok := iocf.(*DBConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	ncf, ok := incf.(*DBConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	switch {
	case
		ncf.UserName != ocf.UserName,
		ncf.Password != ocf.Password,
		ncf.Host != ocf.Host,
		ncf.Port != ocf.Port,
		ncf.Name != ocf.Name:
		rst = conf.CompareObjConfRstNeedReplace
		return

	case
		ncf.MaxLifeTime != ocf.MaxLifeTime,
		ncf.MaxIdleConn != ocf.MaxIdleConn,
		ncf.MaxOpenConn != ocf.MaxOpenConn:
		rst = conf.CompareObjConfRstNeedReset
		return

	default:
		rst = conf.CompareObjConfRstNoNeed
		return
	}
}

func newDBObj(icf conf.IConf) (iobj conf.IObject, err error) {
	cf, ok := icf.(*DBConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	iobj, err = NewDB(cf)
	if err != nil {
		err = fmt.Errorf("new db: %w", err)
		return
	}

	return
}

func resetDBObj(iobj conf.IObject, iocf, incf conf.IConf) (err error) {
	ocf, ok := iocf.(*DBConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	ncf, ok := incf.(*DBConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	db, ok := iobj.(*DB)
	if !ok {
		err = conf.ErrInvalidObjectType
		return
	}

	switch {
	case ncf.MaxLifeTime != ocf.MaxLifeTime:
		db.SetConnMaxLifetime(time.Duration(ncf.MaxLifeTime) * time.Second)
		ocf.MaxLifeTime = ncf.MaxLifeTime
		fallthrough

	case ncf.MaxIdleConn != ocf.MaxIdleConn:
		db.SetMaxIdleConns(ncf.MaxIdleConn)
		ocf.MaxIdleConn = ncf.MaxIdleConn
		fallthrough

	case ncf.MaxOpenConn != ocf.MaxOpenConn:
		ocf.MaxOpenConn = ncf.MaxIdleConn
		db.SetMaxOpenConns(ncf.MaxOpenConn)
		fallthrough

	default:
	}

	return
}

func (ct *DBContainer) MustGetDB() (db *DB) {
	obj := ct.MustGetObj()

	db, ok := obj.(*DB)
	if !ok {
		panic(conf.ErrInvalidObjectType)
	}

	return
}

func (ct *DBContainer) PutDB(db *DB) {
	ct.PutObj(db)
}
