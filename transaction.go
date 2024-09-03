package main

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/rs/zerolog/log"
)

func Query(a1 string, args ...string) ([]byte, error) {

	data, err := contract.EvaluateTransaction(a1, args...)
	checkErr(err)
	log.Info().Str("func", a1).Strs("args", args).AnErr("err", err).Msg("Query")

	return data, err
}

func Invoke(a1 string, args ...string) ([]byte, error) {

	data, err := contract.SubmitTransaction(a1, args...)
	checkErr(err)

	log.Info().Str("func", a1).Strs("args", args).AnErr("err", err).Msg("Invoke")

	return data, err
}
func InvokeTransistent(a1 string, data map[string][]byte, args ...string) ([]byte, error) {
	res, err := contract.Submit(a1, client.WithTransient(data), client.WithArguments(args...))
	checkErr(err)

	log.Info().Str("func", a1).Strs("args", args).Any("data", data).AnErr("err", err).Msg("InvokeTransistent")

	return res, err
}
