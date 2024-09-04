package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type BuyerInfo struct {
	Org   string `json:"org"`
	Price int    `json:"price"`
	Date  int    `json:"date"`
}

type ResMarket struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Date   int    `json:"date"`

	Res        ComputeRes `json:"resource"`
	Price      int        `json:"price"`
	Duration   int        `json:"duration"`
	MarketType string     `json:"type"`

	OwnerOrg string `json:"ownerOrg"`

	Winner string               `json:"winner"`
	Buyers map[string]BuyerInfo `json:"buyers"`
}

func (s *SmartContract) getResMarketElement(ctx contractapi.TransactionContextInterface, id string) (*ResMarket, error) {

	element, err := s.readState(ctx, assetMarket, id)
	if err != nil {
		return nil, err
	}

	var res ResMarket

	err = json.Unmarshal(element, &res)

	if err != nil {
		return nil, err
	}

	return &res, err
}

func (s *SmartContract) putResMarketElement(ctx contractapi.TransactionContextInterface, id string, res *ResMarket) error {
	data, err := json.Marshal(*res)
	if err != nil {
		return err
	}

	return s.putState(ctx, assetMarket, id, data)
}

func (s *SmartContract) putOnMarket(ctx contractapi.TransactionContextInterface, res ResMarket) (string, error) {

	id := ctx.GetStub().GetTxID()

	var err error

	res.Id = id

	if res.OwnerOrg == "" {
		res.OwnerOrg, err = ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return "", err
		}
	}
	if res.MarketType == "" {
		res.MarketType = "rent"
	}

	res.Status = "open"

	_time, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", err
	}
	res.Date = int(_time.AsTime().UnixMicro())

	res.Buyers = make(map[string]BuyerInfo)

	if err != nil {
		return "", err
	}

	err = s.putResMarketElement(ctx, id, &res)

	if err != nil {

		return "", err
	}

	return id, err
}

func (s *SmartContract) removeFromMarket(ctx contractapi.TransactionContextInterface, id string) error {
	res, err := s.getResMarketElement(ctx, id)
	if err != nil {
		return err
	}

	org, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if res.OwnerOrg != org {
		return fmt.Errorf("only the releaser can remove the element from market")
	}

	if res.Status != "open" && res.Status != "ended" {
		return fmt.Errorf("you can remove market element at this status %s", res.Status)

	}

	ctx.GetStub().DelPrivateData(assetMarket, id)
	return nil
}

func (s *SmartContract) MakePrice(ctx contractapi.TransactionContextInterface, id string, price int) error {
	res, err := s.getResMarketElement(ctx, id)

	if err != nil {
		return err
	}

	org, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return err
	}

	if res.Status != "open" {
		return fmt.Errorf("the market status can't be modified %s", res.Status)
	}

	if res.OwnerOrg == org {
		return fmt.Errorf("you can't make price on your own trades")
	}

	if price < res.Price {
		return fmt.Errorf("you can't make price lower than owner's price")
	}

	_time, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	res.Buyers[org] = BuyerInfo{
		Org:   org,
		Price: price,
		Date:  int(_time.AsTime().UnixMicro()),
	}

	return s.putResMarketElement(ctx, id, res)
}

func (s *SmartContract) LockMarketElement(ctx contractapi.TransactionContextInterface, id string, winner string, price int) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)

	if err != nil {
		return err
	}

	res, err := s.getResMarketElement(ctx, id)

	if err != nil {
		return err
	}

	if res.Res.OwnerOrg != org {
		return fmt.Errorf("only owner can lock")
	}

	if res.Status != "open" {
		return fmt.Errorf("can only lock a opening market element")
	}

	res.Status = "locked"

	by, ok := res.Buyers[winner]
	if !ok {
		return fmt.Errorf("winner not found")
	}

	if by.Price != price {
		return fmt.Errorf("price not consistent")
	}

	res.Winner = winner

	return s.putResMarketElement(ctx, id, res)

}

func (s *SmartContract) PutOnMarket(ctx contractapi.TransactionContextInterface, asset_id string, duration int, price int) (string, error) {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)

	if err != nil {
		return "", err
	}

	asset, err := s.GetComputeRes(ctx, asset_id)

	if err != nil {
		return "", err
	}

	if asset.OwnerOrg != org {
		return "", fmt.Errorf("only owner can market a res")
	}

	asset.SSHAccessDetails = SSHAccessDetails{}
	asset.AccessLogs = []Access{}

	return s.putOnMarket(ctx, ResMarket{
		Status:   "open",
		Res:      *asset,
		Price:    price,
		Duration: duration,
		OwnerOrg: org,
		Winner:   "",
	})

}

func (s *SmartContract) EndMarketElement(ctx contractapi.TransactionContextInterface, id string) error {
	org, err := verifyClientOrgMatchesPeerOrg(ctx)

	if err != nil {
		return err
	}
	res, err := s.getResMarketElement(ctx, id)
	if err != nil {
		return err
	}

	if res.OwnerOrg != org {
		return fmt.Errorf("only owner can end a market element")
	}

	if res.Status != "locked" {
		return fmt.Errorf("can only end a locked element")
	}

	res.Res.UserOrg = res.Winner
	_time, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	res.Res.UserOrgDueDate = int(_time.AsTime().Add(time.Duration(res.Duration)).UnixMicro())

	res.Status = "ended"

	err = s.putResMarketElement(ctx, id, res)
	if err != nil {
		return err
	}

	return s.PutComputeRes(ctx, res.Res.Id, &res.Res)
}

func (s *SmartContract) ListMarketElements(ctx contractapi.TransactionContextInterface) ([]*ResMarket, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(assetMarket, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var res []*ResMarket
	for resultsIterator.HasNext() {
		result, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var element ResMarket
		err = json.Unmarshal(result.Value, &element)
		if err != nil {
			return nil, err
		}
		if element.Status != "ended" {
			res = append(res, &element)
		}
	}

	return res, nil
}

func (s *SmartContract) GetMarketElement(ctx contractapi.TransactionContextInterface, id string) (ResMarket, error) {
	res, err := s.getResMarketElement(ctx, id)
	if err != nil {
		return ResMarket{}, err
	}
	return *res, nil
}
