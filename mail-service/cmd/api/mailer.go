package main

import (
	"bytes"
	"html/template"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

const (
	connectTimeout = time.Second * 10
	sendTimeout    = time.Second * 10
)

type Encryption string

const (
	TLS  Encryption = "tls"
	SSL  Encryption = "ssl"
	None Encryption = "none"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  Encryption
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		// use default address if a specific one is not specified
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		// use default name if a specific one is not specified
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data, // message will be mapped to the templates
	}

	msg.DataMap = data

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// setup SMTP server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption()
	server.KeepAlive = false
	server.ConnectTimeout = connectTimeout
	server.SendTimeout = sendTimeout

	// connect to SMTP server
	smtClient, err := server.Connect()
	if err != nil {
		return err
	}

	// setup new email message
	email := mail.NewMSG()
	email.SetFrom(msg.From)
	email.AddTo(msg.To)
	email.SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	// send email
	err = email.Send(smtClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

// getEncryption returns the encryption type based on the receiver
// encryption field
func (m Mail) getEncryption() mail.Encryption {
	switch m.Encryption {
	case TLS:
		return mail.EncryptionSTARTTLS
	case SSL:
		return mail.EncryptionSSLTLS
	case None:
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
