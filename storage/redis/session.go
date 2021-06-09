package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

func (rs *Storage) Add(id, key string, value string) (err error) {
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := rs.redisClient.HSetNX(ctx, rs.config.Prefix+":"+id, key, value).Result()
	if err != nil {
		return
	}
	if !result {
		return errors.New("操作失败，可能是key已存在")
	}
	return
}

func (rs *Storage) Delete(id, key string) (err error) {
	var count int64
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err = rs.redisClient.HDel(ctx, rs.config.Prefix+":"+id, key).Result()
	if err != nil {
		return
	}
	if count < 1 {
		return errors.New("delete操作失败")
	}
	return
}

func (rs *Storage) Put(id, key string, value string) (err error) {
	var count int64
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err = rs.redisClient.HSet(ctx, rs.config.Prefix+":"+id, key, value).Result()
	if err != nil {
		return
	}
	if count < 1 {
		return errors.New("put操作失败")
	}
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
	err = rs.redisClient.HSet(ctx, rs.config.Prefix+":"+id, key, value).Err()
	if err != nil {
		return
	}
	return
}

func (rs *Storage) Get(id, key string) (value string, err error) {
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value, err = rs.redisClient.HGet(ctx, rs.config.Prefix+":"+id, key).Result()
	if err != nil {
		if err.Error() == redis.Nil.Error() {
			err = errors.New("nil")
			return
		}
	}
	return
}

func (rs *Storage) Refresh(id string, expires time.Time) (err error) {
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = rs.redisClient.ExpireAt(ctx, rs.config.Prefix+":"+id, expires).Err()
	if err != nil {
		return
	}
	return
}

func (rs *Storage) Destroy(id string) (err error) {
	var count int64
	if err = rs.Connect(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err = rs.redisClient.Del(ctx, rs.config.Prefix+":"+id).Result()
	if err != nil {
		return
	}
	if count < 1 {
		return errors.New("destroy操作失败")
	}
	return
}
