package redis

import (
	"context"
	"errors"
	"time"

	"github.com/dxvgef/sessions"
	"github.com/redis/go-redis/v9"
)

// Add 添加键值到指定的session id
func (stg *Storage) Add(id, key string, value string) (err error) {
	var result bool
	if err = stg.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err = stg.redisClient.HSetNX(ctx, stg.config.Prefix+":"+id, key, value).Result()
	if err != nil {
		return
	}
	if !result {
		return errors.New("操作失败，可能是key已存在")
	}
	return
}

// Delete 从session中删除指定的键
func (stg *Storage) Delete(id, key string) (err error) {
	if err = stg.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = stg.redisClient.HDel(ctx, stg.config.Prefix+":"+id, key).Result()
	return
}

// Put 写入session中的值
func (stg *Storage) Put(id, key string, value string) (err error) {
	if err = stg.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = stg.redisClient.HSet(ctx, stg.config.Prefix+":"+id, key, value).Result()
	return
}

// Update 更新session中的值
func (stg *Storage) Update(id, key string, value string) (err error) {
	var result bool
	if err = stg.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if result, err = stg.redisClient.HExists(ctx, stg.config.Prefix+":"+id, key).Result(); err != nil {
		return
	}
	if !result {
		return errors.New("nil")
	}
	return stg.redisClient.HSet(ctx, stg.config.Prefix+":"+id, key, value).Err()
}

// Get 获取session中的值
func (stg *Storage) Get(id, key string) (result sessions.Result) {
	var value string
	err := stg.Connect()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value, err = stg.redisClient.HGet(ctx, stg.config.Prefix+":"+id, key).Result()
	if err != nil {
		if err.Error() == redis.Nil.Error() {
			err = errors.New("nil")
		}
	}
	return sessions.NewResult(value, err)
}

// Refresh 刷新session生命周期
func (stg *Storage) Refresh(id string, expires time.Time) (err error) {
	if err = stg.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return stg.redisClient.ExpireAt(ctx, stg.config.Prefix+":"+id, expires).Err()
}

// Destroy 销毁Session
func (stg *Storage) Destroy(id string) (err error) {
	if err = stg.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = stg.redisClient.Del(ctx, stg.config.Prefix+":"+id).Result()
	return
}
