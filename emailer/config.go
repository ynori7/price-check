package emailer

type Config struct {
	Enabled    bool
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
	From       EmailRecipient
	To         EmailRecipient
}

type EmailRecipient struct {
	Address string
	Name    string
}
