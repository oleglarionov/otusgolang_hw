package main

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf     `mapstructure:"logger"`
	HTTPServer HTTPServerConf `mapstructure:"http_server"`
	GrpcServer GrpcServerConf `mapstructure:"grpc_server"`
	Repository RepositoryConf `mapstructure:"repository"`
	DB         DBConf         `mapstructure:"db"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type HTTPServerConf struct {
	Port string `mapstructure:"port"`
}

type GrpcServerConf struct {
	Port string `mapstructure:"port"`
}

type RepositoryConf struct {
	Type string `mapstructure:"type"`
}

type DBConf struct {
	DSN string `mapstructure:"dsn"`
}
