/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/rs/zerolog"
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	InitWebServer()

	initLedger()
	createAsset(contract)
	getAllAssets(contract)

	select {}

}

// newGrpcConnection creates a gRPC connection to the Gateway server.

// This type of transaction would typically only be run once by an application the first time it was started after its
// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
func initLedger() {
	Invoke("CreateRootUser")

}

// Evaluate a transaction to query ledger state.
func getAllAssets(contract *client.Contract) {
	Query("ListComputeRes")
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createAsset(contract *client.Contract) {
	Invoke("CreateComputeRes", assetId)

}
