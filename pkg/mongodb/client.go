package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Logstay/tech-challange-wherehouse-service/config"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	client   *mongo.Client
	database *mongo.Database
	config   *config.Config
	log      logger.Logger
	mu       sync.RWMutex
}

func New(ctx context.Context, cfg *config.Config, log logger.Logger) (*Client, error) {

	connCtx, cancel := context.WithTimeout(ctx, cfg.Database.Timeout)
	defer cancel()

	uri := buildURI(cfg.Database)

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)

	if cfg.Database.User != "" && cfg.Database.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: cfg.Database.User,
			Password: cfg.Database.Password,
		})
	}

	client, err := mongo.Connect(connCtx, clientOptions)
	if err != nil {
		log.Error("Failed to connect to MongoDB", err, map[string]any{
			"host": cfg.Database.Host,
			"port": cfg.Database.Port,
		})
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = client.Ping(connCtx, readpref.Primary()); err != nil {
		log.Error("Failed to ping MongoDB", err, map[string]any{
			"host": cfg.Database.Host,
			"port": cfg.Database.Port,
		})
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Info("MongoDB connected successfully", map[string]any{
		"host":     cfg.Database.Host,
		"port":     cfg.Database.Port,
		"database": cfg.Database.Database,
	})

	return &Client{
		client:   client,
		database: client.Database(cfg.Database.Database),
		config:   cfg,
		log:      log,
	}, nil
}

func buildURI(cfg config.DatabaseConfig) string {
	if cfg.User != "" && cfg.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	}
	return fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
}

func (c *Client) GetDatabase() *mongo.Database {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.database
}

func (c *Client) GetCollection(name string) *mongo.Collection {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.database.Collection(name)
}

func (c *Client) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.Ping(pingCtx, readpref.Primary())
}

func (c *Client) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		return nil
	}

	c.log.Info("Closing MongoDB connection", nil)

	disconnectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := c.client.Disconnect(disconnectCtx)
	if err != nil {
		c.log.Error("Failed to disconnect from MongoDB", err, nil)
		return err
	}

	c.log.Info("MongoDB disconnected successfully", nil)
	return nil
}

func (c *Client) IsConnected(ctx context.Context) bool {
	if c.client == nil {
		return false
	}

	err := c.Ping(ctx)
	return err == nil
}
