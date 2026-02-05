package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func NewRelicContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	if txn := nrgin.Transaction(c); txn != nil {
		ctx = newrelic.NewContext(ctx, txn)
		c.Request = c.Request.WithContext(ctx)
	}

	return ctx
}
