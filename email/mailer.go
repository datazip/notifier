package email

func NotifyEmail(subject string, fromEmail string, recipient Recipient, template string, args ...interface{}) {
	standardMailer.NotifyEmail(subject, fromEmail, recipient, template, args...)
}

func NotifyEmailE(subject string, fromEmail string, recipient Recipient, template string, args ...interface{}) error {
	return standardMailer.NotifyEmailE(subject, fromEmail, recipient, template, args...)
}
