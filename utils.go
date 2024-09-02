package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog/log"
)

// Submit transaction, passing in the wrong number of arguments ,expected to throw an error containing details of any error responses from the smart contract.
func checkErr(err error) {

	if err == nil {
		log.Debug().Msg("No error")
		return
	}

	var endorseErr *client.EndorseError
	var submitErr *client.SubmitError
	var commitStatusErr *client.CommitStatusError
	var commitErr *client.CommitError

	if errors.As(err, &endorseErr) {
		log.Err(endorseErr).Str("TransactionID", endorseErr.TransactionID).Msg("Endorse error")
	} else if errors.As(err, &submitErr) {
		fmt.Printf("Submit error for transaction %s with gRPC status %v: %s\n", submitErr.TransactionID, status.Code(submitErr), submitErr)
		log.Err(submitErr).Str("TransactionID", submitErr.TransactionID).Msg("Submit error")
	} else if errors.As(err, &commitStatusErr) {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Err(commitStatusErr).Str("TransactionID", commitStatusErr.TransactionID).Msg("Timeout waiting for commit status")
		} else {
			log.Err(commitStatusErr).Str("TransactionID", commitStatusErr.TransactionID).Msg("Error obtaining commit status")
		}
	} else if errors.As(err, &commitErr) {
		log.Err(commitErr).Str("TransactionID", commitErr.TransactionID).Msg("Commit error")
	} else {
		panic(fmt.Errorf("unexpected error type %T: %w", err, err))
	}

	// Any error that originates from a peer or orderer node external to the gateway will have its details
	// embedded within the gRPC status error. The following code shows how to extract that.
	statusErr := status.Convert(err)

	details := statusErr.Details()
	if len(details) > 0 {
		fmt.Println("Error Details:")

		for _, detail := range details {
			switch detail := detail.(type) {
			case *gateway.ErrorDetail:
				fmt.Printf("- address: %s, mspId: %s, message: %s\n", detail.Address, detail.MspId, detail.Message)
			}
		}
	}

}

// // Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
