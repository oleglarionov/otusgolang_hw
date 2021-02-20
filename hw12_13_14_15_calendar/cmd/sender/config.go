package main

type Config struct {
	Logger LoggerConf `mapstructure:"logger"`
	Rabbit RabbitConf `mapstructure:"rabbit"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type RabbitConf struct {
	DSN   string `mapstructure:"dsn"`
	Queue string `mapstructure:"queue"`
}
