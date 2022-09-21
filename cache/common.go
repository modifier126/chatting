package cache

import (
	"fmt"
	"strconv"

	logging "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
)

// RedisClient Redis缓存客户端单例
var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

//Redis 在中间件中初始化redis链接
func Redis() {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		//Password: conf.RedisPw,  // 无密码，就这样就好了
		DB: int(db),
	})
	_, err := client.Ping().Result()
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	RedisClient = client
}

func Init() {
	file, err := ini.Load("./conf/config.ini")

	if err != nil {
		fmt.Println("redis ini.Load err=", err)
		return
	}

	LoadRedis(file)
	Redis()
}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}
