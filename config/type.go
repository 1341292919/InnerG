package config

type config struct {
	MySQL     mySQL
	OSS       oss
	Redis     redis
	Smtp      smtp
	Service   service
	Api       api
	MongoDb   mongodb
	Log       log
	SnowFlake snowflake
	RabbitMq  rabbitmq
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
	Address           string `mapstructure:"address"`
	PrivateKey        string `mapstructure:"private-key"`
	WebsocketShardNum int    `mapstructure:"websocket-shard-num"`
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
	Level        string `mapstructure:"level"`
	LogPath      string `mapstructure:"log_path"`
	LogPrefix    string `mapstructure:"log_prefix"`
	GinLogPrefix string `mapstructure:"gin_log_prefix"`
	LogMaxDays   int    `mapstructure:"log_max_days"`
}
type snowflake struct {
	WorkerID      int64 `mapstructure:"worker-id"`
	DatancenterID int64 `mapstructure:"datancenter-id"`
}

type rabbitmq struct {
	Addr     string
	Username string
	Password string
}
