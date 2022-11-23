package notifier

type Config struct {
	Token   string              `yaml:"token"`
	Success *SlackChannelConfig `yaml:"success"`
	Error   *SlackChannelConfig `yaml:"error"`
	Warn    *SlackChannelConfig `yaml:"warn"`
}

type SlackChannelConfig struct {
	ChannelID string   `yaml:"id"`
	Mentions  []string `yaml:"mentions"` // mentions are expected to have @ excluded before usernames
}
