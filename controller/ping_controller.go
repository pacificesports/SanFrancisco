package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"san_francisco/config"
	"san_francisco/model"
	"strconv"
	"time"
)

func Ping(c *gin.Context) {
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
