package api

import (
	"Ts3Panel/core"
	"io"

	"github.com/gin-gonic/gin"
)

// StreamEvents SSE 实时日志流
func StreamEvents(c *gin.Context) {
	// 设置 SSE 必需的 Header
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	subID, stream := core.SubscribeSSE()
	defer core.UnsubscribeSSE(subID)

	clientGone := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false // 客户端断开，停止流
		case msg, ok := <-stream:
			if !ok {
				return false // 通道关闭
			}
			// 发送事件: type, data
			c.SSEvent(msg.Type, msg.Data)
			return true
		}
	})
}
