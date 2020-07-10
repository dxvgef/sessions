package sessions

import (
	"encoding/hex"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis/v7"
)

// Engine session管理引擎
type Engine struct {
	config     *Config         // 配置
	seedIDNode *snowflake.Node // 种子ID节点
}

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
func NewEngine(config *Config) (*Engine, error) {
	// 实例化一个管理器
	var engine Engine

	// 判断配置是否正确
	if config.CookieName == "" {
		return nil, errors.New("CookieName参数值不正确")
	}
	keyLen := len(config.Key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, errors.New("密钥的长度必须是16、24、32个字节")
	}
	if config.RedisAddr == "" {
		return nil, errors.New("Redis服务器地址参数值不正确")
	}
	if config.RedisDB < 0 {
		return nil, errors.New("Redis数据库参数值不正确")
	}

	// 连接redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		DB:       config.RedisDB,
		Password: config.RedisPassword,
	})
	if err := redisClient.Ping().Err(); err != nil {
		return nil, err
	}

	// 将redis连接对象传入session管理器
	engine.config = config

	// 创建种子ID节点的实例
	seedIDNode, err := newSeedIDNode()
	if err != nil {
		return nil, err
	}
	engine.seedIDNode = seedIDNode

	return &engine, nil
}

// Use 使用session，检查sessionID是否存在，如果不存在则创建一个新的并写入到cookie
func (this *Engine) Use(req *http.Request, resp http.ResponseWriter) (*Session, error) {
	var (
		sess        Session
		cookieValid = true
	)

	// 从cookie中获得sessionID
	cookieObj, err := req.Cookie(this.config.CookieName)
	if err != nil || cookieObj == nil || cookieObj.Value == "" {
		cookieValid = false
	}

	// 如果cookie中的sessionID有效
	if cookieValid {
		// 将cookie中的cid解码成sid
		sid, err := decodeSID(cookieObj.Value, this.config.Key)
		if err != nil {
			return nil, err
		}
		sess.CookieID = cookieObj.Value
		sess.StorageID = sid
	} else {
		// 如果cookies中的sessionID无效
		// 生成种子id
		seedID := this.seedIDNode.Generate().String() + strconv.FormatUint(uint64(rand.New(rand.NewSource(rand.Int63n(time.Now().UnixNano()))).Uint32()), 10)
		// 用种子ID编码成cid
		cid, err := encodeByBytes(strToByte(this.config.Key), strToByte(seedID))
		if err != nil {
			return nil, err
		}
		sess.CookieID = cid
		sess.StorageID = this.config.RedisKeyPrefix + ":" + seedID
		// 创建一个cookie对象并赋值后写入到客户端
		http.SetCookie(resp, &http.Cookie{
			Name:     this.config.CookieName,
			Value:    cid,
			Domain:   this.config.Domain,
			Path:     this.config.Path,
			Expires:  time.Now().Add(this.config.IdleTime),
			Secure:   this.config.Secure,
			MaxAge:   this.config.MaxAge,
			HttpOnly: this.config.HttpOnly,
		})
	}

	sess.resp = resp
	sess.engine = this

	return &sess, nil
}

// 更新session的空闲时间
func (this *Engine) UpdateIdleTime(cid, sid string, resp http.ResponseWriter) error {
	// 更新cookie的超时时间
	http.SetCookie(resp, &http.Cookie{
		Name:     this.config.CookieName,
		Value:    cid,
		Domain:   this.config.Domain,
		Path:     this.config.Path,
		Expires:  time.Now().Add(this.config.IdleTime),
		Secure:   this.config.Secure,
		MaxAge:   this.config.MaxAge,
		HttpOnly: this.config.HttpOnly,
	})

	return redisClient.ExpireAt(sid, time.Now().Add(this.config.IdleTime)).Err()
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

// 清除指定请求的所有会话数据，包含redis数据和cookie中的sessionID
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

// VerityRequest 校验request中的session id是否有效
func (this *Engine) VerityRequest(req *http.Request) (bool, error) {
	// 从cookie中获得sessionID
	cookieObj, err := req.Cookie(this.config.CookieName)
	if err != nil || cookieObj == nil {
		return false, nil
	} else if cookieObj.Value == "" {
		return false, nil
	}
	// 校验session id
	return this.VerityID(cookieObj.Value)
}

// VerityID 校验session id是否有效
func (this *Engine) VerityID(id string) (bool, error) {
	// 将id解码
	sid, err := decodeSID(id, this.config.Key)
	if err != nil {
		return false, err
	}
	if redisClient.Exists(sid).Val() == 0 {
		return false, nil
	}
	return true, nil
}

// 获取指定ID的服务端数据
func (engine *Engine) GetByID(id, key string) *Value {
	var (
		result Value
		value  string
	)
	result.Key = key
	// 将id值解码得到sessionID
	sid, err := decodeSID(id, engine.config.Key)
	if err != nil {
		result.Error = err
		return &result
	}
	value, err = redisClient.HGet(sid, key).Result()
	if err != nil {
		result.Error = err
		return &result
	}
	result.Value = value
	return &result
}

// 设置指定ID的服务端数据，如果键名存在则覆盖其值
func (engine *Engine) SetByID(id, key string, value interface{}) error {
	// 将id值解码得到sessionID
	sid, err := decodeSID(id, engine.config.Key)
	if err != nil {
		return err
	}
	return redisClient.HSet(sid, key, value).Err()
}

// ClearByID 清除指定ID的所有redis中的session数据，不删除cookie中的数据，但能从服务端使会话失效
func (engine *Engine) ClearByID(id string) error {
	// 将id值解码得到sessionID
	sid, err := decodeSID(id, engine.config.Key)
	if err != nil {
		return err
	}
	// 清除redis中的数据
	if err = redisClient.Del(sid).Err(); err != nil {
		return err
	}
	return nil
}

// DeleteByID 删除指定ID的服务端，如果键名不存在则忽略，不会报错
func (engine *Engine) DeleteByID(id, key string) error {
	// 将id值解码得到sessionID
	sid, err := decodeSID(id, engine.config.Key)
	if err != nil {
		return err
	}
	return redisClient.HDel(sid, key).Err()
}

// 设置种子ID的实例
func newSeedIDNode() (*snowflake.Node, error) {
	snowflake.Epoch = time.Now().Unix()
	rand.Seed(rand.Int63n(time.Now().UnixNano()))
	node := 0 + rand.Int63n(1023-0)
	return snowflake.NewNode(node)
}
