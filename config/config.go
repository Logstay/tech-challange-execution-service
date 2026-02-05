package config

import (
	"os"
	"time"

	"github.com/Logstay/tech-challange-wherehouse-service/pkg/logger"
	"github.com/caarlos0/env/v11"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger      logger.Logger
	Database    DatabaseConfig
	Server      ServerConfig
	Email       EmailConfig
	Tracer      TracingConfig
	MercadoPago MercadoPagoConfig
	Kafka       KafkaConfig
	NewRelic    NewRelicConfig
	JWT         JWTConfig
}

type JWTConfig struct {
	Secret string `env:"JWT_SECRET,required"`
}

type DatabaseConfig struct {
	Host     string        `env:"DB_HOST,required"`
	Port     int           `env:"DB_PORT" envDefault:"27017"`
	User     string        `env:"DB_USER,required"`
	Password string        `env:"DB_PASSWORD,required"`
	Database string        `env:"DB_NAME,required"`
	Timeout  time.Duration `env:"DB_TIMEOUT,required"`
}

type ServerConfig struct {
	Host            string `env:"REST_SERVER_HOST,required" envDefault:""`
	Port            string `env:"REST_SERVER_PORT,required" envDefault:"8080"`
	SwaggerBasePath string `env:"SWAGGER_BASE_PATH,required" envDefault:"/v1"`
}

type EmailConfig struct {
	APIKey      string `env:"EMAIL_API_KEY"`
	SenderName  string `env:"EMAIL_SENDER_NAME"`
	SenderEmail string `env:"EMAIL_SENDER_EMAIL"`
}

type TracingConfig struct {
	ServiceName    string `env:"TRACING_SERVICE_NAME" envDefault:"tech-challange-wherehouse-service"`
	ServiceVersion string `env:"TRACING_SERVICE_VERSION" envDefault:"1.0.0"`
}

type MercadoPagoConfig struct {
	AccessToken string `env:"MERCADO_PAGO_ACCESS_TOKEN,required"`
}

type NewRelicConfig struct {
	AppName    string `env:"NEW_RELIC_APP_NAME"`
	LicenseKey string `env:"NEW_RELIC_LICENSE_KEY"`
}

type KafkaConfig struct {
	ProducerHosts string `env:"KAFKA_PRODUCER_HOSTS,required"`
	ConsumerHosts string `env:"KAFKA_CONSUMER_HOSTS,required"`
	Username      string `env:"KAFKA_USERNAME"`
	Password      string `env:"KAFKA_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	cfg.initLogger()

	return cfg, nil
}

func (c *Config) initLogger() error {
	baseLogger := logrus.New()
	baseLogger.SetFormatter(&logrus.JSONFormatter{})
	baseLogger.SetOutput(os.Stdout)
	baseLogger.SetLevel(logrus.InfoLevel)

	c.Logger = logger.NewLogger(c.Tracer.ServiceName, baseLogger)
	return nil
}
