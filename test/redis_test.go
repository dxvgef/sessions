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
	tLog     = log.Default()
	err      error
	redisCfg = redis.Config{
		Prefix: "session_test",
	}
	engine     *sessions.Engine
	testServer http.ServeMux
	sessionID  string
)

func TestMain(m *testing.M) {
	tLog.SetFlags(log.Ltime | log.Lshortfile)

	var storage sessions.Storage
	// 创建引擎
	storage, err = redis.New(&redisCfg)
	if err != nil {
		tLog.Println(err)
		return
	}
	engine, err = sessions.New(&sessions.Config{
		IdleTimeout: 20,
	}, storage)
	if err != nil {
		tLog.Println(err)
		return
	}

	// 注册路由
	testServer.HandleFunc("/add", regAdd)
	testServer.HandleFunc("/put", regPut)
	testServer.HandleFunc("/get", regGet)
	testServer.HandleFunc("/update", regUpdate)
	testServer.HandleFunc("/delete", regDelete)
	testServer.HandleFunc("/destroy", regDestroy)

	m.Run()
	os.Exit(0)
}

func TestAll(t *testing.T) {
	t.Run("testAdd", testAdd)
	t.Run("testGet", testGet)
	t.Run("testPut", testPut)
	t.Run("testGet", testGet)
	t.Run("testUpdate", testUpdate)
	t.Run("testGet", testGet)
	t.Run("testDelete", testDelete)
	t.Run("testGet", testGet)
	t.Run("testDestroy", testDestroy)
	t.Run("testGet", testGet)
}

func regAdd(resp http.ResponseWriter, req *http.Request) {
	var sess *sessions.Session
	sess, err = engine.Use(req, resp)
	if err != nil {
		tLog.Println(err)
		return
	}
	if err = sess.Add("username", "dxvgef"); err != nil {
		tLog.Println(err)
		return
	}
	sessionID = sess.GetID()
}
func testAdd(t *testing.T) {
	var (
		req  *http.Request
		resp = httptest.NewRecorder()
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", "/add", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	testServer.ServeHTTP(resp, req)
}

func regPut(resp http.ResponseWriter, req *http.Request) {
	var sess *sessions.Session
	sess, err = engine.Use(req, resp)
	if err != nil {
		tLog.Println(err)
		return
	}
	if err = sess.Put("password", "123456"); err != nil {
		tLog.Println(err)
		return
	}
}
func testPut(t *testing.T) {
	var (
		req  *http.Request
		resp = httptest.NewRecorder()
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", "/put", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	cfg := engine.GetConfig()
	req.AddCookie(&http.Cookie{
		Name:    cfg.Key,
		Value:   sessionID,
		Expires: time.Now().Add(time.Duration(cfg.IdleTimeout) * time.Minute),
		MaxAge:  int(cfg.IdleTimeout),
	})

	testServer.ServeHTTP(resp, req)
}

func regGet(resp http.ResponseWriter, req *http.Request) {
	var (
		sess               *sessions.Session
		username, password string
	)
	sess, err = engine.Use(req, resp)
	if err != nil {
		tLog.Println(err)
		return
	}
	username, err = sess.Get("username")
	if err != nil {
		if err.Error() != "nil" {
			tLog.Println(err)
		}
	}
	tLog.Println("username:", username)
	password, err = sess.Get("password")
	if err != nil {
		if err.Error() != "nil" {
			tLog.Println(err)
		}
	}
	tLog.Println("password:", password)
}

func testGet(t *testing.T) {
	var (
		req  *http.Request
		resp = httptest.NewRecorder()
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", "/get", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	cfg := engine.GetConfig()
	req.AddCookie(&http.Cookie{
		Name:    cfg.Key,
		Value:   sessionID,
		Expires: time.Now().Add(time.Duration(cfg.IdleTimeout) * time.Minute),
		MaxAge:  int(cfg.IdleTimeout),
	})
	testServer.ServeHTTP(resp, req)
}

func regUpdate(resp http.ResponseWriter, req *http.Request) {
	var sess *sessions.Session
	sess, err = engine.Use(req, resp)
	if err != nil {
		tLog.Println(err)
		return
	}
	if err = sess.Update("password", "abcdefg"); err != nil {
		tLog.Println(err)
		return
	}
}

func testUpdate(t *testing.T) {
	var (
		req  *http.Request
		resp = httptest.NewRecorder()
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", "/update", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	cfg := engine.GetConfig()
	req.AddCookie(&http.Cookie{
		Name:    cfg.Key,
		Value:   sessionID,
		Expires: time.Now().Add(time.Duration(cfg.IdleTimeout) * time.Minute),
		MaxAge:  int(cfg.IdleTimeout),
	})

	testServer.ServeHTTP(resp, req)
}

func regDelete(resp http.ResponseWriter, req *http.Request) {
	var sess *sessions.Session
	sess, err = engine.Use(req, resp)
	if err != nil {
		tLog.Println(err)
		return
	}
	if err = sess.Delete("username"); err != nil {
		tLog.Println(err)
		return
	}
}

func testDelete(t *testing.T) {
	var (
		req  *http.Request
		resp = httptest.NewRecorder()
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", "/delete", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	cfg := engine.GetConfig()
	req.AddCookie(&http.Cookie{
		Name:    cfg.Key,
		Value:   sessionID,
		Expires: time.Now().Add(time.Duration(cfg.IdleTimeout) * time.Minute),
		MaxAge:  int(cfg.IdleTimeout),
	})

	testServer.ServeHTTP(resp, req)
}

func regDestroy(resp http.ResponseWriter, req *http.Request) {
	var sess *sessions.Session
	sess, err = engine.Use(req, resp)
	if err != nil {
		tLog.Println(err)
		return
	}
	if err = sess.Destroy(); err != nil {
		tLog.Println(err)
		return
	}
}

func testDestroy(t *testing.T) {
	var (
		req  *http.Request
		resp = httptest.NewRecorder()
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", "/destroy", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	cfg := engine.GetConfig()
	req.AddCookie(&http.Cookie{
		Name:    cfg.Key,
		Value:   sessionID,
		Expires: time.Now().Add(time.Duration(cfg.IdleTimeout) * time.Minute),
		MaxAge:  int(cfg.IdleTimeout),
	})

	testServer.ServeHTTP(resp, req)
}
