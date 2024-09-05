package main

import (
	"errors"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/rs/zerolog/log"
)

func Query(a1 string, args ...string) ([]byte, error) {

	data, err := contract.EvaluateTransaction(a1, args...)
	msg := checkErr(err)
	log.Info().Str("func", a1).Strs("args", args).AnErr("err", err).Msg("Query")

	if err != nil {
		return data, errors.Join(err, errors.New(msg))
	}
	return data, nil
}

func Invoke(a1 string, args ...string) ([]byte, error) {

	data, err := contract.SubmitTransaction(a1, args...)
	msg := checkErr(err)

	log.Info().Str("func", a1).Strs("args", args).AnErr("err", err).Msg("Invoke")

	if err != nil {
		return data, errors.Join(err, errors.New(msg))
	}
	return data, nil
}
func InvokeTransistent(a1 string, data map[string][]byte, args ...string) ([]byte, error) {
	res, err := contract.Submit(a1, client.WithTransient(data), client.WithArguments(args...))

	msg := checkErr(err)

	log.Info().Str("func", a1).Strs("args", args).Any("data", data).AnErr("err", err).Msg("InvokeTransistent")

	if err != nil {
		return res, errors.Join(err, errors.New(msg))
	}
	return res, nil
}
