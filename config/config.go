package config

import "github.com/spf13/viper"

type Config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`

	SMTPEmail    string `mapstructure:"SMTP_EMAIL"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`

	JWTIDTokenSecret      string `mapstructure:"ID_TOKEN_HMAC_SECRET"`
	JWTAccessTokenSecret  string `mapstructure:"ACCESS_TOKEN_HMAC_SECRET"`
	JWTRefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_HMAC_SECRET"`
}

func LoadConfig(path string) (c Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}
