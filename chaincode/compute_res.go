package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type Access struct {
	AccessTime int    `json:"AccessTime"`
	AccessUser string `json:"AccessUser"`
}

type SSHAccessDetails struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Addr string `json:"addr"`
	Idn  string `json:"idn"`
}

type ComputeRes struct {
	Id string `json:"Id"`

	Name string `json:"Name"`

	State string `json:"State"`

	OwnerOrg string `json:"OwnerOrg"`
	UserOrg  string `json:"UserOrg"`

	UserOrgDueDate int `json:"UserOrgDueDate"`

	User string `json:"User"`

	Details ComputeResUpdate `json:"Details"`

	SSHAccessDetails SSHAccessDetails `json:"SSHAccessDetails"`

	AccessLogs []Access `json:"AccessLogs"`
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

	if asset.UserOrg != org {
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

	if org != res.UserOrg {
		return fmt.Errorf("org != UserOrg")
	}

	return s.putState(ctx, assetComputeRes, id, data)
}

func (s *SmartContract) CreateComputeRes(ctx contractapi.TransactionContextInterface, name string) (string, error) {

	org, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return "", err
	}

	user, err := s.getUserInfo(ctx, org)

	if err != nil {
		return "", fmt.Errorf("failed to get user info: %v", err)
	}

	if user.Role != "admin" {
		return "", fmt.Errorf("user %s is not authorized to create a compute resource", user.UserName)
	}

	id := ctx.GetStub().GetTxID()

	rid, err := ctx.GetStub().CreateCompositeKey(resKeyType, []string{id})

	if err != nil {
		return "", fmt.Errorf("failed to create composite key: %v", err)
	}

	existing, err := s.readState(ctx, assetComputeRes, rid)
	if err == nil && existing != nil {
		return "", fmt.Errorf("the asset %s(%s) already exists", id, rid)
	}

	res := ComputeRes{
		Id:       id,
		Name:     name,
		State:    "uninitialized",
		OwnerOrg: org,
		UserOrg:  org,

		AccessLogs:       []Access{},
		SSHAccessDetails: SSHAccessDetails{},
	}

	assetJSON, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	fmt.Println("CreateComputeRes", string(assetJSON))

	return id, s.putState(ctx, assetComputeRes, id, assetJSON)
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

		if asset.UserOrg != org {
			asset.User = ""
			asset.AccessLogs = []Access{}

			// asset.Details.Ip = ""
			asset.State = "rented"
			asset.Details = ComputeResUpdate{}
			asset.SSHAccessDetails = SSHAccessDetails{}
		}

		assets = append(assets, &asset)
	}

	fmt.Println("ListComputeRes", assets)

	return assets, nil
}

func (s *SmartContract) QueryComputeRes(ctx contractapi.TransactionContextInterface, id string) (ComputeRes, error) {
	a, b := s.GetComputeRes(ctx, id)
	a.SSHAccessDetails = SSHAccessDetails{}
	return *a, b
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
	Hostname   string `json:"hostname"`
}

func (s *SmartContract) UpdateComputeRes(ctx contractapi.TransactionContextInterface, Id string) error {
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

	asset, err := s.GetComputeRes(ctx, Id)

	if err != nil {
		return err
	}

	data, err := ctx.GetStub().GetTransient()

	if err != nil {
		return err
	}

	var u_res ComputeResUpdate

	u_r, ok := data["update"]
	if ok {
		err = json.Unmarshal(u_r, &u_res)
		if err != nil {
			return err
		}

		asset.Details = u_res
		asset.State = "normal"
	}

	var s_res SSHAccessDetails
	s_r, ok := data["ssh"]
	if ok {
		err = json.Unmarshal(s_r, &s_res)
		if err != nil {
			return err
		}

		asset.SSHAccessDetails = s_res
	}

	return s.PutComputeRes(ctx, Id, asset)
}

func (s *SmartContract) DelComputeRes(ctx contractapi.TransactionContextInterface, Id string) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return err
	}

	usr, err := s.getUserInfo(ctx, org)

	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}

	if usr.Role != "admin" {
		return fmt.Errorf("user %s is not authorized to delete a compute resource", usr.UserName)
	}

	asset, err := s.GetComputeRes(ctx, Id)

	if err != nil {
		return err
	}

	if org != asset.OwnerOrg {
		return fmt.Errorf("only owner can delete a compute resource")
	}

	if asset.OwnerOrg != asset.UserOrg {
		return fmt.Errorf("can't delete a rented compute resource")
	}

	return ctx.GetStub().DelPrivateData(assetComputeRes, Id)
}

func (s *SmartContract) GetConnectDetails(ctx contractapi.TransactionContextInterface, Id string) (SSHAccessDetails, error) {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)

	if err != nil {
		return SSHAccessDetails{}, err
	}

	usr, err := s.getUserInfo(ctx, org)

	if err != nil {
		return SSHAccessDetails{}, err
	}

	if usr.Role != "admin" && !contains(Id, usr.ComputeResList) {
		return SSHAccessDetails{}, fmt.Errorf("unauthorized access")
	}

	asset, err := s.GetComputeRes(ctx, Id)

	if err != nil {
		return SSHAccessDetails{}, err
	}

	times, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return SSHAccessDetails{}, err
	}

	asset.AccessLogs = append(asset.AccessLogs, Access{
		AccessTime: int(times.AsTime().UnixMicro()),
		AccessUser: usr.UserName,
	})

	err = s.PutComputeRes(ctx, Id, asset)

	if err != nil {
		return SSHAccessDetails{}, err
	}

	return asset.SSHAccessDetails, nil
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

func (s *SmartContract) GetConnectionLogs(ctx contractapi.TransactionContextInterface, id string) ([]Access, error) {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)

	if err != nil {
		return nil, err
	}

	usr, err := s.getUserInfo(ctx, org)

	if err != nil {
		return nil, err
	}

	if usr.Role != "admin" {
		return nil, fmt.Errorf("unauthorized access")
	}

	asset, err := s.GetComputeRes(ctx, id)

	if err != nil {
		return nil, err
	}

	return asset.AccessLogs, nil
}
