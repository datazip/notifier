package email

var standardMailer *Mailer

type Mailer struct {
	config *Config
}

func Init(config *Config) *Mailer {
	standardMailer = &Mailer{
		config: config,
	}

	return standardMailer
}
