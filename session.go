package sessions

import (
	"net/http"
	"strconv"
	"strings"
)

//session对象
type sessionObject struct {
	ID   string
	req  *http.Request
	resp http.ResponseWriter
}

//设置一个键值，如果键名存在则覆盖
func (this *sessionObject) Set(key string, value interface{}) error {
	return redisClient.HSet(this.ID, key, value).Err()
}

func (this *sessionObject) GetString(key string) (string, error) {
	return redisClient.HGet(this.ID, key).Result()
}

func (this *sessionObject) GetMustString(key string, val string) string {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return str
	}
	return val
}

func (this *sessionObject) GetStrings(key string, sep string) ([]string, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return nil, err
	}
	return strings.Split(str, sep), nil
}

func (this *sessionObject) GetMustStrings(key string, sep string, val []string) []string {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	return strings.Split(str, sep)
}

func (this *sessionObject) GetInt8(key string) (int8, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return 0, err
	}
	i, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}
func (this *sessionObject) GetMustInt8(key string, val int8) int8 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	i, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return val
	}
	return int8(i)
}
func (this *sessionObject) GetInt8s(key string, sep string) ([]int8, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return nil, err
	}
	value := strings.Split(str, sep)
	var a []int8
	for k := range value {
		i, err := strconv.ParseInt(value[k], 10, 8)
		if err != nil {
			return nil, err
		}
		a = append(a, int8(i))
	}

	return a, nil
}
func (this *sessionObject) GetMustInt8s(key string, sep string, val []int8) []int8 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	value := strings.Split(str, sep)
	var a []int8
	for k := range value {
		i, err := strconv.ParseInt(value[k], 10, 8)
		if err != nil {
			return val
		}
		a = append(a, int8(i))
	}

	return a
}

func (this *sessionObject) GetUint8(key string) (uint8, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return 0, err
	}
	i, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}
func (this *sessionObject) GetMustUint8(key string, val uint8) uint8 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	i, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		return val
	}
	return uint8(i)
}
func (this *sessionObject) GetUint8s(key string, sep string) ([]uint8, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return nil, err
	}
	value := strings.Split(str, sep)
	var a []uint8
	for k := range value {
		i, err := strconv.ParseUint(value[k], 10, 8)
		if err != nil {
			return nil, err
		}
		a = append(a, uint8(i))
	}

	return a, nil
}
func (this *sessionObject) GetMustUint8s(key string, sep string, val []uint8) []uint8 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	value := strings.Split(str, sep)
	var a []uint8
	for k := range value {
		i, err := strconv.ParseInt(value[k], 10, 8)
		if err != nil {
			return val
		}
		a = append(a, uint8(i))
	}

	return a
}

func (this *sessionObject) GetInt16(key string) (int16, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return 0, err
	}
	i, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}
func (this *sessionObject) GetMustInt16(key string, val int16) int16 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	i, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return val
	}
	return int16(i)
}
func (this *sessionObject) GetInt16s(key string, sep string) ([]int16, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return nil, err
	}
	value := strings.Split(str, sep)
	var a []int16
	for k := range value {
		i, err := strconv.ParseInt(value[k], 10, 16)
		if err != nil {
			return nil, err
		}
		a = append(a, int16(i))
	}

	return a, nil
}
func (this *sessionObject) GetMustInt16s(key string, sep string, val []uint16) []uint16 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	value := strings.Split(str, sep)
	var a []uint16
	for k := range value {
		i, err := strconv.ParseUint(value[k], 10, 16)
		if err != nil {
			return val
		}
		a = append(a, uint16(i))
	}

	return a
}

func (this *sessionObject) GetUint16(key string) (uint16, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return 0, err
	}
	i, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}
func (this *sessionObject) GetMustUint16(key string, val uint16) uint16 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	i, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return val
	}
	return uint16(i)
}
func (this *sessionObject) GetUint16s(key string, sep string) ([]uint16, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return nil, err
	}
	value := strings.Split(str, sep)
	var a []uint16
	for k := range value {
		i, err := strconv.ParseUint(value[k], 10, 16)
		if err != nil {
			return nil, err
		}
		a = append(a, uint16(i))
	}

	return a, nil
}
func (this *sessionObject) GetMustUint16s(key string, sep string, val []uint16) []uint16 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	value := strings.Split(str, sep)
	var a []uint16
	for k := range value {
		i, err := strconv.ParseInt(value[k], 10, 16)
		if err != nil {
			return val
		}
		a = append(a, uint16(i))
	}

	return a
}

func (this *sessionObject) GetInt(key string) (int, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return 0, err
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return i, nil
}
func (this *sessionObject) GetMustInt(key string, val int) int {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return val
	}
	return i
}

func (this *sessionObject) GetInts(key string, sep string) ([]int, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return nil, err
	}
	value := strings.Split(str, sep)
	var a []int
	for k := range value {
		i, err := strconv.Atoi(value[k])
		if err != nil {
			return nil, err
		}
		a = append(a, i)
	}

	return a, nil
}

func (this *sessionObject) GetMustInts(key string, sep string, val []int) []int {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	value := strings.Split(str, sep)
	var a []int
	for k := range value {
		i, err := strconv.Atoi(value[k])
		if err != nil {
			return val
		}
		a = append(a, i)
	}

	return a
}

func (this *sessionObject) GetInt64(key string) (int64, error) {
	return redisClient.HGet(this.ID, key).Int64()
}

func (this *sessionObject) GetMustInt64(key string, val int64) int64 {
	value, err := redisClient.HGet(this.ID, key).Int64()
	if err == nil {
		return value
	}
	return val
}
func (this *sessionObject) GetInt64s(key string, sep string) ([]int64, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return nil, err
	}
	value := strings.Split(str, sep)
	var a []int64
	for k := range value {
		i, err := strconv.ParseInt(value[k], 10, 64)
		if err != nil {
			return nil, err
		}
		a = append(a, i)
	}

	return a, nil
}
func (this *sessionObject) GetMustInt64s(key string, sep string, val []int64) []int64 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	value := strings.Split(str, sep)
	var a []int64
	for k := range value {
		i, err := strconv.ParseInt(value[k], 10, 64)
		if err != nil {
			return val
		}
		a = append(a, i)
	}

	return a
}

func (this *sessionObject) GetUint64(key string) (uint64, error) {
	return redisClient.HGet(this.ID, key).Uint64()
}

func (this *sessionObject) GetMustUint64(key string, val uint64) uint64 {
	value, err := redisClient.HGet(this.ID, key).Uint64()
	if err == nil {
		return value
	}
	return val
}

func (this *sessionObject) GetFloat32(key string) (float32, error) {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return 0, err
	}
	i, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, err
	}
	return float32(i), nil
}

func (this *sessionObject) GetMustFloat32(key string, val float32) float32 {
	str, err := redisClient.HGet(this.ID, key).Result()
	if err == nil {
		return val
	}
	i, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return val
	}
	return float32(i)
}

func (this *sessionObject) GetFloat64(key string) (float64, error) {
	return redisClient.HGet(this.ID, key).Float64()
}

func (this *sessionObject) GetMustFloat64(key string, val float64) float64 {
	value, err := redisClient.HGet(this.ID, key).Float64()
	if err == nil {
		return value
	}
	return val
}

func (this *sessionObject) GetBytes(key string) ([]byte, error) {
	return redisClient.HGet(this.ID, key).Bytes()
}

func (this *sessionObject) GetMustBytes(key string, val []byte) []byte {
	value, err := redisClient.HGet(this.ID, key).Bytes()
	if err == nil {
		return value
	}
	return val
}

func (this *sessionObject) GetBool(key string) (bool, error) {
	s, err := redisClient.HGet(this.ID, "a").Result()
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return b, nil
}
func (this *sessionObject) GetMustBool(key string, val bool) bool {
	s, err := redisClient.HGet(this.ID, "a").Result()
	if err != nil {
		return val
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return val
	}
	return b
}

//删除一个键值，如果键名不存在则忽略，不会报错
func (this *sessionObject) Delete(key string) error {
	return redisClient.HDel(this.ID, key).Err()
}

//清除所有redis中的session数据，但不删除cookie中的sessionID
func (this *sessionObject) ClearData() error {
	return redisClient.Del(this.ID).Err()
}
