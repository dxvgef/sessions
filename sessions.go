package sessions

import (
	"net/http"
	"time"
)

// 会话
type Session struct {
	id     string
	engine *Engine
	req    *http.Request
	resp   http.ResponseWriter
}

// 存储器接口
type Storage interface {
	Add(id, key string, value string) error           // 添加k/v，如果key存在则报错
	Delete(id, key string) error                      // 删除k
	Put(id, key string, value string) error           // 创建或更新
	Update(id, key string, value string) error        // 更新k/v，如果key不存在则报错
	Get(id, key string) (result Result)               // 获取key
	Refresh(id string, expires time.Time) (err error) // 刷新生命周期
	Destroy(id string) (err error)                    // 销毁会话
}

// 在会话中设置一个键值，如果键存在则返回错误
func (sess *Session) Add(key string, value string) (err error) {
	if err = sess.engine.storage.Add(sess.id, key, value); err != nil {
		return
	}
	return sess.Refresh()
}

// 在会话中删除一个键值
func (sess *Session) Delete(key string) (err error) {
	if err = sess.engine.storage.Delete(sess.id, key); err != nil {
		return
	}
	return sess.Refresh()
}

// 在会话中修改一个键值，如果键不存在则返回错误
func (sess *Session) Update(key string, value string) (err error) {
	if err = sess.engine.storage.Update(sess.id, key, value); err != nil {
		return
	}
	return sess.Refresh()
}

// 在会话中设置一个键值，如果键不存在则创建，存在则替换
func (sess *Session) Put(key string, value string) (err error) {
	if err = sess.engine.storage.Put(sess.id, key, value); err != nil {
		return
	}
	return sess.Refresh()
}

// 从会话中读取一个键值
func (sess *Session) Get(key string) (result Result) {
	if result = sess.engine.storage.Get(sess.id, key); result.Err != nil {
		return
	}
	if err := sess.Refresh(); err != nil {
		result.Err = err
	}
	return
}

// 刷新会话，延长生命周期
func (sess *Session) Refresh() (err error) {
	expires := time.Now().Add(time.Duration(int(sess.engine.config.IdleTimeout)) * time.Second)
	if sess.resp.Header().Get("Set-Cookie") == "" {
		http.SetCookie(sess.resp, &http.Cookie{
			Name:     sess.engine.config.Key,
			Value:    sess.id,
			Domain:   sess.engine.config.Domain,
			Path:     sess.engine.config.Path,
			Expires:  expires,
			MaxAge:   int(sess.engine.config.IdleTimeout),
			Secure:   sess.engine.config.Secure,
			HttpOnly: sess.engine.config.HTTPOnly,
		})
	}
	return sess.engine.storage.Refresh(sess.id, expires)
}

// 销毁会话
func (sess *Session) Destroy() (err error) {
	if err = sess.engine.storage.Destroy(sess.id); err != nil {
		return
	}
	http.SetCookie(sess.resp, &http.Cookie{
		Name:    sess.engine.config.Key,
		Value:   "",
		Expires: time.Unix(1, 0),
		MaxAge:  -1,
	})
	return nil
}

// 获取会话的ID
func (sess *Session) GetID() string {
	return sess.id
}

// 获取会话的http.Request
func (sess *Session) GetRequest() *http.Request {
	return sess.req
}

// 获取会话的http.ResponseWriter
func (sess *Session) GetResponseWriter() http.ResponseWriter {
	return sess.resp
}
