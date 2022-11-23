package notifier

// Notify error logs on Slack
func NotifyError(errorAt, description, errString string) {
	standardNotifier.NotifyError(errorAt, description, errString)
}

// Notify success logs on Slack
func NotifySuccess(successAt, description, successString string) {
	standardNotifier.NotifySuccess(successAt, description, successString)
}

// Notify success logs on Slack
func NotifyWarn(warnAt, description, warnString string) {
	standardNotifier.NotifyWarn(warnAt, description, warnString)
}

// Notify error logs on Slack
func NotifyErrorE(errorAt, description, errString string) error {
	return standardNotifier.NotifyErrorE(errorAt, description, errString)
}

// Notify success logs on Slack
func NotifySuccessE(successAt, description, successString string) error {
	return standardNotifier.NotifySuccessE(successAt, description, successString)
}

// Notify success logs on Slack
func NotifyWarnE(warnAt, description, warnString string) error {
	return standardNotifier.NotifyWarnE(warnAt, description, warnString)
}
