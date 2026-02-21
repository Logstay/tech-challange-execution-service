package kafka

import (
	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/message-broker"
	"github.com/Logstay/tech-challange-execution-service/pkg/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Config struct {
	ProducerHosts     string
	ConsumerHosts     string
	Username          string
	Password          string
	SecurityProtocol  string
	SASLMechanism     string
	CompressionType   string
	EnableIdempotence bool
	MaxInFlightReqs   int
}

type Kafka struct {
	config Config
	logger logger.Logger
}

func New(config Config, logger logger.Logger) (k *Kafka, err error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	if config.SecurityProtocol == "" {
		config.SecurityProtocol = "PLAINTEXT"
	}
	if config.SASLMechanism == "" {
		config.SASLMechanism = "PLAIN"
	}
	if config.CompressionType == "" {
		config.CompressionType = "snappy"
	}
	if config.MaxInFlightReqs == 0 {
		config.MaxInFlightReqs = 5
	}

	logger.Info("Kafka adapter initialized", map[string]any{
		"producer_hosts":    config.ProducerHosts,
		"consumer_hosts":    config.ConsumerHosts,
		"security_protocol": config.SecurityProtocol,
		"compression":       config.CompressionType,
	})

	return &Kafka{
		config: config,
		logger: logger,
	}, nil
}

func validateConfig(config Config) error {

	if config.ProducerHosts == "" {
		return ports.InvalidConfigurationError("producer hosts cannot be empty")
	}

	if config.ConsumerHosts == "" {
		return ports.InvalidConfigurationError("consumer hosts cannot be empty")
	}

	if config.Username != "" && config.Password == "" {
		return ports.InvalidConfigurationError("password is required when username is provided")
	}

	return nil
}

func (k Kafka) NewProducer() (p ports.Producer, err error) {
	k.logger.Info("Creating Kafka producer", nil)
	kafkaProducer, err := kafka.NewProducer(k.buildProducerConfig())

	if err != nil {
		k.logger.Error("Failed to create Kafka producer", err)
		return
	}

	producer := Producer{
		Producer: kafkaProducer,
		logger:   k.logger,
	}

	k.logger.Info("Kafka producer created successfully", nil)
	return producer, nil
}

func (k Kafka) NewConsumer(groupID string) (ports.Consumer, error) {
	var err error

	k.logger.Info("Creating Kafka consumer", map[string]any{
		"group_id": groupID,
	})

	c := &Consumer{
		config: k.buildConsumerConfig(groupID),
		logger: k.logger,
	}

	if err = c.build(); err != nil {
		k.logger.Error("Failed to create Kafka consumer", err, map[string]any{
			"group_id": groupID,
		})
		c = nil
	}

	k.logger.Info("Kafka consumer created successfully", map[string]any{
		"group_id": groupID,
	})
	return c, err
}

func (k Kafka) buildConsumerConfig(groupID string) (kafkaConfig *kafka.ConfigMap) {
	kafkaConfig = &kafka.ConfigMap{
		"bootstrap.servers":         k.config.ConsumerHosts,
		"group.id":                  groupID,
		"auto.offset.reset":         "earliest",
		"enable.auto.commit":        false,
		"max.poll.interval.ms":      600000,
		"session.timeout.ms":        10000,
		"heartbeat.interval.ms":     3000,
		"fetch.min.bytes":           1,
		"fetch.wait.max.ms":         500,
		"max.partition.fetch.bytes": 1048576,
	}

	k.setCredentials(kafkaConfig)

	k.logger.Debug("Consumer config built", map[string]any{
		"group_id":     groupID,
		"auto_commit":  false,
		"offset_reset": "earliest",
	})

	return
}

func (k Kafka) setCredentials(kafkaConfig *kafka.ConfigMap) {
	if k.config.Username != "" && k.config.Password != "" {
		_ = kafkaConfig.SetKey("security.protocol", k.config.SecurityProtocol)
		_ = kafkaConfig.SetKey("sasl.mechanism", k.config.SASLMechanism)
		_ = kafkaConfig.SetKey("sasl.username", k.config.Username)
		_ = kafkaConfig.SetKey("sasl.password", k.config.Password)
		k.logger.Debug("Credentials configured", map[string]any{
			"username":          k.config.Username,
			"security_protocol": k.config.SecurityProtocol,
			"sasl_mechanism":    k.config.SASLMechanism,
		})
	} else {
		_ = kafkaConfig.SetKey("security.protocol", "PLAINTEXT")
	}
}

func (k Kafka) buildProducerConfig() (kafkaConfig *kafka.ConfigMap) {
	kafkaConfig = &kafka.ConfigMap{
		"bootstrap.servers":                     k.config.ProducerHosts,
		"acks":                                  "all",
		"compression.type":                      k.config.CompressionType,
		"linger.ms":                             10,
		"batch.size":                            16384,
		"max.in.flight.requests.per.connection": k.config.MaxInFlightReqs,
		"retries":                               3,
		"retry.backoff.ms":                      100,
		"request.timeout.ms":                    30000,
		"delivery.timeout.ms":                   120000,
	}

	if k.config.EnableIdempotence {
		_ = kafkaConfig.SetKey("enable.idempotence", true)

		if k.config.MaxInFlightReqs > 5 {
			_ = kafkaConfig.SetKey("max.in.flight.requests.per.connection", 5)
		}
	}

	k.setCredentials(kafkaConfig)

	k.logger.Debug("Producer config built", map[string]any{
		"compression": k.config.CompressionType,
		"idempotence": k.config.EnableIdempotence,
	})

	return
}
