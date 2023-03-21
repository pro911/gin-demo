package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //隐式调用 init()
	"time"
)

var db *sql.DB

func initMySQL() (err error) {
	//dsn:data source name
	dsn := "root:6nv2lxTHDVuQUQN9@tcp(127.0.0.1:3305)/api"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	//做完错误检查之后,确保db不为nil
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to db failed, err:%v\n", err)
		return
	}

	db.SetConnMaxLifetime(time.Second * 10) //连接池 链接的生命周期,在连接在池中存活时间超过 60 秒后，连接会自动关闭并从池中移除。
	db.SetMaxOpenConns(200)                 //最大连接数
	db.SetConnMaxIdleTime(10)               //最大空闲连接数 当业务处于休息状态是 会有指定数量的空闲连接数在等待
	return
}

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func queryRowDemo() {
	sqlStr := "select id,name from user where id=?"
	var u user
	if err := db.QueryRow(sqlStr, 1).Scan(&u.ID, &u.Name); err != nil {
		fmt.Printf("scan failed,err:%v\n", err)
		return
	}
	fmt.Printf("id:%d and name:%v\n", u.ID, u.Name)
}

func queryMultiRowDemo() {
	sqlStr := "select id,name from user where id < ?"
	rows, err := db.Query(sqlStr, 1)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u user

		if err = rows.Scan(&u.ID, &u.Name); err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d and name:%s\n", u.ID, u.Name)
	}
}

func prepareQueryDemo() {
	sqlStr := "select id,name from user where id < ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed,err:%v\n", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(3)
	if err != nil {
		fmt.Printf("query failed,err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err = rows.Scan(&u.ID, &u.Name)
		if err != nil {
			fmt.Printf("scan failed,err:%v\n", err)
			return
		}
		fmt.Printf("id:%d and name:%s\n", u.ID, u.Name)
	}
}

func transactionDemo() {
	tx, err := db.Begin() //开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() //回滚
		}
		fmt.Printf("begin trans failed,err:%v\n", err)
		return
	}
	sqlStr := "update user set name='悬崖之上1' where id = ?"
	ret1, err := tx.Exec(sqlStr, 1)
	if err != nil {
		tx.Rollback() //回滚
		fmt.Printf("exec sql1 failed,err:%v\n", err)
		return
	}

	affRow1, err := ret1.RowsAffected() //返回影响的行数 int
	if err != nil {
		tx.Rollback() //回滚
		fmt.Printf("exec ret1.RowsAffected() failed,err:%v\n", err)
		return
	}

	sqlStr2 := "update user set name='悬崖之上2' where id = ?"
	ret2, err := db.Exec(sqlStr2, 2)
	if err != nil {
		tx.Rollback() //回滚
		fmt.Printf("exec sql2 failed,err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() //回滚
		fmt.Printf("exec ret2.RowsAffected() failed,err:%v\n", err)
		return
	}
	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交...")
		tx.Commit()
	} else {
		tx.Rollback()
		fmt.Println("事务回滚...")
	}
	fmt.Println("exec trans success!")
}
func main() {

	if err := initMySQL(); err != nil {
		fmt.Printf("connect to db failed,err:%v\n", err)
	}
	defer db.Close()
	fmt.Println("connect to db success")

	queryRowDemo()
	//queryMultiRowDemo()
	//prepareQueryDemo()
	//transactionDemo()
	fmt.Println("查询结束了...")
}
