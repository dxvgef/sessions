package sessions

import (
	"time"

	"github.com/go-redis/redis"
)

// session管理器
type SessionManager struct {
	config Config // 配置
}

// 键不存在时的错误类型
const Nil = RedisError("redis: nil")

type RedisError string

// redis连接对象
var redisClient *redis.Client

// 配置参数
type Config struct {
	// cookie参数
	CookieName                string        // sessionID的cookie键名
	HttpOnly                  bool          // 仅用于http（无法被js读取）
	Domain                    string        // sessionID的cookie作用域名
	Path                      string        // sessionid的cookie作用路径
	MaxAge                    int           // 最大生命周期（秒）
	IdleTime                  time.Duration // 空闲生命周期
	Secure                    bool          // 启用https
	Key                       string        // sessionID值加密的密钥
	DisableAutoUpdateIdleTime bool          // 禁止自动更新空闲时间
	// redis参数
	RedisAddr      string // redis地址
	RedisDB        int    // redis数据库
	RedisPassword  string // redis密码
	RedisKeyPrefix string // redis键名前缀
}
