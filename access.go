package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func connectToBackend(c *gin.Context) {

	id := c.Param("id")

	ip, err := Query("GetConnectDetails", id)

	if err != nil {
		log.Err(err).Str("ip", string(ip)).Msg("failed to get connection details")

		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	hdr := sshHandler{
		addr:    string(ip),
		user:    "root",
		keyfile: "ssh.key",
	}

	hdr.webSocket(c.Writer, c.Request)
}
