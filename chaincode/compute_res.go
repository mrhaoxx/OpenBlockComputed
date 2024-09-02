package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type ComputeRes struct {
	Id string `json:"Id"`

	State string `json:"State"`

	OwnerOrg string `json:"OwnerOrg"`
	UserOrg  string `json:"UserOrg"`

	UserOrgDueDate int `json:"UserOrgDueDate"`

	User        string `json:"User"`
	UserDueDate int    `json:"UserDueDate"`

	CPUSKU string `json:"CPUSKU"`
	CPUNum int    `json:"CPUNum"`

	GPUSKU string `json:"GPUSKU"`
	GPUNum int    `json:"GPUNum"`

	MemorySize int `json:"MemorySize"`

	ConnectionAbilities string `json:"ConnectionAbilities"`
}

func (c *ComputeRes) IsAvailable() bool {
	return c.User == ""
}

func (s *SmartContract) CreateComputeRes(ctx contractapi.TransactionContextInterface, id string) error {

	org, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return err
	}

	user, err := s.getUserInfo(ctx, org)

	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}

	if user.Role != "admin" {
		return fmt.Errorf("user %s is not authorized to create a compute resource", user.UserName)
	}

	rid, err := ctx.GetStub().CreateCompositeKey(resKeyType, []string{id})

	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	existing, err := s.readState(ctx, assetComputeRes, rid)
	if err == nil && existing != nil {
		return fmt.Errorf("the asset %s(%s) already exists", id, rid)
	}

	res := ComputeRes{
		Id:                  id,
		State:               "unverified",
		OwnerOrg:            org,
		UserOrg:             org,
		UserOrgDueDate:      0,
		User:                "",
		UserDueDate:         0,
		CPUSKU:              "",
		CPUNum:              0,
		GPUSKU:              "",
		GPUNum:              0,
		MemorySize:          0,
		ConnectionAbilities: "",
	}

	assetJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}

	fmt.Println("CreateComputeRes", string(assetJSON))

	return s.putState(ctx, assetComputeRes, id, assetJSON)
}

func (s *SmartContract) ListComputeRes(ctx contractapi.TransactionContextInterface) ([]*ComputeRes, error) {
	_, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return nil, err
	}

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(assetComputeRes, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*ComputeRes
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		fmt.Println("ListComputeRes", string(queryResponse.Value))
		var asset ComputeRes
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	fmt.Println("ListComputeRes", assets)

	return assets, nil
}

func (s *SmartContract) AssignUser(ctx contractapi.TransactionContextInterface, id string, user string, userDueDate int) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return err
	}

	usr, err := s.getUserInfo(ctx, org)

	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}

	if usr.Role != "admin" {
		return fmt.Errorf("user %s is not authorized to assign a compute resource", usr.UserName)
	}

	asset, err := s.readState(ctx, assetComputeRes, id)
	if err != nil {
		return err
	}

	if asset == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	var res ComputeRes
	err = json.Unmarshal(asset, &res)
	if err != nil {
		return err
	}

	res.User = user
	res.UserDueDate = userDueDate

	assetJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}

	return s.putState(ctx, assetComputeRes, id, assetJSON)
}
