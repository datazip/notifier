package notifier

// Notify error logs on Slack
func NotifyError(errorAt, description, errString string, fields ...string) {
	standardNotifier.NotifyError(errorAt, description, errString, fields...)
}

// Notify success logs on Slack
func NotifySuccess(successAt, description, successString string, fields ...string) {
	standardNotifier.NotifySuccess(successAt, description, successString, fields...)
}

// Notify success logs on Slack
func NotifyWarn(warnAt, description, warnString string, fields ...string) {
	standardNotifier.NotifyWarn(warnAt, description, warnString, fields...)
}

// Notify error logs on Slack
func NotifyErrorE(errorAt, description, errString string, fields ...string) error {
	return standardNotifier.NotifyErrorE(errorAt, description, errString, fields...)
}

// Notify success logs on Slack
func NotifySuccessE(successAt, description, successString string, fields ...string) error {
	return standardNotifier.NotifySuccessE(successAt, description, successString, fields...)
}

// Notify success logs on Slack
func NotifyWarnE(warnAt, description, warnString string, fields ...string) error {
	return standardNotifier.NotifyWarnE(warnAt, description, warnString, fields...)
}
