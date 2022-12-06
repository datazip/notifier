package email

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Recipient struct {
	ToEmails  []string
	CCEmails  []string
	BCCEmails []string
}

func getHTMLTemplate(template string, args ...interface{}) (string, error) {
	htmlData, err := os.ReadFile(template)
	if err != nil {
		return "", err
	}

	data := string(htmlData)

	for i, val := range args {
		data = strings.ReplaceAll(data, fmt.Sprintf("$%d", i+1), fmt.Sprintf("%v", val))
	}

	return data, nil
}

func (m *Mailer) newSession() (*session.Session, error) {
	// create new AWS session
	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(m.config.AccessKey, m.config.SecretKey, ""),
		Region:      aws.String(m.config.Region)},
	)

}

func (m *Mailer) NotifyEmail(subject string, fromEmail string, recipient Recipient, template string, args ...interface{}) {
	if err := standardMailer.NotifyEmailE(subject, fromEmail, recipient, template, args...); err != nil {
		logrus.Error("failed to send email : %s", err)
	}
}

// NotifyEmailE sends email to specified email IDs
func (m *Mailer) NotifyEmailE(subject string, fromEmail string, recipient Recipient, template string, args ...interface{}) error {
	msgBody, err := getHTMLTemplate(template, args...)
	if err != nil {
		return fmt.Errorf("failed to read html template: %s : %s", template, err)
	}

	// set to section
	var recipients []*string
	for _, r := range recipient.ToEmails {
		recipient := r
		recipients = append(recipients, &recipient)
	}

	// set cc section
	var ccRecipients []*string
	if len(recipient.CCEmails) > 0 {
		for _, r := range recipient.CCEmails {
			ccrecipient := r
			ccRecipients = append(ccRecipients, &ccrecipient)
		}
	}

	// set bcc section
	var bccRecipients []*string
	if len(recipient.BCCEmails) > 0 {
		for _, r := range recipient.BCCEmails {
			bccrecipient := r
			recipients = append(recipients, &bccrecipient)
		}
	}

	sess, err := m.newSession()
	if err != nil {
		return err
	}

	// create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{

		// Set destination emails
		Destination: &ses.Destination{
			CcAddresses:  ccRecipients,
			ToAddresses:  recipients,
			BccAddresses: bccRecipients,
		},

		// Set email message and subject
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(msgBody),
				},
			},

			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},

		// send from email
		Source: aws.String(fromEmail),
	}

	// Call AWS send email function which internally calls to SES API
	_, err = svc.SendEmail(input)
	if err != nil {
		return err
	}

	logrus.Println("Email sent successfully to: ", recipient.ToEmails)
	return nil
}

// Experimental Feature
//
// NotifyRawEmailE sends email to specified email IDs with attachments
func (m *Mailer) NotifyRawEmailE(subject string, fromEmail string, recipient Recipient, attachments []string, template string, args ...interface{}) error {
	// create new AWS session
	sess, err := m.newSession()
	if err != nil {
		return err
	}

	msgBody, err := getHTMLTemplate(template, args...)
	if err != nil {
		return fmt.Errorf("failed to read html template: %s : %s", template, err)
	}

	// create raw message
	msg := gomail.NewMessage()

	// set to section
	var recipients []*string
	for _, r := range recipient.ToEmails {
		recipient := r
		recipients = append(recipients, &recipient)
	}

	// Set to emails
	msg.SetHeader("To", recipient.ToEmails...)

	// cc mails mentioned
	if len(recipient.CCEmails) != 0 {
		// Need to add cc mail IDs also in recipient list
		for _, r := range recipient.CCEmails {
			recipient := r
			recipients = append(recipients, &recipient)
		}
		msg.SetHeader("cc", recipient.CCEmails...)
	}

	// bcc mails mentioned
	if len(recipient.BCCEmails) != 0 {
		// Need to add bcc mail IDs also in recipient list
		for _, r := range recipient.BCCEmails {
			recipient := r
			recipients = append(recipients, &recipient)
		}
		msg.SetHeader("bcc", recipient.BCCEmails...)
	}

	// create an SES session.
	svc := ses.New(sess)

	msg.SetAddressHeader("From", fromEmail, "<name>")
	msg.SetHeader("To", recipient.ToEmails...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", msgBody)

	// If attachments exists
	if len(attachments) != 0 {
		for _, f := range attachments {
			msg.Attach(f)
		}
	}

	// create a new buffer to add raw data
	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	// create new raw message
	message := ses.RawMessage{Data: emailRaw.Bytes()}

	input := &ses.SendRawEmailInput{Source: &fromEmail, Destinations: recipients, RawMessage: &message}

	// send raw email
	_, err = svc.SendRawEmail(input)
	if err != nil {
		return err
	}

	logrus.Println("Email sent successfully to: ", recipient.ToEmails)
	return err
}
