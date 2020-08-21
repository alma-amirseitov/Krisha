package client

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"Krisha/advertisement_scrapper/proxy"
	"Krisha/advertisement_scrapper/utils"

)


func GetClient() *http.Client{
	utils.Logging("Info","Creating a new client")
	randomProxy,err := proxy.GetProxy()
	if randomProxy == "" || err !=nil{
		utils.Logging("Error","Received empty proxy "+err.Error())
		return nil
	}
	proxyUrl,err := url.Parse(randomProxy)

	if err != nil{
		utils.Logging("error","Failed to parse randomProxy"+err.Error())
		return nil
	}

	transport := &http.Transport{
		Proxy:http.ProxyURL(proxyUrl),
		DialContext:(&net.Dialer{
			Timeout:       30 * time.Second,
			KeepAlive:     30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		IdleConnTimeout: 90 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport:    transport,
		Timeout:       time.Second * 30,
	}
	return client
}