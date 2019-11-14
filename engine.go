package sessions

import (
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// Engine session管理引擎
type Engine struct {
	config *Config // 配置
}

// 键不存在时的错误类型
const Nil = RedisError("redis: nil")

// RedisError Redis错误
type RedisError string

// redis连接对象
var redisClient *redis.Client

// Config 配置参数
type Config struct {
	// cookie参数
	CookieName                string        // sessionID的cookie键名
	Domain                    string        // sessionID的cookie作用域名
	Path                      string        // sessionID的cookie作用路径
	Key                       string        // sessionID值加密的密钥
	RedisAddr                 string        // redis地址
	RedisPassword             string        // redis密码
	RedisKeyPrefix            string        // redis键名前缀
	MaxAge                    int           // 最大生命周期（秒）
	IdleTime                  time.Duration // 空闲生命周期
	RedisDB                   int           // redis数据库
	HttpOnly                  bool          // 仅用于http（无法被js读取）
	Secure                    bool          // 启用https
	DisableAutoUpdateIdleTime bool          // 禁止自动更新空闲时间

}

// NewEngine 根据配置实例化一个引擎
func NewEngine(config *Config) (Engine, error) {
	// 实例化一个管理器
	var engine Engine

	// 判断配置是否正确
	if config.CookieName == "" {
		return engine, errors.New("CookieName参数值不正确")
	}
	keyLen := len(config.Key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return engine, errors.New("密钥的长度必须是16、24、32个字节")
	}
	if config.RedisAddr == "" {
		return engine, errors.New("Redis服务器地址参数值不正确")
	}
	if config.RedisDB < 0 {
		return engine, errors.New("Redis数据库参数值不正确")
	}

	// 连接redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		DB:       config.RedisDB,
		Password: config.RedisPassword,
	})
	if err := redisClient.Ping().Err(); err != nil {
		return engine, err
	}

	// 将redis连接对象传入session管理器
	engine.config = config
	return engine, nil
}

// Use 使用session，检查sessionID是否存在，如果不存在则创建一个新的并写入到cookie
func (this *Engine) Use(req *http.Request, resp http.ResponseWriter) (*Session, error) {
	var sess Session
	var cookieValid = true
	var sidValue string

	// 从cookie中获得sessionID
	cookieObj, err := req.Cookie(this.config.CookieName)
	if err != nil || cookieObj == nil {
		cookieValid = false
	} else if cookieObj.Value == "" {
		cookieValid = false
	}

	// 如果cookie中的sessionID有效
	if cookieValid {
		// 将cookie中的值解码
		sid, err := decodeSID(cookieObj.Value, this.config.Key)
		if err != nil {
			return nil, err
		}
		// 将uuid作为sessionID赋值给session对象
		sess.ID = sid
	} else {
		var err error
		// 生成一个uuid并赋值给session对象
		sess.ID = uuid.New().String()
		// 将uuid结合key加密成sid
		if sidValue, err = encodeByBytes(strToByte(this.config.Key), strToByte(sess.ID)); err != nil {
			return nil, err
		}

		// 创建一个cookie对象并赋值后写入到客户端
		http.SetCookie(resp, &http.Cookie{
			Name:     this.config.CookieName,
			Value:    sidValue,
			Domain:   this.config.Domain,
			Path:     this.config.Path,
			Expires:  time.Now().Add(this.config.IdleTime),
			Secure:   this.config.Secure,
			MaxAge:   this.config.MaxAge,
			HttpOnly: this.config.HttpOnly,
		})
	}

	sess.ID = this.config.RedisKeyPrefix + ":" + sess.ID
	sess.req = req
	sess.resp = resp

	// 自动更新空闲时间
	if !this.config.DisableAutoUpdateIdleTime {
		if err := this.UpdateIdleTime(req, resp); err != nil {
			return nil, err
		}
	}

	return &sess, nil
}

// 更新session的空闲时间
func (this *Engine) UpdateIdleTime(req *http.Request, resp http.ResponseWriter) error {
	// 从cookie中获得sessionID
	cookieObj, err := req.Cookie(this.config.CookieName)
	if err != nil || cookieObj == nil {
		return nil
	} else if cookieObj.Value == "" {
		return nil
	}

	// 更新cookie的超时时间
	http.SetCookie(resp, &http.Cookie{
		Name:     this.config.CookieName,
		Value:    cookieObj.Value,
		Domain:   this.config.Domain,
		Path:     this.config.Path,
		Expires:  time.Now().Add(this.config.IdleTime),
		Secure:   this.config.Secure,
		MaxAge:   this.config.MaxAge,
		HttpOnly: this.config.HttpOnly,
	})

	// 将cookie中的值解码
	sid, err := decodeSID(cookieObj.Value, this.config.Key)
	if err != nil {
		return err
	}
	// 更新redis的超时时间
	return redisClient.ExpireAt(this.config.RedisKeyPrefix+":"+sid, time.Now().Add(this.config.IdleTime)).Err()
}

// 解码得到sessionID
func decodeSID(hexStr, key string) (string, error) {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	// 将解码后的sid值解密成uuid
	sid, err := decodeByBytes(strToByte(key), b)
	if err != nil {
		return sid, err
	}

	return sid, nil
}

// 清除当前session的所有redis数据和cookie中的sessionID
func (engine *Engine) ClearAll(req *http.Request, resp http.ResponseWriter) error {
	// 从cookie中获得sessionID
	cookieObj, err := req.Cookie(engine.config.CookieName)
	if err != nil || cookieObj == nil {
		return nil
	} else if cookieObj.Value == "" {
		return nil
	}
	// 将cookie中的值解码得到sessionID
	sid, err := decodeSID(cookieObj.Value, engine.config.Key)
	if err != nil {
		return err
	}
	// 清除redis中的数据
	if err = redisClient.Del(sid).Err(); err != nil {
		return err
	}
	// 清除cookie
	http.SetCookie(resp, &http.Cookie{
		Name:     engine.config.CookieName,
		Value:    "",
		Domain:   engine.config.Domain,
		Path:     engine.config.Path,
		Expires:  time.Now().AddDate(-1, 0, 0),
		Secure:   engine.config.Secure,
		MaxAge:   -1,
		HttpOnly: engine.config.HttpOnly,
	})

	// 更新redis的超时时间
	return nil
}
