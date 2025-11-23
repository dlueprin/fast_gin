package config

type Image struct {
	Size int64  `yaml:"size"` //这里用int64是因为后面要乘2次1024转化为字节
	Path string `yaml:"path"`
}
