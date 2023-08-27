package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"san_francisco/config"
	"san_francisco/model"
	"san_francisco/utils"
	"strconv"
	"time"
)

func Ping(c *gin.Context) {
	// Start tracing span
	span := utils.BuildSpan(c.Request.Context(), "Ping", oteltrace.WithAttributes(attribute.Key("Request-ID").String(c.GetHeader("Request-ID"))))
	defer span.End()

	startTime, _ := c.Get("startTime")
	c.AbortWithStatusJSON(200, model.Response{
		Status:    "SUCCESS",
		Ping:      strconv.FormatInt(time.Now().Sub(startTime.(time.Time)).Milliseconds(), 10) + "ms",
		Gateway:   "SanFrancisco v" + config.Version,
		Service:   "SanFrancisco v" + config.Version,
		Timestamp: time.Now().Format("Mon Jan 02 15:04:05 MST 2006"),
		Data:      json.RawMessage("{\"message\": \"SanFrancisco v" + config.Version + " is online!\"}"),
	})
}
