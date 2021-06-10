# Sessions
Go语言的Sessions库。使用`Storage`接口来实现数据存储层的抽象化，可通过实现该接口来将数据存储到Memcached、Redis、LevelDB、RocksDB、MySQL等数据库。还支持值类型转换的链式方法，有助于减少代码量。

### 已内置的数据存储器：
- [x] Redis（支持v6）

### 与`Tsing`框架结合使用的示例

```Go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dxvgef/sessions"
	"github.com/dxvgef/sessions/storage/redis"
	"github.com/dxvgef/tsing"
)

// 定义session引擎
var engine *sessions.Engine

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	// 创建session引擎
	if err := newSessEngine(); err != nil {
		log.Fatalln(err.Error())
	}

	app := tsing.New(&tsing.Config{})
	
	// 定义一个路由处理器用于写入session
	app.Router.GET("/", func(ctx *tsing.Context) error {
		// 启用session
		session, err := engine.Use(ctx.Request, ctx.ResponseWriter)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// 写入值
		err = session.Put("test", "ok")
		return err
	})

	// 定义一个路由处理器用于演示sessions的其它操作
	app.Router.GET("/test", func(ctx *tsing.Context) error {
		// 启用session
		session, err := engine.Use(ctx.Request, ctx.ResponseWriter)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// 读取string类型的值
		var result sessions.Result
		result = session.Get("test")
		if result.Err() != nil {
			log.Println(result.Err())
			return result.Err()
		}
		ctx.ResponseWriter.WriteHeader(200)
		ctx.ResponseWriter.Write(result.Bytes())

		return nil
	})

	// 启动HTTP服务
	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatal(err.Error())
	}
}

// 创建session引擎
func newEngine() error {
	var (
		err     error
		storage sessions.Storage
	)
	// 创建存储器
	storage, err = redis.New(&redis.Config{
		Addr: "127.0.0.1:6379", // redis server的地址
		Prefix: "sess",         // redis的key前缀
		Username: "",           // redis 6支持的用户名
		Password: "",           // redis的密码
		DB: 0,                  // redis的库
	})
	if err != nil {
		return err
	}
	// 创建session引擎
	engine, err = sessions.New(&sessions.Config{
		Key:         "sessionid", // cookie中的Session ID的键名，默认为"sessionid"
		HTTPOnly:    false,       // 仅允许HTTP读取，JS无法读取
		Domain:      "",          // 作用域名，默认为空
		Path:        "",          // 作用路径，默认为空
		IdleTimeout: 20 * 60,     // 空闲超时时间(秒)
		Secure:      false,       // 仅在HTTPS协议中有效
	}, storage)
	if err != nil {
		return err
	}
	return nil
}
```
更多示例见`/test/`目录中的单元测试文件