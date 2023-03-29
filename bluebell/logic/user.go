package logic

import (
	"github.com/pro911/gin-demo/bluebell/dao/mysql"
	"github.com/pro911/gin-demo/bluebell/pkg/snowflake"
)

func SignUp() {
	//判断用户存不存在
	mysql.QueryUserByUserName()
	//生成uuid
	snowflake.GenID()
	//保存数据
	mysql.InsertUser()
	return
}
