package config

type Upload struct {
	Path         string   `json:"path"`
	ImageMaxSize int      `json:"imageMaxSize"`
	ImageType    []string `json:"imageType"`
}
