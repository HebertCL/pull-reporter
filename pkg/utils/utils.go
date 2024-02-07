package utils

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

func LoadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Couldn't load env: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		kv := strings.SplitN(line, "=", 2)

		if len(kv) != 2 {
			err := errors.New("couldn't get key/value pair from env")
			return err
		}

		os.Setenv(kv[0], kv[1])
	}

	return scanner.Err()
}
