package main

import (
	"context"
	"fmt"
	"github.com/pro911/gin-demo/web_skeleton2/dao/mysql"
	"github.com/pro911/gin-demo/web_skeleton2/dao/redis"
	"github.com/pro911/gin-demo/web_skeleton2/logger"
	"github.com/pro911/gin-demo/web_skeleton2/routers"
	"github.com/pro911/gin-demo/web_skeleton2/settings"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		return
	}
	defer zap.L().Sync() //当程序退出是执行同步
	zap.L().Debug("logger init success...")

	//3.初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init MySQL failed,err:%v\n", err)
		return
	}
	defer mysql.Close() //当程序停止时关闭mysql资源

	//4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init Redis failed,err:%v\n", err)
		return
	}
	defer redis.Close()

	//5.注册路由
	r := routers.Setup()

	//6.启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.HttpServer.Port),
		Handler: r,
	}
	fmt.Printf("httpServer.addr:%s\n", fmt.Sprintf(":%d", settings.Conf.HttpServer.Port))

	go func() {
		//开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen", zap.Error(err))
		}
	}()

	//等待中断信号来优雅地关闭服务器,为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) //创建一个接收信号的通道
	//kill 默认会发送 syscall.SIGTERM 信号
	//kill -2 发送 syscall.SIGINT 信号,我们常用的Ctrl+C就是触发系统SIGINT信号
	//kill -9 发送 syscall.SIGKILL 信号,但是不能被捕获,所以不需要添加它
	//signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不会阻塞
	<-quit                                               //阻塞在此,当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	//创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//5秒内优雅关闭服务(将未处理完的请求处理完再关闭服务),超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
