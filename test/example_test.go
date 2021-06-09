package test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dxvgef/sessions"
	"github.com/dxvgef/sessions/storage/redis"
)

var (
	err      error
	redisCfg = redis.Config{
		Prefix: "session_test",
	}
	engine     sessions.Engine
	testServer http.ServeMux
)

func TestMain(m *testing.M) {
	var storage sessions.Storage
	// 创建引擎
	storage, err = redis.New(&redisCfg)
	if err != nil {
		log.Println(err)
		return
	}
	engine, err = sessions.New(&sessions.Config{}, storage)
	if err != nil {
		log.Println(err)
		return
	}

	m.Run()
	os.Exit(0)
}

func TestAll(t *testing.T) {
	t.Run("testAdd", testAdd)
}

func testAdd(t *testing.T) {
	var (
		req  *http.Request
		resp = httptest.NewRecorder()
	)
	testServer.HandleFunc("/add", func(resp http.ResponseWriter, req *http.Request) {
		var sess *sessions.Session
		sess, err = engine.Use(req, resp)
		if err != nil {
			resp.WriteHeader(500)
			_, _ = resp.Write([]byte(err.Error()))
			return
		}
		if err = sess.Add("username", "dxvgef"); err != nil {
			resp.WriteHeader(500)
			_, _ = resp.Write([]byte(err.Error()))
			return
		}
		if err = sess.Add("password", "123456"); err != nil {
			resp.WriteHeader(500)
			_, _ = resp.Write([]byte(err.Error()))
			return
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", "/add", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	testServer.ServeHTTP(resp, req)

}
