package api

import (
	"github.com/AntNoHuabei/Remo/pkg/api/request"
	"github.com/AntNoHuabei/Remo/pkg/chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SessionCreate(c *gin.Context) {

	s := chat.CreateSession()

	c.JSON(http.StatusOK, Success(s))
}

func SessionList(c *gin.Context) {

	var req request.SessionListRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Fail(err.Error()))
		return
	} else {

		if req.Page < 1 {
			req.Page = 1
		}
		if req.Size < 1 {
			req.Size = 10
		}

		sessions, err := chat.SessionList((req.Page-1)*req.Size, req.Size)

		if err != nil {
			c.JSON(http.StatusInternalServerError, Fail(err.Error()))
		} else {
			c.JSON(http.StatusOK, Success(sessions))
		}

	}

}
func SessionDelete(c *gin.Context) {

	var req request.SessionDeleteRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Fail(err.Error()))
		return
	} else {
		err = chat.DeleteSession(req.Id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, Fail(err.Error()))
		} else {
			c.JSON(http.StatusOK, Success(nil))
		}
	}
}

func SessionMessages(c *gin.Context) {

	var req request.SessionMessagesRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Fail(err.Error()))
		return
	} else {
		messages, err := chat.Messages(req.Session)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Fail(err.Error()))
		} else {
			c.JSON(http.StatusOK, Success(messages))
		}
	}
}
