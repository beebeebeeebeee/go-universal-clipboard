package cfg

var Cfg Config

type AppConfig struct {
	Port      int    `yaml:"port"`
	ApiKey    string `yaml:"apiKey"`
	SignerKey string `yaml:"signerKey"`
}

type JWTConfig struct {
	PublicKey                     string `yaml:"publicKey"`
	PrivateKey                    string `yaml:"privateKey"`
	KID                           string `yaml:"kid"`
	AccessTokenDurationInMinute   int    `yaml:"accessTokenDurationInMinute"`
	RefreshTokenDurationInMinute  int    `yaml:"refreshTokenDurationInMinute"`
	RegisterTokenDurationInMinute int    `yaml:"registerTokenDurationInMinute"`
	Issuer                        string `yaml:"issuer"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	App      AppConfig      `yaml:"app"`
	JWT      JWTConfig      `yaml:"jwt"`
	Database DatabaseConfig `yaml:"database"`
}
