package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

func main() {

	domains, _ := getDomains()

	for {

		for _, url := range domains {
			_, err := net.LookupHost(url)

			if err != nil {
				fmt.Println(fmt.Sprintf("Cannot resolve host: %v error: %v", url, err))
			} else {
				fmt.Println(fmt.Sprintf("Resolved host: %v", url))
			}

			time.Sleep(170 * time.Millisecond)
		}

		time.Sleep(7 * time.Minute)
	}
}

func getDomains() ([]string, error) {

	content, err := ioutil.ReadFile("random-domains.txt")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	return lines, nil
}
