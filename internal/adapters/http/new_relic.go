package http

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"

	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
)

func NewRelic(app *newrelic.Application) gin.HandlerFunc {
	return nrgin.Middleware(app)
}
