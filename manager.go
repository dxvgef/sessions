package sessions

import (
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// 实例化一个session对象
func NewSessions(config *Config) (*SessionManager, error) {
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
	err := redisClient.Ping().Err()
	if err != nil {
		return nil, err
	}

	// 实例化一个管理器
	var tempManager SessionManager
	// 将redis连接对象传入session管理器
	tempManager.config = config
	manager = &tempManager
	return &tempManager, nil
}

// 使用session，检查sessionID是否存在，如果不存在则创建一个新的并写入到cookie
func (this *SessionManager) UseSession(req *http.Request, resp http.ResponseWriter) (*sessionObject, error) {
	var sessObj sessionObject
	var cookieValid = true
	var sidValue string

	// 从cookie中获得sessionID
	cookieObj, _ := req.Cookie(this.config.CookieName)
	if cookieObj == nil {
		cookieValid = false
	} else if cookieObj.Value == "" {
		cookieValid = false
	}

	// 如果cookie中的sessionID有效
	if cookieValid == true {
		// 将cookie中的值解码
		sid, err := decodeSID(cookieObj.Value, this.config.Key)
		if err != nil {
			return nil, err
		}
		// 将uuid作为sessionID赋值给session对象
		sessObj.ID = sid
	} else {
		var err error
		// 生成一个uuid并赋值给session对象
		sessObj.ID = uuid.New().String()
		// 将uuid结合key加密成sid
		sidValue, err = encodeByBytes(strToByte(this.config.Key), strToByte(sessObj.ID))
		if err != nil {
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

	sessObj.ID = this.config.RedisKeyPrefix + ":" + sessObj.ID
	sessObj.req = req
	sessObj.resp = resp

	return &sessObj, nil
}

// 更新session的空闲时间
func (this *SessionManager) UpdateIdleTime(req *http.Request, resp http.ResponseWriter) error {
	// 从cookie中获得sessionID
	cookieObj, _ := req.Cookie(this.config.CookieName)
	if cookieObj == nil {
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
	sid, err := decodeSID(cookieObj.Value, manager.config.Key)
	if err != nil {
		return err
	}
	// 更新redis的超时时间
	return redisClient.ExpireAt(manager.config.RedisKeyPrefix+":"+sid, time.Now().Add(manager.config.IdleTime)).Err()
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
func (this *SessionManager) ClearAll(req *http.Request, resp http.ResponseWriter) error {
	// 从cookie中获得sessionID
	cookieObj, _ := req.Cookie(this.config.CookieName)
	if cookieObj == nil {
		return nil
	} else if cookieObj.Value == "" {
		return nil
	}
	// 将cookie中的值解码得到sessionID
	sid, err := decodeSID(cookieObj.Value, manager.config.Key)
	if err != nil {
		return err
	}
	// 清除redis中的数据
	err = redisClient.Del(sid).Err()
	if err != nil {
		return err
	}
	// 清除cookie
	http.SetCookie(resp, &http.Cookie{
		Name:     this.config.CookieName,
		Value:    "",
		Domain:   this.config.Domain,
		Path:     this.config.Path,
		Expires:  time.Now().AddDate(-1, 0, 0),
		Secure:   this.config.Secure,
		MaxAge:   -1,
		HttpOnly: this.config.HttpOnly,
	})

	// 更新redis的超时时间
	return nil
}
