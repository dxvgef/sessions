# 使用Golang实现的Sessions功能包，仅支持Redis存储

##使用示例
```Go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dxvgef/httpdispatcher"
	"github.com/dxvgef/sessions"
)

//定义session管理器
var sessManager *sessions.SessionManager

func main() {
	log.SetFlags(log.Lshortfile)

	//设置session管理器
	err := setSessManager()
	if err != nil {
		log.Fatalln(err.Error())
	}

	//定义一个http调度器
	dispatcher := httpdispatcher.New()
	//启用500错误
	dispatcher.EventConfig.ServerError = true
	//定义事件处理器
	dispatcher.Handler.Event = func(e *httpdispatcher.Event) {
		log.Println(e.Source)
		log.Println(e.Message)
	}

	//定义一个路由处理器用于写入session
	dispatcher.Router.GET("/", func(ctx *httpdispatcher.Content) error {
		//启用session
		session, err := sessManager.UseSession(ctx.Request, ctx.ResponseWriter)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		//写入值
		err = session.Set("test", "ok")
		return err
	})

	//定义一个路由处理器用于演示sessions的其它操作
	dispatcher.Router.GET("/test", func(ctx *httpdispatcher.Content) error {
		//启用session
		session, err := sessManager.UseSession(ctx.Request, ctx.ResponseWriter)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		//读取string类型的值
		var value string
		value, err = session.Get("test").String()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		ctx.ResponseWriter.WriteHeader(200)
		ctx.ResponseWriter.Write([]byte(value))

		//删除指定key的数据
		err = session.Delete("test")
		if err != nil {
			log.Println(err.Error())
			return err
		}
		//清除session数据但不删除cookie中的sessionID
		err = session.ClearData()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		//清除session数据以及cookie中的sessionID
		//下次请求时会重新生成新的sessionID
		err = sessManager.ClearAll(ctx.Request, ctx.ResponseWriter)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		return nil
	})

	//启动HTTP服务
	if err := http.ListenAndServe(":8080", dispatcher); err != nil {
		log.Fatal(err.Error())
	}
}

//设置sessionManager
func setSessManager() error {
	var err error
	//创建session管理器
	sessManager, err = sessions.NewSessions(&sessions.Config{
		CookieName:     "sessionid",        //cookie中的sessionID名称
		HttpOnly:       true,               //仅允许HTTP读取，js无法读取
		Path:           "/",                //作用路径
		MaxAge:         60 * 60,            //最大生命周期（秒）
		IdleTime:       20 * time.Minute,   //空闲超时时间
		Secure:         false,              //启用HTTPS
		RedisAddr:      "127.0.0.1:32771",  //Redis地址
		RedisDB:        0,                  //Redis数据库
		RedisPassword:  "",                 //Redis密码
		RedisKeyPrefix: "sess",             //Redis中的键名前缀，必须
		Key:            "abcdefghijklmnop", //用于加密sessionID的密钥，密钥的长度16,24,32对应AES-128,AES-192,AES-256算法
	})
	if err != nil {
		return err
	}
	return nil
}
```
