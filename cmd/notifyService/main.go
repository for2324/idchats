package main

import (
	"Open_IM/internal/chainop/notifyservice"
	"Open_IM/internal/chainop/notifyservice/emailnotify"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/xlog"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/** 本服务是为了邮件发送服务， 内容
 */
func main() {
	_ = config.Config
	xlog.InitLevel(xlog.DefaultOption())
	xlog.CInfo("config.Config.Kafka.BusinessTop.Topic:", config.Config.Kafka.BusinessTop.Topic)
	xlog.CInfo("config.Config.EmailSend.EmailSmtpHost:", config.Config.EmailSend.EmailSmtpHost)
	server := NewServer(config.Config.Kafka.BusinessTop.Group,
		config.Config.Kafka.BusinessTop.Addr,
		"business", config.Config.Kafka.BusinessTop.Topic,
		&notifyservice.ListenerSendServer{
			SenderServiceEmail: emailnotify.NewMailUseCase(
				config.Config.EmailSend.EmailSmtpHost,
				config.Config.EmailSend.EmailSmtpPort,
				config.Config.EmailSend.FromAddress,
				config.Config.EmailSend.FromPassword),
			SenderServiceSms: nil,
		}, nil, nil)
	if server != nil {
		//GO-223.8617.58
		go server.Start(context.Background())
	}

	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 在Init()中的另一goroutine中打开server
EXIT:
	for {
		sig := <-sc
		xlog.CInfof("接收到信号[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if server != nil {
				server.Stop(context.Background())
			}

			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	time.Sleep(time.Second)
	os.Exit(state)
	return
}
