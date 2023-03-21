package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	//设置默认值
	viper.SetDefault("fileDir", "./")
	//读取配置文件
	viper.SetConfigName("config") //配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   //如果配置文件的名称中没有扩展名，则需要配置此项
	//viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("/etc/appname/")  //查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.appname") //多次调用以添加多个搜索路径
	viper.AddConfigPath("./config")       //还可以在工作目录中查找配置

	err := viper.ReadInConfig() //查找并读取配置文件
	if err != nil {             //处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file:%s\n", err))
	}

	//实时监控配置文件的变化
	viper.WatchConfig()
	//当配置发生变化之后调用的一个回调函数
	viper.OnConfigChange(func(e fsnotify.Event) {
		//配置文件发送变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})

	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("version"))
	})
	r.Run(":8000")
}
