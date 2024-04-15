package config

type Config struct {
	IP    string `json:"ip"`
	Port  uint64 `json:"port" validate:"range:[1,65535],default:80"`
	MySQL MySQL  `json:"mysql"`
}
