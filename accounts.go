package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func getAccs() []Acc {

	file, err := os.Open("accs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	account := Acc{}
	accs := []Acc{}
	for scanner.Scan() {
		acc := strings.Split(scanner.Text(), ":")
		account.Login = acc[0]
		account.Pass = acc[1]
		accs = append(accs, account)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return accs
}
