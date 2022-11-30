package email

import "github.com/sirupsen/logrus"

func NotifyEmail(subject string, fromEmail string, recipient Recipient, template string, args ...interface{}) {
	if err := standardMailer.NotifyEmailE(subject, fromEmail, recipient, template, args...); err != nil {
		logrus.Error("failed to send email : %s", err)
	}
}

func NotifyEmailE(subject string, fromEmail string, recipient Recipient, template string, args ...interface{}) error {
	return standardMailer.NotifyEmailE(subject, fromEmail, recipient, template, args...)
}
