package main

import "github.com/rs/zerolog/log"

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
