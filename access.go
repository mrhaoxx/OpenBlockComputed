package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SSHAccessDetails struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Addr string `json:"addr"`
	Idn  string `json:"idn"`
}

func connectToBackend(c *gin.Context) {

	id := c.Param("id")

	details, err := Invoke("GetConnectDetails", id)

	if err != nil {
		log.Err(err).Str("data", string(details)).Msg("failed to get connection details")

		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	var sshDetails SSHAccessDetails

	err = json.Unmarshal(details, &sshDetails)

	if err != nil {
		log.Err(err).Str("data", string(details)).Msg("failed to unmarshal connection details")
	}

	log.Printf("ssh details: %+v", sshDetails)

	hdr := sshHandler{
		addr:    sshDetails.Addr,
		user:    sshDetails.User,
		secret:  sshDetails.Pass,
		keyfile: sshDetails.Idn,
	}

	hdr.webSocket(c.Writer, c.Request)
}
