package config

type config struct {
	MySQL   mySQL
	OSS     oss
	Redis   redis
	Smtp    smtp
	Service service
	Api     api
	MongoDb mongodb
	Log     log
}
type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type redis struct {
	Addr     string
	Username string
	Password string
}

type oss struct {
	Bucket    string
	AccessKey string
	SecretKey string
	Domain    string
	Region    string
}

type smtp struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	FromName string `mapstructure:"from_name"`
}

type service struct {
	Address    string `mapstructure:"address"`
	PrivateKey string `mapstructure:"private-key"`
}

type api struct {
	Key   string
	Model string
	Url   string
}

type mongodb struct {
	Addr     string
	Database string
	Username string
	Password string
}

type log struct {
	LogPath      string `mapstructure:"log_path"`
	LogPrefix    string `mapstructure:"log_prefix"`
	GinLogPrefix string `mapstructure:"gin_log_prefix"`
}
