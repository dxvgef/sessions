package sessions

import (
	"errors"
	"net/http"
	"path"
	"time"

	"github.com/google/uuid"
)

// Config 管理器配置
type Config struct {
	Key         string        // Session ID的Cookie键名
	GenerateID  func() string // 生成 Session ID的值
	Domain      string        // Session ID的cookie作用域名
	Path        string        // Session ID的cookie作用路径
	IdleTimeout uint          // 空闲超时（秒）
	HTTPOnly    bool          // 仅用于HTTP传输（无法被JS脚本读取）
	Secure      bool          // 启用安全，使Session仅在HTTPS下才会有效
}

// Engine 引擎
type Engine struct {
	config  *Config // 管理器配置
	storage Storage // 存储器实例
}

// Storage 存储器接口
type Storage interface {
	Add(id, key string, value string) error           // 添加k/v，如果key存在则报错
	Delete(id, key string) error                      // 删除k
	Put(id, key string, value string) error           // 创建或更新
	Update(id, key string, value string) error        // 更新k/v，如果key不存在则报错
	Get(id, key string) (result Result)               // 获取key
	Refresh(id string, expires time.Time) (err error) // 刷新生命周期
	Destroy(id string) (err error)                    // 销毁会话
}

// New 创建新的引擎
func New(config *Config, storage Storage) (engine *Engine, err error) {
	if config == nil {
		err = errors.New("必须定义Session的配置")
		return
	}
	if storage == nil {
		err = errors.New("必须定义存储器")
		return
	}
	if config.GenerateID != nil && config.GenerateID() == "" {
		err = errors.New("ID生成回调函数返回值不能为空")
		return
	}
	if config.Key == "" {
		config.Key = "sessionid"
	}
	if config.Path != "" && string(path.Dir(config.Path)[0]) != "/" {
		err = errors.New("作用路径不正确")
		return
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = 20 * 60
	}
	engine = &Engine{
		config:  config,
		storage: storage,
	}
	return
}

// GetConfig 获取引擎的配置
func (engine *Engine) GetConfig() Config {
	return *engine.config
}

// Use 使用会话
func (engine *Engine) Use(req *http.Request, resp http.ResponseWriter) (*Session, error) {
	var (
		err    error
		hasKey bool
		ck     *http.Cookie
		id     uuid.UUID
	)
	if req == nil || resp == nil {
		return nil, errors.New("req和resp参数不为是空指针")
	}

	// 从cookie中读取session id
	ck, err = req.Cookie(engine.config.Key)
	if err == nil &&
		ck != nil &&
		ck.Value != "" {
		hasKey = true
	}

	// 创建一个session实例
	var sess Session
	sess.engine = engine
	sess.req = req
	sess.resp = resp
	// 获取或生成新的session id
	if hasKey {
		sess.id = ck.Value
		return &sess, nil
	}
	if engine.config.GenerateID != nil {
		// 执行回调函数生成session id
		sess.id = engine.config.GenerateID()
	} else {
		// 使用UUID V4生成session id
		id, err = uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		sess.id = id.String()
	}
	return &sess, nil
}
