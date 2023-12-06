package gredis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func InitRedis() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	// 消息订阅
	SubChannel()
	return
}

func SubChannel() {
	// 为客户端订阅指定的频道。可以省略频道来创建空订阅。(没有提供任何 channels 参数)
	// 请注意，此方法不会等待 Redis 的响应，因此订阅可能不会立即激活。
	// 要强制连接等待，您可以在返回的 *Pub Sub 上调用 Receive() 方法
	// sub := client.Subscribe(queryResp)
	// iface, err := sub.Receive()
	// if err != nil {
	// 	 handle error
	// }
	sub := rdb.Subscribe(ctx, "post_cache")
	ch := sub.Channel()
	go func() {
		for msg := range ch {
			if err := DeleteKey(msg.Payload); err != nil {
				log.Println("delete key error: ", err)
				return
			}
		}
	}()
}

/*
	type Message struct {
	    Channel      string
	    Pattern      string
	    Payload      string
	    PayloadSlice []string
	}
*/
func ExistKey(key string) bool {
	n, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		log.Println("find exists user key error", err)
		return false
	}
	if n == 0 {
		log.Println(key, "key not exist")
		return false
	}
	log.Println(key, "key exist")
	return true
}

func SetKeyExpired(key string) error {
	err := rdb.ExpireAt(ctx, key, time.Now().Add(-10*time.Second)).Err()
	if err != nil {
		log.Println("Set Key Expired Error:", err)
		return err
	}
	return nil
}

func DeleteKey(key string) error {
	if err := rdb.Del(ctx, key).Err(); err != nil {
		log.Println("Delete key error: ", err)
		return err
	}
	return nil
}
