package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-redis/redis/v8"
)

// 存储器本置
type Config struct {
	TLS      *tls.Config // TLS配置
	Network  string      // TCP或者UNIX，默认TCP
	Addr     string      // Redis地址
	Username string      // 账号
	Password string      // 密码
	DB       uint8       // 数据库
	Prefix   string      // 键名前缀
}

// Redis存储器结构
type Storage struct {
	config      *Config       // 存储器配置
	redisClient *redis.Client // Redis客户端实例
}

// 创建存储器
func New(config *Config) (*Storage, error) {
	if config.Addr == "" {
		config.Addr = "127.0.0.1:6379"
	}
	return &Storage{
		config: config,
	}, nil
}

// 连接Redis
func (rs *Storage) Connect() (err error) {
	if rs.redisClient != nil {
		return
	}
	rs.redisClient = redis.NewClient(&redis.Options{
		TLSConfig: rs.config.TLS,
		Network:   rs.config.Network,
		Addr:      rs.config.Addr,
		Username:  rs.config.Username,
		DB:        int(rs.config.DB),
		Password:  rs.config.Password,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return rs.redisClient.Ping(ctx).Err()
}
