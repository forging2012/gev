package gev

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	fsnotify "gopkg.in/fsnotify.v1"
)

var Log = log.New(os.Stdout, "[ gev ] ", log.Ltime|log.Lshortfile)

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func AutoRestart() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Log.Println("启动热更新失败", err)
		return
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					time.AfterFunc(1*time.Second, func() {
						cmd := exec.Command("kill", "-1", strconv.Itoa(syscall.Getpid()))
						err := cmd.Start()
						Log.Println("重启中...", err)
					})
				}
			}
		}
	}()

	err = watcher.Add(os.Args[0])
	if err == nil {
		Log.Println("启动热更新成功")
	} else {
		Log.Println("启动热更新失败", err)
	}
}

func stack() []byte {
	buf := make([]byte, 10240)
	n := runtime.Stack(buf, false)
	if n > 710 {
		copy(buf, buf[710:n])
		return buf[:n-710]
	}
	return buf[:n]
}

func CrossDomainMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "x-auth-token,x-device,x-uuid,content-type")
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(200)
			}
		}
	}
}
