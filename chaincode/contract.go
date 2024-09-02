package main

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

const (
	userKeyType = "User"
	resKeyType  = "ComputeRes"

	rootuser = "RootUser"

	assetComputeRes = "assetComputeRes"
	assetUser       = "assetUsers"
)
