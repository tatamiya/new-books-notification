package notifier

import "github.com/slack-go/slack"

type SlackNotifier struct {
	webhookURL string
}

func (s *SlackNotifier) Post(message string) error {
	msg := slack.WebhookMessage{
		Text: message,
	}
	err := slack.PostWebhook(s.webhookURL, &msg)
	return err
}

func NewSlackNotifier(webhookURL string) (*SlackNotifier, error) {

	return &SlackNotifier{
		webhookURL: webhookURL,
	}, nil
}
