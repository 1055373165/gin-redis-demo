# 在 gin 中使用 redis 缓存

随着访问量的逐渐提交，使用关系型数据库的站点会遇到一些性能上的瓶颈---元凶：磁盘的 IO 开销。
如今网站的需求主要体现在低延迟的读写速度、支持海量的数据和流量、大规模集群的管理和运营成本的考量；
Redis 是一种内存型数据库，性能高效、支持分布式、可扩展性强、支持丰富的数据结构、以 KV 键值对存储等特征称为了当前很受
欢迎的 NoSQL 数据库之一；

- 应用场景丰富
  - 缓存系统（缓存热点数据：高频读、低频写场景）
  - 消息队列
  - 实时投票系统

在Gin框架中简单将redis作为缓存系统，这种场景读写 redis 一般步骤是
1. 查询 redis 是否存在某个 key，如果存在直接从内存中返回结果
2. 如果 redis 中不存在某个 key 则去 mysql 中查询数据并写进缓存中并设定一个过期时间
3. 当 mysql 数据更新或者删除数据时则将 redis 中对应的 key-value 修改掉。

# 缓存系统场景介绍

一张 post 的数据表，用来存储
- 文章 id
- 标题
- 内容
- 作者 id

等；
写好 model、logic、controller 层的代码能正常运行后，在本地环境运行：

**目录结构**
- controller
  - post.go
  - response.go
  - user.go

go.mod
go.sum
- gredis
  - keys.go
  - post.go
  - redis.go
  - user.go

- logic
  - post.go
  - user.go

main.go
- middleware
  - cache.go

- model
  - post.go
  - user.go

- mysql
  - mysql.go
  - post.go
  - user.go


实现缓存的思路：在 gredis 目录的. keys.go 文件中写常量
在 gin 处理请求时先从 redis 获取数据，如果不存在才放行去查询数据库
然后通过不同的 key 设置到 redis 中

# 创建 redis 客户端

```go
var (
    rdb *redis.Client
    ctx = context.Background()
)

func InitRedis() (err error) {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "", // no password
        DB: 0, // default db
    })

    // connectivity test
    _, err = rdb.Ping(ctx).Result()
    if err != nil {
        return err
    }
    return nil
}
```

# 创建 reids 中的 key 常量（对缓存数据分类定义 key）
在 gredis 目录中创建一个 keys.go 文件，用来存放一些常量
```go
package gredis

const (
    KeyUserIdSet = "CACHE/users/"
    KeyPostIdSet = "CACHE/posts/"
    CACHE_POSTS = "CACHE/all-posts"
    CACHE_USERS = "CACHE/all-users"
)
```

# 判断 redis 中 key 存在
判断 redis 数据库中是否存在 key，用 exists key 来判断 key 是否存在；（对应 redis 的 Exists API）
```go
func ExistKey(key string) bool {
    n, err := rdb.Exists(ctx, key).Result()
    if err != nil {
        log.Println("find exist user key error: ", err)
    }
    if n == 0 {
        log.Println(key, "key not exist")
        return false
    }
    log.Println(key, "key exist")
    return true
}
```

# 设置和获取缓存
## 根据文章 id 获取缓存数据（序列化后的）
```go
func GetCachePostById(key string) (data *model.Post, err error ) {
    res, err := rdb.Get(ctx, key).Result()
    if err != nil {
        log.Println("GET redis post error", err)
        return nil, err
    }
    // 缓存中存储的是序列化后的数据，因此需要进行反序列化为指定的 model
    err = json.Unmarshal([]byte(res), &data)
    return data, nil
}
```
## 根据文章 id 设置缓存
```go
func SetCachePostById(data *model.Post, postid string) (err error) {
    // 序列化
    strdata, _ := json.Marshal(data)
    // 构造数据的唯一 key，组装文章类型缓存的 key 和文章的 id ；KeyPostIdSet = "CACHE/posts/"
    key := fmt.Sprintf("%s%s", KeyPostIdSet, postid)
    err = rdb.Set(ctx, key, strdata, 10*time.Second).Err()
    // 设置 key 10s 的过期时间
    if err != nil {
        log.Println("SET redis ERROR:", err)
        return err
    }
    return nil
}
```

# 设置为 gin 的中间件
我们知道在 gin 框架中可以自定义中间键：
```go
func CacheMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ...
    }
}
```
我们在目录 middleware 创建 cache.go 文件用于判断 redis 中是否存在某个 key；
如果 key 是存在的，那么直接返回 redis 中的数据，并执行 c.Abort()；
否则执行 c.Next() 去数据库查询并将查询结果缓存到 redis；

```go
/*
1. 功能是判断redis中是否存在key,如果存在则取出缓存并返回数据；c.Abort
2. 如果redis中key不存在，则c.Next()继续查询数据库并设置上值
*/
func CacheMiddleware(key string) gin.HandlerFunc {
    return func(c *gin.Context) {
        isExists := gredis.ExistUserKey(key)
        if !isExists {
            c.Next()
        } else {
            switch key {
                case gredis.CACHE_POSTS:
                    data, _ := gredis.GetCacheAllPosts(key)
                    controller.ResponseSuccess(c, data)
                case gredis.CACHE_USERS:
                    data, _ := gredis.GetCacheAllUsers(key)
                    controller.ResponseSuccess(c, data)
            }
            c.Abort()
        }
    }
}
```

# 在路由中加入缓存的中间件
两种路由注册，一种是获取所有的文章 "CACHE/all-posts"
一种是根据文章 id 获取指定的文章 "CACHE/posts/"
```go
r.GET("/getallpost", middleware.CacheMiddleware(gredis.CACHE_POSTS), controller.HandleGetAllPost)
r.GET("/getpostbypostid/:postid", middleware.CacheMiddleware(gredis.KeyPostIdSet), controller.HandleGetPostByPostId)
```
# 对比测试
做两次请求，第一次请求 redis 中不存在 CACHE/all-posts 这个 key，第二次请求将结果已经缓存到 redis 中

# 执行更新和删除操作时
党对数据库中某条数据修改或者删除时，需要删除对应的缓存，或者设置缓存失效时间，或者更新缓存；
对应两种实现
- 更新完数据库后直接删除 redis 中对应的 key 或者为 key 设置过期时间，过期自动触发删除
- 使用发布订阅实现异步缓存失效

第一种方法比较简单直接，但因为不是原子操作所以存在间隔时间窗口，可能导致短暂从缓存获取到旧数据的情况；
第二种方法则使用 redis 的订阅发布功能，配合 go channel 使用，一旦某条数据更新后，将更新信息通过发布订阅机制以 channel 消息的形式推送给 redis，redis 根据接收的消息删除或者更新相关缓存；这种方式以原子操作实现，可以确保缓存和数据库的同步。

## 使用删除或者设置缓存过期时间
```go
func DeleteKey(key string) (err error) {
    if err = rdb.Del(ctx, key).Err(); err != nil {
        log.Println("Delete key Error:", err)
        return err
    }
    return nil
}
```

```go
func SetKeyExpired(key string) (err error) {
    err = rdb.ExpireAt(ctx, key, time.Now().Add(10 * time.Second())).Err()
    if err != nil {
        log.Println("Set Key Expired Error:", err)
        return err
    }
    return nil
}
```

## 使用发布订阅实现异步缓存失效
UpdatePost 函数发布消息
```go
// 当数据库某个数据更新时，调用 redis 的 Publish 函数将这条数据对应 key 的更新消息发布出去
// key 的主题和订阅时设置的主题一致，都是 post_cache，表示订阅有关文章缓存中 key 的更新消息
// 一旦文章缓存中某个文章的数据发生改变，就会通知订阅方更新或者删除缓存数据
// 当订阅成功后会返回一个订阅对象，我们调用这个对象的 Channel 方法可以得到一个用来接收发布消息的 channel，
// 为了避免阻塞主线程，我们可以专门开一个 goroutine 负责去从这个通道中接收订阅的更新消息
// 一旦从通道中读取到消息，那么调用 redis 的 Del 函数删除消息 Payload 中存储的 key 对应的缓存数据
func UpdatePost(key string) {
    rdb.Publish(ctx, "post_cache", key)
}
```
订阅对应的 channel
```go
func SubChannel() {
    sub := rdb.Subscribe(ctx, "post_cache")
    ch := sub.Channel()
    go func() {
        for msg := range ch {
            if err := DeleteKey(msg.Payload); err != nil {
                log.Println("delete key error: ", err)
            }
        }
    }()
}
```

在 controller 中调用函数
```go
func Handlexxx(c *gin.Context) {
    ...
    gredis.UpdatePost(currentKey)
}
```

