package ads

import (
	client2 "Krisha/advertisement_scrapper/client"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"Krisha/advertisement_scrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

type Advertisement struct {
	Link string `json:"link"`
	City string `json:"city"`
	Street string `json:"street"`
	Title string `json:"title"`
	Sum string `json:"sum"`
	Phone []string `json:"phone"`
	Year string `json:"year"`
	Area string `json:"area"`
	Count_floors int `json:"count_floors"`
	Count_rooms int `json:"count_rooms"`
	House_number int `json:"house_number"`
	Characteristic string `json:"characteristic"`
	Pledget bool `json:"pledget"`
}

// get value from otto object
func getValueFromObject(val otto.Value, key string) (*otto.Value, error) {
	if !val.IsObject() {
		return nil, errors.New("passed val is not an Object -- "+key)
	}

	valObj := val.Object()

	obj, err := valObj.Get(key)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}


func takeData(doc *goquery.Document) *otto.Value{
	data := strings.TrimSpace(doc.Find("#jsdata").Text())
	vm := otto.New()

	_,err := vm.Run(data)
	if err != nil{
		utils.Logging("Error",err.Error())
		return nil
	}

	val,err := vm.Get("data")
	if err != nil{
		utils.Logging("Error",err.Error())
		return nil
	}

	pdata,err := getValueFromObject(val,"advert")
	if err != nil{
		utils.Logging("Error",err.Error())
		return nil
	}
	return pdata
}

func getStreet(data *otto.Value)string{
	address,err := getValueFromObject(*data,"address")
	if err != nil{
		utils.Logging("Error",err.Error())
		return ""
	}

	result,err := getValueFromObject(*address,"street")
	if err != nil{
		utils.Logging("Error",err.Error())
		return ""
	}
	street := fmt.Sprintf("%s",result)

	return street
}

func getCity(data *otto.Value)string{

	adress,err := getValueFromObject(*data,"address")
	if err != nil{
		utils.Logging("Error",err.Error())
		return ""
	}

	result,err := getValueFromObject(*adress,"city")
	if err != nil{
		utils.Logging("Error",err.Error())
		return ""
	}
	city := fmt.Sprintf("%s",result)

	return city
}

func getTitle(data *otto.Value)string{

	result,err := getValueFromObject(*data,"title")
	if err != nil{
		utils.Logging("Error",err.Error())
		return ""
	}
	title := fmt.Sprintf("%s",result)

	return title
}

func getPrice(data *otto.Value)string{

	result,err := getValueFromObject(*data,"price")
	if err != nil{
		utils.Logging("Error",err.Error())
		return ""
	}
	title := fmt.Sprintf("%s",result)

	return title
}

func getPhone(data *otto.Value,url string,client http.Client)[]string{
	var phones []string
	result,err :=  getValueFromObject(*data,"id")
	if err != nil{
		utils.Logging("Error",err.Error())
		return  phones
	}
	id := fmt.Sprintf("%s",result)

	req, err := http.NewRequest("GET", "https://krisha.kz/a/ajaxPhones?id="+id, nil)
	if err != nil {
		utils.Logging("Error",err.Error())
		return phones
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://krisha.kz/a/show/"+id)
	req.Header.Set("Cookie",
		"krssid=t6baa17ssjg8lrn7a2tm0h1902; " +
			"krishauid=d6a7774aa6aceca614fb7e9779f1a70a589a6130; " +
			"_gcl_au=1.1.530140204.1596613086;" +
			" _ga=GA1.2.625127310.1596613091; " +
			"_gid=GA1.2.1130648494.1596613091; " +
			"_ym_uid=1596613091132319815;" +
			" _ym_d=1596613091; " +
			"_ym_visorc_10575199=w; " +
			"_ym_isad=2;" +
			" _fbp=fb.1.1596613092837.1410708002; " +
			"__gads=ID=2a83db5ff6e98976-22ad1560abb600d1:" +
			"T=1596613092:S=ALNI_MZLHWSuoaSNl3iu0nVDx5tMMm7BJQ; " +
			"saw_region_hint=1; " +
			"_gat=1; " +
			"ssaid=0b136980-d6f3-11ea-a7c9-f52242ec781e; " +
			"__tld__=null; " +
			"tutorial=%7B%22advPage%22%3A%22viewed%22%7D; " +
			"kr_cdn_host=//alakcell-kz.kcdn.online; " +
			"cto_bundle=Zklvy19XVHJDNU41aVJDc2pVYVg2Y1NwYVJxc24wbjRHcXJublF3SEElMkZEVHpsMmwlMkZlYiUyQmJxOVN2aTM1alRSb1RwWk5JN2p0c2hsV1U4WEtscXElMkIlMkZwZlNUcUVTeG1ad2tweG5xcEtSVWFaUHJmRUtBZ2lNWXgzMmUxJTJCT1lDZTRETWRGZQ; " +
			"_ym_visorc_49456573=w; " +
			"_gat_UA-20095530-1=1")
	resp, err := client.Do(req)
	if err != nil {
		utils.Logging("Error",err.Error())
		return phones
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Logging("Error",err.Error())
		return phones
	}
	defer resp.Body.Close()

	value := strings.ReplaceAll(string(bodyText),"\"","")
	value = strings.ReplaceAll(value,"[","")
	value = strings.ReplaceAll(value,"]","")
	value = strings.ReplaceAll(value," ","")
	value = strings.ReplaceAll(value,"+","")
	phones = strings.Split(value,",")


	return phones
}

func getYear(doc *goquery.Document) string{
	var yearInt string
	doc.Find(".offer__info-item").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find(".offer__info-title").Text()
		doesMatch := strings.Contains(title, "Дом")
		if doesMatch {
			year := selection.Find(".offer__advert-short-info").Text()
			yearInt = year
		}
	})
	return yearInt
}

func getArea(data *otto.Value) string{
	result,err :=  getValueFromObject(*data,"square")
	if err != nil{
		utils.Logging("Error",err.Error())
		return ""
	}
	area := fmt.Sprintf("%s",result)
	return area
}

func getFloorNumber(doc *goquery.Document)int{
	return 0
}

func getRoomNumber(data *otto.Value)int{
	result,err :=  getValueFromObject(*data,"rooms")
	if err != nil{
		utils.Logging("Error",err.Error())
		return 0
	}
	count_rooms,_ := result.ToInteger()
	return int(count_rooms)
}

func getHouseNumber(data *otto.Value)int{

	adress,err := getValueFromObject(*data,"address")
	if err != nil{
		utils.Logging("Error",err.Error())
		return 0
	}

	result,err := getValueFromObject(*adress,"house_num")
	if err != nil{
		utils.Logging("Error",err.Error())
		return 0
	}
	house_number,_ := result.ToInteger()
	return int(house_number)
}

func getCharacteristic(doc *goquery.Document)string{
	characteristic := strings.TrimSpace(doc.Find(".text").Text())
	characteristic = strings.ReplaceAll(characteristic,"\n","")
	return characteristic
}

func isPledget(doc *goquery.Document)bool{
	text := doc.Find(".offer__description").Text()
	pledged := strings.Contains(text,"залоге")
	return pledged
}


func GetAdd(url string) (string,error){
	client := client2.GetClient()
	req,err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Logging("Error","Failed to make new Request - "+err.Error())
	}

	res, err := client.Do(req)
	if err != nil{
		utils.Logging("Error",err.Error())
		return "",err
	}
	defer res.Body.Close()

	status := res.StatusCode
	utils.Logging("Info","Status code -- "+ strconv.Itoa(status))

	doc,err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		utils.Logging("Error",err.Error())
		return "",err
	}

	data := takeData(doc)
	link := url
	city := getCity(data)
	street := getStreet(data)
	title := getTitle(data)
	sum := getPrice(data)
	phone := getPhone(data,url,*client)
	year := getYear(doc)
	area := getArea(data)
	count_floors := getFloorNumber(doc)
	count_rooms := getRoomNumber(data)
	house_number := getHouseNumber(data)
	characteristic := getCharacteristic(doc)
	pledget := isPledget(doc)

	advertisement := Advertisement{
		Link:link,
		City:city,
		Street:street,
		Title:title,
		Sum:sum,
		Phone:phone,
		Year:year,
		Area:area,
		Count_floors:count_floors,
		Count_rooms:count_rooms,
		House_number:house_number,
		Characteristic:characteristic,
		Pledget:pledget,
	}

	jsonData,err := json.Marshal(advertisement)
	if err != nil {
		utils.Logging("Error","Failed to Marshal json - " +err.Error())
		return "",err
	}
	return string(jsonData),nil
}


