package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	u := "https://127.0.0.1:8123/aaaa"
	resp, err := c.Get(u)
	if err != nil {
		log.Printf("Get(%s) %s", u, err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll %s", err)
	}
	log.Printf("response data %s", b)
}
