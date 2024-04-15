package config

type MySQL struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Schema   string `json:"schema" yaml:"schema"`
}
