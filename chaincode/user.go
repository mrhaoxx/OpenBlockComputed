package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type User struct {
	UserName string `json:"UserName"`

	Role string `json:"Role"`
	Org  string `json:"Org"`

	ComputeResList []string `json:"ComputeResList"`
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
