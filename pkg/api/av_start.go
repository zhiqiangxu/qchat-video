package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/qchat-video/pkg/instance"
	"github.com/zhiqiangxu/qchat-video/pkg/udp"
	"github.com/zhiqiangxu/qchat/pkg/core"
)

type (
	// AVStartInput for input
	AVStartInput struct {
		Session udp.Session
	}

	// AVStartOutput for output
	AVStartOutput struct {
		core.BaseResp
	}
)

// AVStart for av start
func AVStart(c *gin.Context) {
	var input AVStartInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var output AVStartOutput
	err := instance.UDPServer().AVStart(input.Session)
	if err != nil {
		output.SetBase(core.ErrAPI, err.Error())
	}

	c.JSON(200, output)
}
