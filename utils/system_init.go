package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql:", viper.Get("mysql"))
	fmt.Println("config redis:", viper.Get("redis"))
}

func InitMySQL() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")),
		&gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("connect mysql error:", err)
		//return nil
	}
	log.Println("Connect mysql success!")
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
	//return DB
}

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	/*
		pong, err := RDB.Ping().Result()

		if err != nil {
			// 如果发生错误，打印错误信息
			fmt.Println("connect redis error:", err)
		} else {
			// 如果没有错误，打印成功的返回结果
			logs.Println("Connect redis success!", pong)
		}
	*/
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("pub", msg)
	err = RDB.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

// Subscribe 订阅消息到Redis
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := RDB.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("sub", msg.Payload)
	return msg.Payload, err
}
