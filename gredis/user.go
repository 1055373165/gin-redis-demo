package gredis

import (
	"encoding/json"
	"ginredis/model"
	"log"
	"time"
)

// 定义一个可以根据用户 id 获取缓存 key 的方法
func GetCacheKey(id string) string {
	return KeyUserIdSet + ":" + id
}

// 缓存全部的用户
func SetCacheAllUsers(data []*model.User) error {
	strdata, _ := json.Marshal(data)
	err := rdb.Set(ctx, CACHE_USERS, strdata, 5*time.Second).Err()
	if err != nil {
		log.Println("redis set error", err)
		return err
	}
	return nil
}

// GET 获取全部用户缓存
func GetCacheAllUsers(key string) ([]*model.User, error) {
	res, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Println("redis get key ", key, "failed, err: ", err)
		return nil, err
	}
	// res 是序列化后的数据
	data := []*model.User{}
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

// Set 单个用户信息的缓存
func SetCacheUserById(key string, data *model.User) error {
	strdata, _ := json.Marshal(data)
	err := rdb.Set(ctx, key, strdata, 10*time.Minute).Err()
	if err != nil {
		log.Println("set single user info cache failed err: ", err)
		return err
	}
	return nil
}

func GetCacheUserById(key string) (model.User, error) {
	res, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Println("redis client get key", key, "failed, err: ", err)
		return model.User{}, err
	}
	data := model.User{}
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		return model.User{}, err
	}
	return data, nil
}

func DelCacheUserById(key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
