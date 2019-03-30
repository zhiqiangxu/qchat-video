package api

import (
	"net/http"

	"github.com/zhiqiangxu/qchat-video/pkg/udp"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/qchat-video/pkg/instance"
)

type (
	// AVEndInput for input
	AVEndInput udp.AVEndInput

	// AVEndOutput for output
	AVEndOutput udp.AVEndOutput
)

// AVEnd for av end
func AVEnd(c *gin.Context) {

	var input AVEndInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := instance.UDPServer().AVEnd(udp.AVEndInput(input))

	c.JSON(200, output)
}
