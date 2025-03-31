package cfg

var Cfg Config

type AppConfig struct {
	Port                    int `yaml:"port"`
	RoomEmptyCleanupMinutes int `yaml:"roomEmptyCleanupMinutes"`
	RoomIdleCleanupMinutes  int `yaml:"roomIdleCleanupMinutes"`
}

type Config struct {
	App AppConfig `yaml:"app"`
}
