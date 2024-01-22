package main

import (
	cronTask "Open_IM/internal/cron_task"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("start cronTask")
	//cronTask.CheckChainBBTPledgePoolVolume()
	//开启定时器：
	cronTask.CheckPledgePoolVolume()
	InitSignal(func() {})
}

type SignalQuitFunc func()

func InitSignal(fn SignalQuitFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		select {
		case sig := <-ch:
			{
				switch sig {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
					if fn != nil {
						fn()
					}
					close(ch)
					return
				case syscall.SIGHUP:
					close(ch)
				default:
					close(ch)
				}
			}
		}
	}
}
