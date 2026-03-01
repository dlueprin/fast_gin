package config

type config struct {
	ApiKey  string `yaml:"api_key"`
	BaseUrl string `yaml:"base_url"`
	Model   string `yaml:"model"`
}
type AI struct {
	LLM       config `yaml:"llm"`
	Embedding config `yaml:"embedding"`
}
