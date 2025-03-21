package lib

import (
	"log"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	ServerPort              string `mapstructure:"SERVER_PORT"`
	Environment             string `mapstructure:"ENV"`
	LogOutput               string `mapstructure:"LOG_OUTPUT"`
	LogLevel                string `mapstructure:"LOG_LEVEL"`
	DBUsername              string `mapstructure:"DB_USER"`
	DBPassword              string `mapstructure:"DB_PASS"`
	DBHost                  string `mapstructure:"DB_HOST"`
	DBPort                  string `mapstructure:"DB_PORT"`
	DBName                  string `mapstructure:"DB_NAME"`
	RedisHost               string `mapstructure:"REDIS_HOST"`
	RedisPort               string `mapstructure:"REDIS_PORT"`
	RedisPassword           string `mapstructure:"REDIS_PASSWORD"`
	JWTSecret               string `mapstructure:"JWT_SECRET"`
	JwtRefreshSecret        string `mapstructure:"JWT_REFRESH_SECRET"`
	JwtResetSecret          string `mapstructure:"JWT_RESET_SECRET"`
	FirebaseConfigPath      string `mapstructure:"FIREBASE_CONFIG_PATH"`
	GeolocationAPIKey       string `mapstructure:"GEOLOCATION_API_KEY"`
	TokenLifetime           int    `mapstructure:"TOKEN_LIFETIME"`
	RefreshTokenLifetime    int    `mapstructure:"REFRESH_TOKEN_LIFETIME"`
	TCPSocketURL            string `mapstructure:"TCP_SOCKET_URL"`
	TCPSocketPort           string `mapstructure:"TCP_SOCKET_PORT"`
	SMTPHost                string `mapstructure:"SMTP_HOST"`
	SMTPPort                int    `mapstructure:"SMTP_PORT"`
	SMTPUser                string `mapstructure:"SMTP_USER"`
	SMTPPass                string `mapstructure:"SMTP_PASS"`
	SMTPSenderName          string `mapstructure:"SMTP_SENDER_NAME"`
	ClientSecretKey         string `mapstructure:"CLIENT_SECRET_KEY"`
	AWSBucketName           string `mapstructure:"AWS_BUCKET_NAME"`
	AWSRegion               string `mapstructure:"AWS_REGION"`
	AWSAccessKeyID          string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey      string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	VOICES_PATH             string `mapstructure:"VOICES_PATH"`
	FOTA_PATH               string `mapstructure:"FOTA_PATH"`
	MQTTBroker              string `mapstructure:"MQTT_BROKER"`
	MQTTUsername            string `mapstructure:"MQTT_USERNAME"`
	MQTTPassword            string `mapstructure:"MQTT_PASSWORD"`
	MQTTClientID            string `mapstructure:"MQTT_CLIENT_ID"`
	NotifyTopic             string `mapstructure:"NOTIFY_TOPIC"`
	SoundPushTopic          string `mapstructure:"SOUND_PUSH_TOPIC"`
	SoundPushHistoryTopic   string `mapstructure:"SOUND_PUSH_HISTORY_TOPIC"`
	PublishDevInfoTopic     string `mapstructure:"PUBLISH_DEV_INFO_TOPIC"`
	SubscribeDevInfoTopic   string `mapstructure:"SUBSCRIBE_DEV_INFO_TOPIC"`
	PublishDeviceManagement string `mapstructure:"PUBLISH_DEVICE_MANAGEMENT_TOPIC"`
	FotaPushTopic           string `mapstructure:"FOTA_PUSH_TOPIC"`
	TtsPushTopic            string `mapstructure:"TTS_PUSH_TOPIC"`
	CaCertificate           string `mapstructure:"CA_CERT_PATH"`
	Certificate             string `mapstructure:"CERT_PATH"`
	Key                     string `mapstructure:"KEY_PATH"`
}

// NewEnv creates a new environment
func NewEnv() Env {

	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	return env
}
