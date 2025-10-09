package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"iFall/internal/config"
	"iFall/pkg/errs"
	"net"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendMessage(ctx context.Context, sub string, content []byte, to []string, attachFiles []string) error
}

type emailSender struct {
	Auth        smtp.Auth
	EmailConfig config.EmailConfig
}

func NewEmailSender(a smtp.Auth, cfg config.EmailConfig) EmailSender {
	return &emailSender{
		Auth:        a,
		EmailConfig: cfg,
	}
}

func (es *emailSender) SendMessage(ctx context.Context, sub string, content []byte, to []string, attachFiles []string) error {
	op := "emailSender.SendMessage"
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", es.EmailConfig.Name, es.EmailConfig.Address)
	e.Subject = sub
	e.HTML = content
	e.To = to
	for _, f := range attachFiles {
		if _, err := e.AttachFile(f); err != nil {
			return errs.NewAppError(op, err)
		}
	}
	dialer := &net.Dialer{}

	conn, err := dialer.DialContext(ctx, "tcp", es.EmailConfig.SmtpServerAddress)
	if err != nil {
		return errs.NewAppError(op, err)
	}
	defer conn.Close()
	tlsConfig := &tls.Config{
		ServerName: es.EmailConfig.SmtpAddress,
	}

	tlsConn := tls.Client(conn, tlsConfig)
	client, err := smtp.NewClient(tlsConn, es.EmailConfig.SmtpAddress)
	if err != nil {
		return errs.NewAppError(op, err)
	}
	if err := client.Auth(es.Auth); err != nil {
		return errs.NewAppError(op, err)
	}
	if err := e.SendWithTLS(es.EmailConfig.SmtpServerAddress, es.Auth, tlsConfig); err != nil {
		return errs.NewAppError(op, err)
	}
	return nil
}
