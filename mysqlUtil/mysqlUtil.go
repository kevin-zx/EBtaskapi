package mysqlUtil

import (
	"database/sql"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	// db, _ = sql.Open("mysql", "root:Iknowthat@tcp(127.0.0.1:3306)/eb_optimizetest?charset=utf8")
	db, _ = sql.Open("mysql", "remote:Iknowthat@tcp(115.159.79.85:3306)/eb_optimizetest?charset=utf8")
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func Insert(prepareSql string, args ...interface{}) error {
	stmt, err := db.Prepare(prepareSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil

}

func ExecBatch(prepareSql string, args [][]interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(prepareSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, arg := range args {
		_, err = stmt.Exec(arg...)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func Select(prepareSql string, args ...interface{}) ([][]sql.RawBytes, error) {
	rows, err := db.Query(prepareSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var valueArr [][]sql.RawBytes
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var value = make([]sql.RawBytes, len(columns))
		copy(value, values[:])
		valueArr = append(valueArr, value)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return valueArr, nil
}

func SelectAll(sqlstr string, args ...interface{}) (*[]map[string]string, error) {
	stmtOut, err := db.Prepare(sqlstr)
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(args...)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	ret := make([]map[string]string, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
		var value string
		vmap := make(map[string]string, len(scanArgs))
		for i, col := range values {

			if col == nil {
				value = "NULL"
			} else {
				value = string(col)

			}
			vmap[columns[i]] = value
		}

		ret = append(ret, vmap)
	}
	return &ret, nil
}
