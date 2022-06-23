package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func helloworld(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World")
}

func wait(c *gin.Context) {
	sec, err := strconv.Atoi(c.Query("sec"))
	if err != nil {
		sec = 1
	}
	time.Sleep(time.Duration(sec) * time.Second)
	c.String(http.StatusOK, fmt.Sprintf("wait %d sec", sec))
}

func cpu(c *gin.Context) {
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil || count > 100 {
		count = 100
	}
	for i := 0; i < count; i++ {
		bcrypt.GenerateFromPassword([]byte("aaaaaaaaaa"), 10)
	}
	c.String(http.StatusOK, fmt.Sprintf("run bcrypt. count: %d", count))
}

var mem = sync.Map{}

func memory(c *gin.Context) {
	count, err := strconv.Atoi(c.Query("mb"))
	if err != nil {
		count = 1
	}
	count *= 6800
	for i := 0; i < count; i++ {
		id := uuid.New()
		mem.Store(id, id)
	}
	c.String(http.StatusOK, fmt.Sprintf("store uuid. count: %d", count))
}

func clear(c *gin.Context) {
	mem = sync.Map{}
	runtime.GC()
	c.String(http.StatusOK, "clear mem")
}

func route(r *gin.Engine) {
	r.GET("/", helloworld)
	r.GET("/wait", wait)
	r.GET("/cpu", cpu)
	r.GET("/memory", memory)
	r.GET("/memory/clear", clear)
}

func main() {
	port := os.Getenv("PORT")
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	route(r)
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	l.Info("listen", zap.String("port", port))
	l.Fatal(s.ListenAndServe().Error())

}
