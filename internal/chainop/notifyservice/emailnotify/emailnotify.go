package emailnotify

import (
	"Open_IM/internal/chainop/notifyservice"
	"Open_IM/pkg/xlog"
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func NewMailUseCase(host string, port int, fromEmail string, fromEmailPassword string) notifyservice.Sender {
	return &MailUseCase{
		fromUser:          fromEmail,
		host:              host,
		port:              port,
		fromEmail:         fromEmail,
		fromEmailPassword: fromEmailPassword,
	}
}

type MailUseCase struct {
	fromUser                           string
	host, fromEmail, fromEmailPassword string
	port                               int
}

func (m *MailUseCase) Send(receivers []string, subject, content string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)
	d := gomail.NewDialer(m.host, m.port, m.fromEmail, m.fromEmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	sc, err := d.Dial()
	if err != nil {
		xlog.CError(err)
	}

	if err := sc.Send(m.fromUser, receivers, msg); err != nil {
		sc.Close()
		xlog.CError("send to email error", err.Error())
		return err
	}
	xlog.CInfo("send to email ", receivers, subject, content)
	sc.Close()
	return nil
}
