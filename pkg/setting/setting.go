package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	PrefixUrl      string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

var RedisSetting = &Redis{}

func Setup() {
	Cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	err = Cfg.Section("app").MapTo(AppSetting)
	logFatal(err, "AppSetting")
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	err = Cfg.Section("server").MapTo(ServerSetting)
	logFatal(err, "ServerSetting")
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	logFatal(err, "DatabaseSetting")

	err = Cfg.Section("redis").MapTo(RedisSetting)
	logFatal(err, "RedisSetting")
}

func logFatal(err error, settingKey string) {
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", settingKey, err)
	}
}
