package httputils

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Get(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("HTTP GET %s ERROR: %s\n", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("HTTP GET %s ERROR: %s\n", url, err)
	}
	return body
}
