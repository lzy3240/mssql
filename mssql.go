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

// NewMsssql 初始化连接
func NewMsssql(server, user, password, database string) Mssql {
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
	//打印切片
	//fmt.Println(maps)
	//fmt.Printf("%T\n", maps) //[]map[string]interface {}
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
	//Exec代入参数执行 data[0],data[1]
	//`update config..Version set Issueoperid="zwx" where VID="test"`
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

/*
//main
func main() {
	//创建连接
	ms := newMsssql("localhost", "sa", "111111Aa", "config")
	//defer db.Close()
	defer ms.db.Close()
	//查询
	sqlstr := "select top 3 fid,fieldcode,tradeid from config..fieldinfo where tradeid=016601"
	maps := ms.queryby(ms.db, sqlstr)
	fmt.Println(maps)
	//插入、更新、删除
	sqlstr1 := "update config..Version set Issueoperid=? where VID=?"
	num := ms.modifyby(ms.db, sqlstr1, "zwx", "test")
	fmt.Printf("succeed,%v line affected.\n", num)
}
*/
