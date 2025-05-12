package emailer

import (
	"github.com/mailjet/mailjet-apiv3-go"
)

type Mailer struct {
	config      Config
	emailClient *mailjet.Client
}

func NewMailer(conf Config) Mailer {
	return Mailer{
		config:      conf,
		emailClient: mailjet.NewMailjetClient(conf.PublicKey, conf.PrivateKey),
	}
}

func (m Mailer) SendMail(subject string, body string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: m.config.From.Address,
				Name:  m.config.From.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: m.config.To.Address,
					Name:  m.config.To.Name,
				},
			},
			Subject:  subject,
			TextPart: body,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := m.emailClient.SendMailV31(&messages)

	return err
}

func (m Mailer) SendHTMLMail(subject string, body string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: m.config.From.Address,
				Name:  m.config.From.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: m.config.To.Address,
					Name:  m.config.To.Name,
				},
			},
			Subject:  subject,
			HTMLPart: body,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := m.emailClient.SendMailV31(&messages)

	return err
}