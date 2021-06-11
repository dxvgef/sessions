package redis

import (
	"context"
	"errors"
	"time"

	"github.com/dxvgef/sessions"
	"github.com/go-redis/redis/v8"
)

func (rs *Storage) Add(id, key string, value string) (err error) {
	var result bool
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err = rs.redisClient.HSetNX(ctx, rs.config.Prefix+":"+id, key, value).Result()
	if err != nil {
		return
	}
	if !result {
		return errors.New("操作失败，可能是key已存在")
	}
	return
}

func (rs *Storage) Delete(id, key string) (err error) {
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = rs.redisClient.HDel(ctx, rs.config.Prefix+":"+id, key).Result()
	return
}

func (rs *Storage) Put(id, key string, value string) (err error) {
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = rs.redisClient.HSet(ctx, rs.config.Prefix+":"+id, key, value).Result()
	return
}

func (rs *Storage) Update(id, key string, value string) (err error) {
	var result bool
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if result, err = rs.redisClient.HExists(ctx, rs.config.Prefix+":"+id, key).Result(); err != nil {
		return
	}
	if !result {
		return errors.New("nil")
	}
	return rs.redisClient.HSet(ctx, rs.config.Prefix+":"+id, key, value).Err()
}

func (rs *Storage) Get(id, key string) (result sessions.Result) {
	var value string
	err := rs.Connect()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value, err = rs.redisClient.HGet(ctx, rs.config.Prefix+":"+id, key).Result()
	if err != nil {
		if err.Error() == redis.Nil.Error() {
			err = errors.New("nil")
		}
	}
	return sessions.NewResult(value, err)
}

func (rs *Storage) Refresh(id string, expires time.Time) (err error) {
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return rs.redisClient.ExpireAt(ctx, rs.config.Prefix+":"+id, expires).Err()
}

func (rs *Storage) Destroy(id string) (err error) {
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = rs.redisClient.Del(ctx, rs.config.Prefix+":"+id).Result()
	return
}
