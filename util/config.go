package util

import (
	"github.com/spf13/viper"
	"time"
)

// The values are read by viper from the config file or environment variables
type Config struct {
	//Database config
	Environment  string `mapstructure:"ENVIRONMENT"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	//Server config
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	//Token config
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	//Email config
	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	//AWS config
	AwsBucketName string `mapstructure:"AWS_BUCKET_NAME"`
	AwsBucketRegion string `mapstructure:"AWS_BUCKET_REGION"`
	AwsAccessKeyID string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey string `mapstructure:"AWS_SECRET_KEY"`
}

// LoadConfig loads the config from the config file and environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
