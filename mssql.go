// Package mssql 基于ODBC的Microsoft SQL server
package mssql

import (
	"database/sql"
	"fmt"

	_ "github.com/alexbrainman/odbc" //odbd连接
	"github.com/demdxx/gocast"
)

// Mssql 结构体
type Mssql struct {
	server   string
	user     string
	password string
	database string
	Db       *sql.DB
}

// NewMssql 初始化连接
func NewMssql(server, user, password, database string) Mssql {
	dsn := fmt.Sprintf("driver={sql server};server=%s;port=1433;uid=%s;pwd=%s;database=%s", server, user, password, database)
	db, err := sql.Open("odbc", dsn) //"driver={sql server};server=s;port=1433;uid=u;pwd=p;database=d"
	checkErr(err)
	//defer db.Close()
	return Mssql{
		server:   server,
		user:     user,
		password: password,
		database: database,
		Db:       db,
	}
}

//Queryby 查询数据操作
func (m Mssql) Queryby(db *sql.DB, sqlstr string) *[]map[string]interface{} {
	rows, err := db.Query(sqlstr)
	checkErr(err)
	defer rows.Close()

	//遍历每一行
	colNames, _ := rows.Columns()
	var cols = make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		cols[i] = new(interface{})
	}
	var maps = make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(cols...)
		checkErr(err)
		var rowMap = make(map[string]interface{})
		for i := 0; i < len(colNames); i++ {
			rowMap[colNames[i]] = convertRow(*(cols[i].(*interface{})))
		}
		maps = append(maps, rowMap)
	}
	//fmt.Println(maps)
	return &maps //返回指针
}

// convertRow 行数据转换
func convertRow(row interface{}) interface{} {
	switch row.(type) {
	case int:
		return gocast.ToInt(row)
	case string:
		return gocast.ToString(row)
	case []byte:
		return gocast.ToString(row)
	case bool:
		return gocast.ToBool(row)
	}
	return row
}

//Modifyby 修改数据操作
func (m Mssql) Modifyby(db *sql.DB, sqlstr string, args ...interface{}) int64 {
	stmt, err := db.Prepare(sqlstr) // Exec、Prepare均可实现增删改
	checkErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	checkErr(err)
	//判断执行结果
	num, err := res.RowsAffected()
	checkErr(err)
	return num
}

//checkErr 检查错误
func checkErr(err error) {
	if err != nil {
		fmt.Println(err) //panic(err)
	}
}
