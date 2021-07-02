// Author: Steve Zhang
// Date: 2020/9/16 5:16 下午

package mysql

func (ct *DBContainer) Query(qs string, to interface{}, args ...interface{}) (err error) {
	db := ct.MustGetDB()
	defer ct.PutDB(db)

	err = db.Query(qs, to, args...)

	return
}

func (ct *DBContainer) QueryRow(qs string, to interface{}, args ...interface{}) (err error) {
	db := ct.MustGetDB()
	defer ct.PutDB(db)

	err = db.QueryRow(qs, to, args...)
	return
}

func (ct *DBContainer) QueryRowAndScan(qs string, args []interface{}, to ...interface{}) (err error) {
	db := ct.MustGetDB()
	defer ct.PutDB(db)

	err = db.QueryRowAndScan(qs, args, to...)

	return
}

func (ct *DBContainer) Exec(qs string, args ...interface{}) (affected, lastID int64, err error) {
	db := ct.MustGetDB()
	defer ct.PutDB(db)

	affected, lastID, err = db.Exec(qs, args...)

	return
}
