package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/v2/pkg/statebased"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
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

func (s *SmartContract) readState(ctx contractapi.TransactionContextInterface, _type string, id string) ([]byte, error) {
	assetJSON, err := ctx.GetStub().GetPrivateData(_type, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	return assetJSON, nil
}

func (s *SmartContract) putState(ctx contractapi.TransactionContextInterface, _type string, id string, data []byte) error {
	return ctx.GetStub().PutPrivateData(_type, id, data)
}

func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func (s *SmartContract) CreateRootUser(ctx contractapi.TransactionContextInterface) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return err
	}

	created, _ := s.readState(ctx, assetUser, rootuser)

	if created != nil {
		return fmt.Errorf("the root user already exists")
	}

	thisuser, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("failed to get submitting client identity: %v", err)
	}

	uid, err := ctx.GetStub().CreateCompositeKey(userKeyType, []string{org, thisuser})

	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	root := User{
		UserName: thisuser,
		Role:     "admin",
	}

	assetJSON, err := json.Marshal(root)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("CreateRootUser", assetJSON)
	s.putState(ctx, assetUser, rootuser, assetJSON)

	return s.putState(ctx, assetUser, uid, assetJSON)
}

func verifyClientOrgMatchesPeerOrg(ctx contractapi.TransactionContextInterface) (string, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed getting the peer's MSPID: %v", err)
	}

	if clientMSPID != peerMSPID {
		return "", fmt.Errorf("client from org %v is not authorized to read or write private data from an org %v peer", clientMSPID, peerMSPID)
	}

	return clientMSPID, nil
}

func setAssetStateBasedEndorsement(ctx contractapi.TransactionContextInterface, auctionID string, orgToEndorse string) error {

	endorsementPolicy, err := statebased.NewStateEP(nil)
	if err != nil {
		return err
	}
	err = endorsementPolicy.AddOrgs(statebased.RoleTypePeer, orgToEndorse)
	if err != nil {
		return fmt.Errorf("failed to add org to endorsement policy: %v", err)
	}
	policy, err := endorsementPolicy.Policy()
	if err != nil {
		return fmt.Errorf("failed to create endorsement policy bytes from org: %v", err)
	}
	err = ctx.GetStub().SetStateValidationParameter(auctionID, policy)
	if err != nil {
		return fmt.Errorf("failed to set validation parameter on auction: %v", err)
	}

	return nil
}

// uid, role, err
func (s *SmartContract) getUserInfo(ctx contractapi.TransactionContextInterface, org string) (*User, error) {
	user, err := s.GetSubmittingClientIdentity(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to read submitting client identity: %v", err)
	}

	uid, err := ctx.GetStub().CreateCompositeKey(userKeyType, []string{org, user})

	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %v", err)
	}

	state, err := s.readState(ctx, assetUser, uid)

	if err != nil {
		return nil, fmt.Errorf("failed to read state: %v", err)
	}

	var u User
	err = json.Unmarshal(state, &u)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %v", err)
	}

	return &u, nil

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

		var asset ComputeRes
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("ListComputeRes", assets)

	return assets, nil
}
