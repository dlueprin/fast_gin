package config

type AI struct {
	ApiKey  string `yaml:"api_key"`
	BaseUrl string `yaml:"base_url"`
	Model   string `yaml:"model"`
}
