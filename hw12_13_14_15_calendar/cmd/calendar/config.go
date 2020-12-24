package main

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf     `mapstructure:"logger"`
	Server     ServerConf     `mapstructure:"server"`
	Repository RepositoryConf `mapstructure:"repository"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type ServerConf struct {
	Port string `mapstructure:"port"`
}

type RepositoryConf struct {
	Type        string `mapstructure:"type"`
	Credentials string `mapstructure:"credentials"`
}
