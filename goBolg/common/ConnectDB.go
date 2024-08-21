package common

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"goBolg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB 连接数据库
func ConnectDB(config *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.Charset,
		config.ParseTime,
		config.Loc,
	)

	fmt.Println("DSN:", dsn) // 打印 DSN 以进行调试

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true, // 禁用自动软删除支持
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Printf("Error getting database connection: %v\n", err)
		return
	}
	err = dbSQL.Close()
	if err != nil {
		fmt.Printf("Error closing database connection: %v\n", err)
	} else {
		fmt.Println("Database connection closed.")
	}
}

// ConnectRedis 连接 Redis
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("Successfully connected to Redis")
	return rdb, nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis(rdb *redis.Client) {
	if err := rdb.Close(); err != nil {
		fmt.Printf("Error closing Redis connection: %v\n", err)
	} else {
		fmt.Println("Redis connection closed.")
	}
}

// ConnectRabbitMQ 连接 RabbitMQ
func ConnectRabbitMQ(cfg *config.RabbitMQConfig) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	fmt.Println("Successfully connected to RabbitMQ")
	return conn, nil
}

// CloseRabbitMQ 关闭 RabbitMQ 连接
func CloseRabbitMQ(conn *amqp.Connection) {
	if err := conn.Close(); err != nil {
		fmt.Printf("Error closing RabbitMQ connection: %v\n", err)
	} else {
		fmt.Println("RabbitMQ connection closed.")
	}
}
