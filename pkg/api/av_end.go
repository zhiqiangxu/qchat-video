package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/qchat-video/pkg/instance"
	"github.com/zhiqiangxu/qchat-video/pkg/udp"
	"github.com/zhiqiangxu/qchat/pkg/core"
)

type (
	// AVEndInput for input
	AVEndInput struct {
		Session udp.Session
	}

	// AVEndOutput for output
	AVEndOutput struct {
		core.BaseResp
	}
)

// AVEnd for av end
func AVEnd(c *gin.Context) {

	var input AVEndInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var output AVEndOutput
	err := instance.UDPServer().AVEnd(input.Session)
	if err != nil {
		output.SetBase(core.ErrAPI, err.Error())
	}

	c.JSON(200, output)
}
