package redis

import (
	"context"
	"fmt"
	"github.com/pro911/gin-demo/bluebell/settings"
	"github.com/redis/go-redis/v9"
	"time"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// Init 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DBName,
		PoolSize: cfg.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err = rdb.Ping(ctx).Result()
	return
}

// Close 关闭客户端，释放所有打开的资源。
// 关闭客户端很少见，因为客户端是长期存在的，
// 并且在许多goroutine之间共享
func Close() {
	_ = rdb.Close()
}
