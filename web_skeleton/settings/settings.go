package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("config")               //指定配置文件名称 (不需要带后缀)
	viper.SetConfigType("yaml")                 //指定配置文件后缀
	viper.AddConfigPath("./")                   //指定查找配置文件的路径 (这里使用相对路径)
	if err = viper.ReadInConfig(); err != nil { //读取配置信息
		//读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return
	}
	viper.WatchConfig() //监听配置文件变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件修改...\n")
	})
	return
}
