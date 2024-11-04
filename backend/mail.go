package main

import (
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

var mailServer *mail.SMTPServer

func init() {
	mailServer = mail.NewSMTPClient()

	mailServer.Host = config.Mail.Server
	mailServer.Port = config.Mail.Port
	mailServer.Encryption = mail.EncryptionSSLTLS

	mailServer.Username = config.Mail.User
	mailServer.Password = config.Mail.Password

	mailServer.ConnectTimeout = 10 * time.Second
	mailServer.SendTimeout = 10 * time.Second
}
