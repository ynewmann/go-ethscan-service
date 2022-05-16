package useful_servise

import (
	"log"
	"math/big"
	"regexp"
	"strconv"
)

const Eth = 1e18

func strictMatch(pattern string, b []byte) bool {
	pattern = "^" + pattern + "$"
	return match(pattern, b)
}

func match(pattern string, b []byte) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Error while match the path %s: %s\n", b, err)
		return false
	}

	return re.Match(b)
}

func getBlockNumber(path []byte) (uint64, error) {
	r, err := regexp.Compile("\\d+")
	if err != nil {
		log.Printf("Cannot get block number: %s\n", err)
		return 0, err
	}

	n, err := strconv.ParseUint(string(r.Find(path)), 10, 64)
	if err != nil {
		log.Printf("Cannot get block number: %s\n", err)
		return 0, err
	}

	return n, nil
}

func weiToEth(wei *big.Float) float64 {
	amount, _ := new(big.Float).Quo(wei, big.NewFloat(Eth)).Float64()

	return amount
}
