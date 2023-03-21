package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func initDB() (err error) {
	dsn := "root:6nv2lxTHDVuQUQN9@tcp(127.0.0.1:3305)/api"

	//也可以使用MustConnect连接不成功就panic
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed,err:%v\n", err)
		return
	}

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)
	return
}

type user struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func queryRawDemo() {
	sqlStr := "select id, name from user where id = ?"
	var u user
	err := DB.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s\n", u.ID, u.Name)
}

func queryMultiRowDemo() {
	sqlStr := "select id,name,age from user where id > ?"
	var users []user
	err := DB.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("querys failed,err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}

func insertRowDemo() {
	sqlStr := "insert into user (name,age) values (?,?)"
	ret, err := DB.Exec(sqlStr, "周义凯", 33)
	if err != nil {
		fmt.Printf("insert failed,err:%v\n", err)
		return
	}

	retID, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastInsertId() failed,err:%v\n", err)
		return
	}
	fmt.Printf("insert success,the id is %d.\n", retID)
}

func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := DB.Exec(sqlStr, 100, 1)
	if err != nil {
		fmt.Printf("update failed,err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() //操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected() failed,err:%v\n", err)
		return
	}
	fmt.Printf("updated success,affected rows:%d\n", n)
}

func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := DB.Exec(sqlStr, 2)
	if err != nil {
		fmt.Printf("delete failed,err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected() //查询影响行数
	if err != nil {
		fmt.Printf("get RowsAffected() failed,err:%v\n", err)
		return
	}
	fmt.Printf("deleted success,affected rows:%d\n", n)
}

// InsertUserDemo 方法用来绑定SQL语句与结构体或map中的同名字段
func InsertUserDemo() (err error) {
	sqlStr := "INSERT INTO `user` (`name`,`age`) VALUES (:name,:age)"
	_, err = DB.NamedExec(sqlStr, map[string]interface{}{
		"name": "七米",
		"age":  30,
	})
	return
}

// namedQuery 与DB.NamedExec同理，这里是支持查询。
func namedQuery() {
	sqlStr := "SELECT * FROM `user` WHERE `name` = :name"
	//使用map做命名查询
	rows, err := DB.NamedQuery(sqlStr, map[string]interface{}{"name": "七米"})
	if err != nil {
		fmt.Printf("db.NamedQuery failed,err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err = rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed,err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}

	u := user{Name: "七米"}
	//使用结构体命名查询,根据结构体字段的db tag进行映射
	rows, err = DB.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("db.NamedQuery failed,err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err = rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed,err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}

func transactionDemo2() (err error) {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Printf("begin trans failed,err:%v\n", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback() //回滚
			panic(p)      //re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() //err is non-nil,don`t change it
		} else {
			err = tx.Commit() //err is  nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()

	sqlStr := "UPDATE `user` set `age`=20 WHERE `id`=?"

	rs, err := tx.Exec(sqlStr, 1)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr failed")
	}
	sqlStr2 := "UPDATE `user` set `age`=50 WHERE id=?"
	rs, err = tx.Exec(sqlStr2, 4)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr2 failed")
	}
	return err
}

type user2 struct {
	Name string `json:"name" db:"name"`
	Age  int    `json:"age" db:"age"`
}

func (u2 user2) value() (driver.Value, error) {
	return []interface{}{u2.Name, u2.Age}, nil
}

// BatchInsertUsers 实现批量插入
func BatchInsertUsers(users2 []*user2) error {
	_, err := DB.NamedExec("INSERT INTO `user` (`name`,`age`) VALUES (:name,:age)", users2)
	return err
}

func main() {
	if err := initDB(); err != nil {
		fmt.Printf("init DB failed, err:%v\n", err)
		return
	}
	fmt.Println("init DB success...")
	//queryRawDemo()
	//insertRowDemo()
	//updateRowDemo()
	//deleteRowDemo()
	//InsertUserDemo()
	//queryMultiRowDemo()
	//namedQuery()
	//transactionDemo2()
	u1 := user2{
		Name: "x1",
		Age:  18,
	}
	u2 := user2{
		Name: "x2",
		Age:  18,
	}
	u3 := user2{
		Name: "x3",
		Age:  18,
	}
	u4 := user2{
		Name: "x4",
		Age:  18,
	}
	users := []*user2{&u1, &u2, &u3, &u4}
	BatchInsertUsers(users)
}
