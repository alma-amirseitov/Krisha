package elastic

import (
	"context"
	"fmt"
	"os"

	"Krisha/saver/utils"

	"github.com/olivere/elastic"
)
//Creating ESClient
func GetESClient() (*elastic.Client,error){
	val := os.Getenv("KRISHA_ELASTIC")
	client,err := elastic.NewClient(elastic.SetURL(val),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	fmt.Println("ES initialized...")

	return client,err
}

//Send to Elastic
func ToElastic(js string) error{
	ctx := context.Background()
	esClient, err := GetESClient()
	if err != nil {
		utils.Logging("Error","Failed initializing : "+ err.Error() + "Client fail ")
		return err
	}

	ind,err := esClient.Index().
		Index("krishaCrawler").
		BodyJson(js).
		Do(ctx)

	if err != nil {
		utils.Logging("Error", err.Error())
		fmt.Println(ind)
		return err
	}
	utils.Logging("Debug","[Elastic][InsertProduct]Insertion Successful")
	return nil
}