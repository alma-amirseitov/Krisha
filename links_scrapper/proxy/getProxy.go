package proxy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"Krisha/links_scrapper/utils"
)

type data struct {
	Curl string `son:"curl"`
}


func GetProxy() (string,error){

	url := os.Getenv("PROXY")

	req,err:= http.NewRequest("GET", url, nil)

	if err != nil {
		utils.Logging("Error","Failed to make request in getting Proxy" + err.Error())
		return "",err
	}

	req.Header.Add("x-rapidapi-host", "proxy-orbit1.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "c0841e6526msh41636e7e727e5a6p171282jsn5cedcf6ca44e")

	res,err := http.DefaultClient.Do(req)
	if err != nil{
		utils.Logging("Error","Failed to do response in getting Proxy" + err.Error())
		return "",err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err != nil{
		utils.Logging("Error",err.Error())
		return "",err
	}

	var proxy data
	err = json.Unmarshal(body,&proxy)
	if err != nil{
		utils.Logging("Error","Failed to unmarshal json" + err.Error())
		return "",err
	}
	if proxy.Curl != ""{
		return proxy.Curl,nil
	}
	return "",nil
}