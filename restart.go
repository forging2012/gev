package gev

import (
	"os"
	"os/exec"
	"strings"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

func AutoRestart() {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		Log.Println("go run 不启动热更新", os.Args[0])
		return
	}
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
						cmd := exec.Command("/bin/bash", "-c", `ps -ef|grep youyue|grep -v grep|awk '{print "kill -1 "$2|"/bin/bash"}'`)
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
