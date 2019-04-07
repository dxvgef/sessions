# Sessions
使用golang实现的Sessions功能包，仅支持Redis存储，Redis的驱动包使用的是[go-redis/redis](https://github.com/go-redis/redis)，已实现以下功能：
- [x] 写入指定键名的值
- [x] 根据键名读取值并转换类型
- [x] 根据键名删除值
- [x] 清除所有Session数据但保留Cookie中的SessionID
- [x] 清除所有Session数据以及Cookie中的SessionID，下次请求时重新生成新的SessionID
- [x] 默认自动自动更新Session空闲时间，可禁止自动更新

## 使用示例
此示例的HTTP框架使用的是[Tsing](https://github.com/dxvgef/tsing)，也可以和更多框架整合
```Go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dxvgef/tsing"
	"github.com/dxvgef/sessions"
)

// 定义session管理器
var sessManager *sessions.SessionManager

func main() {
	log.SetFlags(log.Lshortfile)

	// 设置session管理器
	err := setSessManager()
	if err != nil {
		log.Fatalln(err.Error())
	}

	app := tsing.New()

	// 定义一个路由处理器用于写入session
	app.Router.GET("/", func(ctx tsing.Context) error {
		// 启用session
		session, err := sessManager.UseSession(ctx.Request, ctx.ResponseWriter)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// 写入值
		err = session.Set("test", "ok")
		return err
	})

	// 定义一个路由处理器用于演示sessions的其它操作
	app.Router.GET("/test", func(ctx tsing.Context) error {
		// 启用session
		session, err := sessManager.UseSession(ctx.Request, ctx.ResponseWriter)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// 读取string类型的值
		var value string
		value, err = session.Get("test").String()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		ctx.ResponseWriter.WriteHeader(200)
		ctx.ResponseWriter.Write([]byte(value))

		// 删除指定key的数据
		err = session.Delete("test")
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// 清除session数据但不删除cookie中的sessionID
		err = session.ClearData()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// 清除session数据以及cookie中的sessionID
		// 下次请求时会重新生成新的sessionID
		err = sessManager.ClearAll(ctx.Request, ctx.ResponseWriter)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		return nil
	})

	// 启动HTTP服务
	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatal(err.Error())
	}
}

// 设置sessionManager
func setSessManager() error {
	var err error
	// 创建session管理器
	sessManager, err = sessions.NewSessions(&sessions.Config{
		CookieName:                 "sessionid",        // cookie中的sessionID名称
		HttpOnly:                   true,               // 仅允许HTTP读取，js无法读取
		Domain:                     "",                 // 作用域名，留空则自动获取当前域名
		Path:                       "/",                // 作用路径
		MaxAge:                     60 * 60,            // 最大生命周期（秒）
		IdleTime:                   20 * time.Minute,   // 空闲超时时间
		Secure:                     false,              // 启用HTTPS
		DisableAutoUpdateIdleTime:  false,              // 禁止自动更新空闲时间
		RedisAddr:                  "127.0.0.1:32771",  // Redis地址
		RedisDB:                    0,                  // Redis数据库
		RedisPassword:              "",                 // Redis密码
		RedisKeyPrefix:             "sess",             // Redis中的键名前缀，必须
		Key:                        "abcdefghijklmnop", // 用于加密sessionID的密钥，密钥的长度16,24,32对应AES-128,AES-192,AES-256算法
	})
	if err != nil {
		return err
	}
	return nil
}
```
