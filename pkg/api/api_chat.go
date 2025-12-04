package api

import (
	"github.com/AntNoHuabei/Remo/pkg/api/request"
	"github.com/AntNoHuabei/Remo/pkg/api/response"
	"github.com/AntNoHuabei/Remo/pkg/chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Chat(ctx *gin.Context) {

	sendError := func(err error, requestId string) {
		ctx.SSEvent("error", response.ChatResponse{Error: err, RequestID: requestId})
		ctx.Writer.Flush()
		ctx.Abort()
	}

	var req request.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		sendError(err, req.RequestId)
		return
	}
	agent, err := chat.NewContinuousAgent(ctx)
	if err != nil {
		sendError(err, req.RequestId)
		return
	}

	err = agent.Recover(ctx, req.Session)
	if err != nil {
		sendError(err, req.RequestId)
		return
	}

	ctx.Status(http.StatusOK)
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	ctx.Writer.Flush()

	output, err := agent.Chat(ctx, &chat.Message{
		Content:   req.Message,
		Role:      "user",
		Session:   req.Session,
		RequestId: req.RequestId,
	})
	if err != nil {
		sendError(err, req.RequestId)
		return
	}
	for res := range output {

		ctx.SSEvent("message", res)
		ctx.Writer.Flush()
	}
}
