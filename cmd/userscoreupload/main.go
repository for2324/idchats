package main

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.NewPrivateLog("updateuserscore")
	_ = config.Config
	s := NewServer()
	s.Start(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		s.Stop()
		os.Exit(0)
	}()
}
