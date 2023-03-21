package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pro911/gin-demo/web_skeleton2/settings"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	//dsn := "root:6nv2lxTHDVuQUQN9@tcp(127.0.0.1:3306)/api"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	//也可以使用MustConnect连接不成功就panic Must带这个关键字的一般都是有panic的
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

// Close 将关闭数据库并阻止启动新查询。关闭，然后等待服务器上已开始处理的所有查询完成。
//
//	 关闭数据库的情况很少见，
//		因为数据库句柄是长期存在的，并且在许多 goroutines 之间共享
func Close() {
	_ = db.Close()
}
