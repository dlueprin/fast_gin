package config

// 这里是这个包的入口
type Config struct {
	DB     DB     `yaml:"db"` //会先走到这里
	Redis  Redis  `yaml:"redis"`
	System System `yaml:"system"`
	Jwt    Jwt    `yaml:"jwt"`
	Image  Image  `yaml:"image"`
	Doc    Doc    `yaml:"doc"`
	Site   Site   `yaml:"site"`
	AI     AI     `yaml:"ai"`
	Chroma Chroma `yaml:"chroma"`
}
