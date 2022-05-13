package useful_servise

import (
	"log"
	"regexp"
	"strconv"
)

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
