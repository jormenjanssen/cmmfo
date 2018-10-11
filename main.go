package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	log.SetOutput(os.Stdout)

	err := DownloadFile("random-domains.txt", "https://github.com/jormenjanssen/cmmfo/raw/master/random-domains.txt")
	if err != nil {
		log.Fatalf("failed to download domain file: %v", err)
	}

	domains, err := GetDomains()
	if err != nil {
		log.Fatalf("failed to read/parse domains file: %v", err)
	}

	for {

		for _, url := range domains {
			_, err := net.LookupHost(url)

			if err != nil {
				log.Println(fmt.Sprintf("Cannot resolve host: %v error: %v", url, err))
			} else {
				log.Println(fmt.Sprintf("Resolved host: %v", url))
				go func() {
					// Visiting domain
					verr := VisitDomain(fmt.Sprintf("http://%v", url))
					if verr != nil {
						log.Println(fmt.Sprintf("Failed to get request(http://%v) %v", url, verr))
					} else {
						log.Println(fmt.Sprintf("Succesfull visited (http://%v)", url))
					}
				}()
			}

			time.Sleep(270 * time.Millisecond)
		}

		time.Sleep(6 * time.Minute)
	}
}

// GetDomains returns a long list of Domain from the domains.txt file
func GetDomains() ([]string, error) {

	content, err := ioutil.ReadFile("random-domains.txt")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	return lines, nil
}

// VisitDomain visits a http domain
func VisitDomain(url string) error {
	// Get the data
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	buffer := make([]byte, 4096)
	_, err = resp.Body.Read(buffer)
	return err
}
