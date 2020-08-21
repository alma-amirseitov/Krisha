package main

import (
	"Krisha/links_scrapper/links"
	"Krisha/links_scrapper/utils"
	"os"
)

func main() {
	utils.Logging("Info","Starting to Crawl...")
	val := os.Getenv("KRISHA_URL")
	for i:=0 ;i<10;i++{
		if err := links.GetRegions(val); err ==nil {
			break
		}
	}
}

