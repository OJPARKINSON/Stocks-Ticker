package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/fabioberger/coinbase-go"
	"github.com/joho/godotenv"
)

func getEnv() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func cmcsa() ([]string, error) {
	getEnv()
	var priceAndProfit []string
	GBPAverage := os.Getenv("averageGBP")
	sellPrice, ammount := 27.26, 330.0

	doc, err := goquery.NewDocument("https://in.finance.yahoo.com/lookup?s=CMCSA")

	if err != nil {
		log.Fatal(err)
		return priceAndProfit, err
	}

	c := coinbase.ApiKeyClient(os.Getenv("COINBASE_KEY"), os.Getenv("COINBASE_SECRET"))

	balance, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance is %f BTC", balance)

	doc.Find("td[data-reactid='59']").Each(func(index int, item *goquery.Selection) {
		currentPrice, _ := strconv.ParseFloat(item.Text(), 64)
		exchangeRate, _ := strconv.ParseFloat(GBPAverage, 64)
		GBPCMCSA := currentPrice * exchangeRate
		CMCSAProfit := (GBPCMCSA*ammount - sellPrice*ammount)

		priceAndProfit = []string{strconv.FormatFloat(GBPCMCSA, 'f', 2, 64), strconv.FormatFloat(CMCSAProfit, 'f', 2, 64)}
	})

	return priceAndProfit, nil
}

func main() {
	CMCSAPrice, err := cmcsa()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("£" + CMCSAPrice[0])
	fmt.Println("£" + CMCSAPrice[1])
}
