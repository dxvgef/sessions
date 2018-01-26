package sessions

import (
	"net/http"
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

//删除一个键值，如果键名不存在则忽略，不会报错
func (this *sessionObject) Delete(key string) error {
	return redisClient.HDel(this.ID, key).Err()
}

//清除所有redis中的session数据，但不删除cookie中的sessionID
func (this *sessionObject) ClearData() error {
	return redisClient.Del(this.ID).Err()
}
