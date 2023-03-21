package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "6nv2lxTHDVuQUQN9", //password
		DB:       12,                 //use default DB
		PoolSize: 100,                //连接池大小
	})

	//设置一个带超时时间的空上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() //
	_, err = rdb.Ping(ctx).Result()
	return err
}

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init redis client failed,err:%v\n", err)
		return
	}
	fmt.Println("connect redis success...")
	defer rdb.Close()

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	err := rdb.Set(ctx, "pro911", "pro911@qq.com", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "pro911").Result()
	if err != nil {
		fmt.Printf("get redis key: pro911 failed,err:%v\n", err)
		return
	}
	fmt.Println(val)

	v, err := rdb.Del(ctx, "pro911").Result()
	if err != nil {
		return
	}
	fmt.Println(v)
	val2, err := rdb.Get(ctx, "pro911").Result()
	if err == redis.Nil {
		fmt.Printf("get redisx key:pro911 failed,err:%v\n", err)
	}
	if err != nil {
		fmt.Printf("get redisx1 key:pro911 failed,err:%v\n", err)
	}
	fmt.Println(val2)
}
