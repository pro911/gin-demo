package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// Init 初始化连接
func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.dbname"),
		PoolSize: viper.GetInt("redis.pool_size"),
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
