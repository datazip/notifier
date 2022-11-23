package notifier

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// Notify error logs on Slack
func (n *Notifier) NotifyError(errorAt, description, errString string, fields ...string) {
	if err := n.NotifyErrorE(errorAt, description, errString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🔴 Error log reported on slack at %s", time.Now())
}

// Notify success logs on Slack
func (n *Notifier) NotifySuccess(successAt, description, successString string, fields ...string) {
	if err := n.NotifySuccessE(successAt, description, successString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🟢 Success log reported on slack at %s", time.Now())
}

// Notify success logs on Slack
func (n *Notifier) NotifyWarn(warnAt, description, warnString string, fields ...string) {
	if err := n.NotifyWarnE(warnAt, description, warnString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🟡 Warn log reported on slack at %s", time.Now())
}

// Notify error logs on Slack and returns error
func (n *Notifier) NotifyErrorE(errorAt, description, errString string, fields ...string) error {
	if !isConfigured(n.config.Error) {
		return errors.New("❌ Slack error config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	attachmentFields := []slack.AttachmentField{
		{Title: "ErrorAt", Value: errorAt},
		{Title: "Description", Value: description},
	}

	for i := 0; i < len(fields); i += 2 {
		attachmentFields = append(attachmentFields, slack.AttachmentField{
			Title: fields[i],
			Value: fields[i+1],
		})
	}

	if len(errString) < 4000 {
		err := n.sendOnSlack(errString, errorColor, n.config.Error, attachmentFields...)
		if err != nil {
			return fmt.Errorf("❌ Failed to report error on slack: %s", err)
		}
	} else {
		err := n.sendOnSlackAsFile(errString, errorColor, n.config.Error, attachmentFields...)
		if err != nil {
			return fmt.Errorf("❌ Failed to report error on slack: %s", err)
		}
	}

	return nil
}

// Notify success logs on Slack and returns error
func (n *Notifier) NotifySuccessE(successAt, description, successString string, fields ...string) error {
	if !isConfigured(n.config.Success) {
		return errors.New("❌ Slack success config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	attachmentFields := []slack.AttachmentField{
		{Title: "SuccessAt", Value: successAt},
		{Title: "Description", Value: description},
	}

	for i := 0; i < len(fields); i += 2 {
		attachmentFields = append(attachmentFields, slack.AttachmentField{
			Title: fields[i],
			Value: fields[i+1],
		})
	}

	err := n.sendOnSlack(successString, successColor, n.config.Success, attachmentFields...)
	if err != nil {
		return fmt.Errorf("❌ Failed to report success on slack: %s", err)
	}

	return nil
}

// Notify warn logs on Slack
func (n *Notifier) NotifyWarnE(warnAt, description, warnString string, fields ...string) error {
	if !isConfigured(n.config.Warn) {
		return errors.New("❌ Slack warn config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	attachmentFields := []slack.AttachmentField{
		{Title: "WarnAt", Value: warnAt},
		{Title: "Description", Value: description},
	}

	for i := 0; i < len(fields); i += 2 {
		attachmentFields = append(attachmentFields, slack.AttachmentField{
			Title: fields[i],
			Value: fields[i+1],
		})
	}

	err := n.sendOnSlack(warnString, warnColor, n.config.Warn, attachmentFields...)
	if err != nil {
		return fmt.Errorf("❌ Failed to report warn on slack: %s", err)
	}

	return nil
}

// sendOnSlackAsFile sends text as file on slack channel
func (n *Notifier) sendOnSlackAsFile(text, messageColor string, channelConfig *SlackChannelConfig, fields ...slack.AttachmentField) error {
	err := n.sendOnSlack("", messageColor, channelConfig, fields...)
	if err != nil {
		return err
	}

	// Create the Slack attachment that we will send to the channel
	fileattachment := slack.FileUploadParameters{
		Content:  text,
		Channels: []string{channelConfig.ChannelID},
	}

	_, err = n.slackClient.UploadFile(fileattachment)
	if err != nil {
		return fmt.Errorf("failed to upload file : %s", err)
	}

	return nil
}

func (n *Notifier) sendOnSlack(text, messageColor string, channelConfig *SlackChannelConfig, fields ...slack.AttachmentField) error {
	mentions := generateMentions(channelConfig.Mentions)

	// Create the Slack attachment that we will send to the channel
	attachment := slack.Attachment{
		Pretext: mentions,
		Text:    text,
		Color:   messageColor,
		Fields:  append([]slack.AttachmentField{}, fields...),
		Footer:  time.Now().Format("2006-01-02 15:04:05"),
	}
	_, _, err := n.slackClient.PostMessage(
		channelConfig.ChannelID,
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		return fmt.Errorf("failed to send message : %s", err)
	}

	return nil
}

func isConfigured(channelConfig *SlackChannelConfig) bool {
	if channelConfig == nil {
		return false
	} else if channelConfig.ChannelID == "" {
		return false
	}

	return true
}

func generateMentions(users []string) string {
	str := ""
	for _, user := range users {
		str += fmt.Sprintf("<@%s> ", user)
	}

	return str
}
