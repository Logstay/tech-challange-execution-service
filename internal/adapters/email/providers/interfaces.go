package providers

import (
	"context"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/email"
)

type Client interface {
	Send(ctx context.Context, e ports.Email, template string) (err error)
}
