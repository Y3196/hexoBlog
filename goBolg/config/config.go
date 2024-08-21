package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Name            string `yaml:"name"`
	Charset         string `yaml:"charset"`
	ParseTime       bool   `yaml:"parseTime"`
	Loc             string `yaml:"loc"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	Port int `yaml:"port"`
}

// RedisConfig Redis 配置结构体
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// RabbitMQConfig RabbitMQ 配置结构体
type RabbitMQConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type WebsiteConfig struct {
	URL string `yaml:"url"`
}

type UploadConfig struct {
	Mode string    `yaml:"mode"`
	Oss  OssConfig `yaml:"oss"`
}

type OssConfig struct {
	URL             string `yaml:"url"`
	Endpoint        string `yaml:"endpoint"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	BucketName      string `yaml:"bucketName"`
}

type SearchConfig struct {
	Mode string `yaml:"mode"`
}

type MailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

// AppConfig 应用程序配置结构体
type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Redis    RedisConfig    `yaml:"redis"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	Website  WebsiteConfig  `yaml:"website"`
	Upload   UploadConfig   `yaml:"upload"`
	Search   SearchConfig   `yaml:"search"`
	Email    MailConfig     `yaml:"email"`
}

// LoadConfig 从 YAML 文件加载配置
func LoadConfig(filename string) (*AppConfig, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	var config AppConfig
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	fmt.Printf("Loaded config: %+v\n", config)

	return &config, nil
}
