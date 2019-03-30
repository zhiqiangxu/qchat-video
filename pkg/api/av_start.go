package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/qchat-video/pkg/instance"
	"github.com/zhiqiangxu/qchat-video/pkg/udp"
)

type (
	// AVStartInput for input
	AVStartInput udp.AVStartInput

	// AVStartOutput for output
	AVStartOutput udp.AVStartOutput
)

// AVStart for av start
func AVStart(c *gin.Context) {

	var input AVStartInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := instance.UDPServer().AVStart(udp.AVStartInput(input))

	c.JSON(200, output)
}
