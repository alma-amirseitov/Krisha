package links

import (
	client2 "Krisha/links_scrapper/client"
	"net/http"
	"path"
	"strconv"
	"strings"

	"Krisha/links_scrapper/rabbitMq"
	"Krisha/links_scrapper/utils"

	"github.com/PuerkitoBio/goquery"
)

// STARTING POINT - gets region's "link" and call GetCategories(link) function
func GetRegions(url string) error{
	utils.Logging("Info","Taking regions links from " + url)
	client := client2.GetClient()
	req,err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Logging("Error","Failed to make new Request - "+err.Error())
		return err
	}

	res, err := client.Do(req)
	if err != nil{
		utils.Logging("Error","problems to make response "+ err.Error())
		return err
	}

	status := res.StatusCode
	utils.Logging("Info","Successfully connected to "+url+", status code = "+ strconv.Itoa(status))

	defer res.Body.Close()

	doc,err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		utils.Logging("Error","Failed to create goquery document - "+err.Error())
		return err
	}

	doc.Find(".se-mainmenu div").Each(func(i int, selection *goquery.Selection) {
		link,exist := selection.Find("a").Attr("href")
		if exist {
			doesMatch,err := path.Match("/sitemap/*",path.Dir(link))
			if err!= nil{
				utils.Logging("Error",err.Error())
			}
			if doesMatch{
				for i:=0 ;i<10;i++{
					if err := GetCategories("https://krisha.kz"+link); err ==nil {
						break
					}
				}
			}
		}
	})
	utils.Logging("Debug","Advertisement's links has taken")
	return nil
}

// gets categories "link" for each region and call GetPageLinks(link) function
func GetCategories(url string)error{
	utils.Logging("Info","Taking categories links from " + url)
	client := client2.GetClient()

	req,err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Logging("Error","Failed to make new Request - "+err.Error())
		return err
	}

	res, err := client.Do(req)
	if err != nil{
		utils.Logging("Error","Failed to make new Response - "+err.Error())
		return err
	}
	defer res.Body.Close()

	doc,err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		utils.Logging("Error","Failed to create goquery document - "+err.Error())
		return err
	}
	doc.Find(".se-mainmenu div").Each(func(i int, selection *goquery.Selection){
		link,exist := selection.Find("a").Attr("href")
		if exist {
			if strings.Contains(link,"prodazha"){
				for i:=0 ;i<10;i++{
					if err := GetPageLinks("https://krisha.kz"+link); err ==nil {
						break
					}
				}
			}else if strings.Contains(link,"arenda"){
				for i:=0 ;i<10;i++{
					if err := GetPageLinks("https://krisha.kz"+link); err ==nil {
						break
					}
				}
			}
		}
	})
	utils.Logging("Debug","Advertisement links has taken from region " + url)
	return nil
}

// gets page "link" for each categories and call GetAddsLink(link) function
func GetPageLinks(url string) error{
	utils.Logging("Info","Taking Pages links from category " + url)
	client := client2.GetClient()
	req,err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Logging("Error","Failed to make new Request - "+err.Error())
		return err
	}

	res, err := client.Do(req)
	if err != nil{
		utils.Logging("Error","Error to make response "+err.Error())
		return err
	}
	defer res.Body.Close()

	doc,err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		utils.Logging("Error",""+err.Error())
		return err
	}

	pages := 0

	value := doc.Find(".a-search-subtitle.search-results-nb span").Text()

	value = strings.ReplaceAll(value,string(160),"")

	quantityAdd,err := strconv.Atoi(value)
	if err != nil{
		utils.Logging("Error","Error to convert string to int "+err.Error())
		return err
	}

	if quantityAdd <=20{
		pages = 1
	}else {
		pages = quantityAdd/20
		reminder := quantityAdd%20
		if reminder != 0 {pages++}
	}
	for i:=1;i<=pages;i++{
		if url[len(url)-1] == '/' {
			link := url + "?page=" + strconv.Itoa(i)
			for i:=0 ;i<10;i++{
				if err := GetAddsLink(link); err ==nil {
					break
				}
			}
		} else {
			link := url + "&page=" + strconv.Itoa(i)
			for i:=0 ;i<10;i++{
				if err := GetAddsLink(link); err ==nil {
					break
				}
			}
		}
	}
	utils.Logging("Debug","Advertisement links has taken from category " + url)
	return nil
}

// gets Advertisement "link" in each page and send to rabbitMq
func GetAddsLink(url string)  error{
	utils.Logging("Info","Taking Adds links from "+url)
	client := client2.GetClient()
	req,err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Logging("Error","Failed to make new Request - "+err.Error())
		return err
	}

	res, err := client.Do(req)
	if err != nil{
		utils.Logging("Error","Error to make response "+err.Error())
		return err
	}
	defer res.Body.Close()

	doc,err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		utils.Logging("Error",""+err.Error())
		return err
	}
	doc.Find("a.a-card__title ").Each(func(_ int, selection *goquery.Selection) {
		link,exist := selection.Attr("href")
		if exist{
			rabbitMq.Send("https://krisha.kz"+link)
		}
	})
	utils.Logging("Debug","Advertisement links has taken from page" + url)
	return nil
}
