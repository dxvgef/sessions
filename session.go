package sessions

import (
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v7"
)

// session对象
type Session struct {
	StorageID string // 存储器中的ID
	CookieID  string // cookies中的ID
	resp      http.ResponseWriter
	engine    *Engine
}

// 键不存在时的错误类型
const Nil = redis.Nil

// session值
type Value struct {
	Key   string
	Value string
	Error error
	*redis.StringCmd
}

// Get 读取参数值
func (obj *Session) Get(key string) *Value {
	var result Value
	result.Key = key

	value, err := redisClient.HGet(obj.StorageID, key).Result()
	if err != nil {
		result.Error = err
		return &result
	}
	result.Value = value

	// 自动更新空闲时间
	if !obj.engine.config.DisableAutoUpdateIdleTime {
		if err := obj.engine.UpdateIdleTime(obj.CookieID, obj.StorageID, obj.resp); err != nil {
			result.Error = err
			return &result
		}
	}

	return &result
}

// 设置一个键值，如果键名存在则覆盖
func (obj *Session) Set(key string, value interface{}) error {
	err := redisClient.HSet(obj.StorageID, key, value).Err()
	if err != nil {
		return err
	}
	// 自动更新空闲时间
	if !obj.engine.config.DisableAutoUpdateIdleTime {
		if err := obj.engine.UpdateIdleTime(obj.CookieID, obj.StorageID, obj.resp); err != nil {
			return err
		}
	}

	return nil
}

// String 将值转为string类型
func (v *Value) String() (string, error) {
	if v.Error != nil {
		return "", v.Error
	}
	return v.Value, nil
}

// Int 将值转为int类型，如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (v *Value) Int(def ...int) (int, error) {
	defLen := len(def)
	if v.Error != nil {
		if defLen == 0 {
			return 0, v.Error
		}
		return def[0], nil
	}
	value, err := strconv.Atoi(v.Value)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return value, nil
}

// Int32 将参数值转为int32类型，如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (v *Value) Int32(def ...int32) (int32, error) {
	defLen := len(def)
	if v.Error != nil {
		if defLen == 0 {
			return 0, v.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseInt(v.Value, 10, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return int32(value), nil
}

// Int64 将参数值转为int64类型，如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (v *Value) Int64(def ...int64) (int64, error) {
	defLen := len(def)
	if v.Error != nil {
		if defLen == 0 {
			return 0, v.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseInt(v.Value, 10, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return value, nil
}

// Uint32 将参数值转为uint32类型，如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (v *Value) Uint32(def ...uint32) (uint32, error) {
	defLen := len(def)
	if v.Error != nil {
		if defLen == 0 {
			return 0, v.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseUint(v.Value, 10, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return uint32(value), nil
}

// Uint64 将参数值转为uint64类型，如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (v *Value) Uint64(def ...uint64) (uint64, error) {
	defLen := len(def)
	if v.Error != nil {
		if defLen == 0 {
			return 0, v.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseUint(v.Value, 10, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return value, nil
}

// Float32 将参数值转为float32类型，如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (v *Value) Float32(def ...float32) (float32, error) {
	defLen := len(def)
	if v.Error != nil {
		if defLen == 0 {
			return 0, v.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseFloat(v.Value, 32)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return float32(value), nil
}

// Float64 将参数值转为float64类型，如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (v *Value) Float64(def ...float64) (float64, error) {
	defLen := len(def)
	if v.Error != nil {
		if defLen == 0 {
			return 0, v.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return 0, err
	}
	return value, nil
}

// Bool 将参数值转为bool类型，如果传入了def参数值，在转换出错时返回def，并且第二个出参永远为nil
func (v *Value) Bool(def ...bool) (bool, error) {
	defLen := len(def)
	if v.Error != nil {
		if defLen == 0 {
			return false, v.Error
		}
		return def[0], nil
	}
	value, err := strconv.ParseBool(v.Value)
	if err != nil {
		if defLen > 0 {
			return def[0], nil
		}
		return false, err
	}
	return value, nil
}

// Delete 删除一个键值，如果键名不存在则忽略，不会报错
func (this *Session) Delete(key string) error {
	return redisClient.HDel(this.StorageID, key).Err()
}

// ClearData 清除所有redis中的session数据，但不删除cookie中的sessionID
func (this *Session) ClearData() error {
	return redisClient.Del(this.StorageID).Err()
}
