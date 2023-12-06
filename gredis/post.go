package gredis

import (
	"encoding/json"
	"fmt"

	"ginredis/model"
	"log"
	"time"
)

// 缓存全部文章
func SetCacheAllPosts(data []*model.Post) (err error) {
	strdata, _ := json.Marshal(data)
	err = rdb.Set(ctx, CACHE_POSTS, strdata, 10*time.Second).Err()
	if err != nil {
		log.Println("SET Redis CACHE/ALL-POSTS error: ", err)
		return
	}
	return nil
}

// 获取文章缓存
func GetCacheAllPosts(key string) ([]*model.Post, error) {
	res, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Println("GET Redis post error: ", err)
		return nil, err
	}
	data := []*model.Post{}
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		log.Println("unmarshal failed, err: ", err)
		return nil, err
	}
	return data, nil
}

// 根据文章 id 获取缓存
func GetCachePostById(key string) (model.Post, error) {
	res, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Println("GET Redis Post error, ", err)
		return model.Post{}, err
	}
	data := model.Post{}
	_ = json.Unmarshal([]byte(res), &data)
	return data, nil
}

// 根据文章 id 设置缓存
func SetCachePostById(data *model.Post, postId string) error {
	strdata, _ := json.Marshal(data)
	key := fmt.Sprintf("%s%s", KeyPostIdSet, postId)
	err := rdb.Set(ctx, key, strdata, 10*time.Second).Err()
	if err != nil {
		log.Println("SET Redis Error: ", err)
		return err
	}
	return nil
}

// 更新文章缓存 PUBLISH将消息发布到频道。
func UpdatePost(key string) {
	rdb.Publish(ctx, "post_cache", key)
}

/* 调用方式
currentKey := fmt.Sprintf("%s%s", gredis.KeyPostIdSet, data.PostId)
gredis.UpdatePost(currentKey)
*/

/* 发布方式
rdb.Publish(ctx, "post_cache", key)
*/

/* 订阅方式
sub := rdb.Subscribe(ctx, "post_cache")
ch := sub.Channel()
go func() {
	for msg := range ch {
		...
	}
}()
*/
