package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type ComputeRes struct {
	Id string `json:"Id"`

	State string `json:"State"`

	OwnerOrg string `json:"OwnerOrg"`
	UserOrg  string `json:"UserOrg"`

	UserOrgDueDate int `json:"UserOrgDueDate"`

	User string `json:"User"`

	OS   string `json:"OS"`
	Arch string `json:"Arch"`

	CPUSKU     string `json:"CPUSKU"`
	CPUSockets string `json:"CPUSockets"`
	CPUCores   string `json:"CPUCores"`

	GPUSKU string `json:"GPUSKU"`
	GPUNum string `json:"GPUNum"`

	MemorySize string `json:"MemorySize"`

	ConnectionAbilities string `json:"ConnectionAbilities"`

	Ip string `json:"Ip"`
}

func (c *ComputeRes) IsAvailable() bool {
	return c.User == ""
}

func (s *SmartContract) GetComputeRes(ctx contractapi.TransactionContextInterface, id string) (*ComputeRes, error) {

	org, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.readState(ctx, assetComputeRes, id)
	if err != nil {
		return nil, err
	}

	var asset ComputeRes
	err = json.Unmarshal(res, &asset)
	if err != nil {
		return nil, err
	}

	if asset.OwnerOrg != org && asset.UserOrg != org {
		return nil, fmt.Errorf("unauthorized access")
	}

	return &asset, nil
}

func (s *SmartContract) PutComputeRes(ctx contractapi.TransactionContextInterface, id string, res *ComputeRes) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return err
	}

	data, err := json.Marshal(*res)

	if err != nil {
		return err
	}

	if org != res.OwnerOrg {
		return fmt.Errorf("org != OwnerOrg")
	}

	return s.putState(ctx, assetComputeRes, id, data)
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
		Id:       id,
		State:    "unverified",
		OwnerOrg: org,
		UserOrg:  org,
	}

	assetJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}

	fmt.Println("CreateComputeRes", string(assetJSON))

	return s.putState(ctx, assetComputeRes, id, assetJSON)
}

func (s *SmartContract) ListComputeRes(ctx contractapi.TransactionContextInterface) ([]*ComputeRes, error) {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)
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

		if asset.OwnerOrg != org && asset.UserOrg != org {
			continue
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

	assetJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}

	return s.putState(ctx, assetComputeRes, id, assetJSON)
}

type ComputeResUpdate struct {
	Os         string `json:"os"`
	Arch       string `json:"arch"`
	CpuSKU     string `json:"cpusku"`
	CpuCores   string `json:"cpucores"`
	CpuSockets string `json:"cpusockets"`
	GpuSKU     string `json:"gpu"`
	GpuNum     string `json:"gpunums"`
	Network    string `json:"network"`
	Ip         string `json:"ip"`
	Ram        string `json:"ram"`
}

func (s *SmartContract) UpdateComputeRes(ctx contractapi.TransactionContextInterface, Id string, data string) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return err
	}

	usr, err := s.getUserInfo(ctx, org)

	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}

	if usr.Role != "admin" {
		return fmt.Errorf("user %s is not authorized to update a compute resource", usr.UserName)
	}

	res, err := s.readState(ctx, assetComputeRes, Id)

	if err != nil {
		return fmt.Errorf("can not find asset %s %v", Id, err)
	}

	var asset ComputeRes
	err = json.Unmarshal(res, &asset)
	if err != nil {
		return fmt.Errorf("unable to unmarshal asset %s %v", Id, err)
	}

	var u_res ComputeResUpdate
	err = json.Unmarshal([]byte(data), &u_res)
	if err != nil {
		return err
	}

	asset.Arch = u_res.Arch
	asset.CPUCores = u_res.CpuCores
	asset.CPUSKU = u_res.CpuSKU
	asset.CPUSockets = u_res.CpuSockets
	asset.ConnectionAbilities = u_res.Network
	asset.GPUNum = u_res.GpuNum
	asset.GPUSKU = u_res.GpuSKU
	asset.Ip = u_res.Ip

	asset.MemorySize = u_res.Ram
	asset.OS = u_res.Os

	asset.State = "verified"

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return s.putState(ctx, assetComputeRes, Id, assetJSON)

}

func (s *SmartContract) GetConnectDetails(ctx contractapi.TransactionContextInterface, Id string) (string, error) {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)

	if err != nil {
		return "", err
	}

	usr, err := s.getUserInfo(ctx, org)

	if err != nil {
		return "", err
	}

	Basset, err := s.readState(ctx, assetComputeRes, Id)

	if err != nil {
		return "", err
	}

	if usr.Role != "admin" && !contains(Id, usr.ComputeResList) {
		return "", fmt.Errorf("unauthorized access")
	}

	var asset ComputeRes
	err = json.Unmarshal(Basset, &asset)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal asset %s %v", Id, err)
	}

	if asset.UserOrg != org {
		return "", fmt.Errorf("permision deined")
	}

	return asset.Ip, nil
}

func (s *SmartContract) ListForRent(ctx contractapi.TransactionContextInterface, Id string, avail int, price int) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)

	if err != nil {
		return err
	}

	usr, err := s.getUserInfo(ctx, org)

	if err != nil {
		return err
	}

	Basset, err := s.readState(ctx, assetComputeRes, Id)

	if err != nil {
		return err
	}

	if usr.Role != "admin" {
		return fmt.Errorf("unauthorized access")
	}

	var asset ComputeRes
	err = json.Unmarshal(Basset, &asset)
	if err != nil {
		return fmt.Errorf("unable to unmarshal asset %s %v", Id, err)
	}

	if asset.OwnerOrg != org {
		return fmt.Errorf("you can only rent your own asset")
	}

	if asset.OwnerOrg != asset.UserOrg {
		return fmt.Errorf("the resource is already rented")
	}

	return nil

}

func (s *SmartContract) ClaimRent(ctx contractapi.TransactionContextInterface, id string) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)

	if err != nil {
		return err
	}

	asset, err := s.GetComputeRes(ctx, id)
	if err != nil {
		return err
	}
	if asset.OwnerOrg != org {
		return fmt.Errorf("only owner can claim a rent")
	}

	if asset.UserOrg == org {
		return fmt.Errorf("not rent")
	}

	_time, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	if time.UnixMicro(int64(asset.UserOrgDueDate)).Before(_time.AsTime()) {
		asset.UserOrg = org
		asset.UserOrgDueDate = 0
		return s.PutComputeRes(ctx, id, asset)
	} else {
		return fmt.Errorf("not time to claim")
	}

}
