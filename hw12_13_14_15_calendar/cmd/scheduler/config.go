package main

type Config struct {
	Logger     LoggerConf     `mapstructure:"logger"`
	DB         DBConf         `mapstructure:"db"`
	Repository RepositoryConf `mapstructure:"repository"`
	Cleaner    CleanerConf    `mapstructure:"cleaner"`
	Rabbit     RabbitConf     `mapstructure:"rabbit"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type DBConf struct {
	DSN string `mapstructure:"dsn"`
}

type RepositoryConf struct {
	Type string `mapstructure:"type"`
}

type CleanerConf struct {
	EventLifespan string `mapstructure:"event_lifespan"`
}

type RabbitConf struct {
	DSN   string `mapstructure:"dsn"`
	Queue string `mapstructure:"queue"`
}
