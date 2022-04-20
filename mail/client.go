package mail

import (
	"fmt"
	"log"

	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/quavious/blog-factory-server/config"
)

type MailClient struct {
	*mailjet.Client
	sender string
}

func NewMailClient(config *config.Config) *MailClient {
	apiKey, secretKey := config.GetEmailKey()
	newClient := mailjet.NewMailjetClient(apiKey, secretKey)
	return &MailClient{
		Client: newClient,
		sender: config.GetEmailAddress(),
	}
}

func (client *MailClient) SendToken(emailToken string, receiver string) bool {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: client.sender,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: receiver,
				},
			},
			Subject:  "Email Verification",
			TextPart: fmt.Sprintf("The email verification token is %s. Input this code in 10 minutes.", emailToken),
		},
	}
	messages := mailjet.MessagesV31{
		Info: messagesInfo,
	}
	_, err := client.SendMailV31(&messages)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
