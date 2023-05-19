package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func main() {
	lambda.Start(handler)
}

func getStockPrice() float64 {
	alpacaClient := marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    os.Getenv("alpacaKey"),
		ApiSecret: os.Getenv("alpacaSecret"),
	})

	alpaca, err := alpacaClient.GetLatestTrade("CMCSA")

	if err != nil {
		fmt.Printf("Failed to get CMCSA price: %v\n", err)
	}

	return alpaca.Price
}

func getExchangeRate() string {
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", "https://api.coinbase.com/v2/exchange-rates?currency=USD", nil)

	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	JSON := ExchangeRateResponse{}
	json.Unmarshal(body, &JSON)

	return JSON.Data.Rates.GBP
}

func handler(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	getEnv()
	var sellPrice, amountOfStock = 27.26, 330.0
	stringExchangeRate := getExchangeRate()
	exchangeRate, _ := strconv.ParseFloat(stringExchangeRate, 64)

	cmcsaPrice := getStockPrice()

	gpbCMCSAPrice := exchangeRate * cmcsaPrice
	UKcmcsaProfit := gpbCMCSAPrice*amountOfStock - sellPrice*amountOfStock

	fmt.Printf("{Comcast Price: $%.2f, Comcast Profit: £%.2f}", cmcsaPrice, UKcmcsaProfit)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("{\"Comcast_Price\": \"$%.2f\", \"Comcast_Profit\": \"£%.2f\"}", cmcsaPrice, UKcmcsaProfit),
		StatusCode: 200,
	}, nil
}

type Response struct {
	Data struct {
		Currency       string `json:"currency"`
		Amount         string `json:"amount"`
		Native_Balance struct {
			Amount string `json:"amount"`
		} `json:"native_balance"`
	} `json:"data"`
	Quote struct {
		LatestPrice float64 `json:"latestPrice"`
	} `json:"quote"`
}

type ExchangeRateResponse struct {
	Data struct {
		Rates struct {
			GBP string `json:"GBP"`
		} `json:"rates"`
	} `json:"data"`
}

func getEnv() {
	env := os.Getenv("NODE_ENV")
	if env != "Production" {
		var err = godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}
}

// func Authentication(req *http.Request, params string) {
// 	apiKey := os.Getenv("apiKey")
// 	timestamp := fmt.Sprintf("%v", time.Now().Unix())

// 	now := strconv.FormatInt(time.Now().Unix(), 10)
// 	sign := createAccessSign(now, "GET", params, "")

// 	req.Header.Add("CB-ACCESS-KEY", apiKey)
// 	req.Header.Add("CB-ACCESS-SIGN", sign)
// 	req.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
// 	req.Header.Add("CB-VERSION", "2017-10-07")
// }

// portfolio := Request("https://api.coinbase.com", "/v2/accounts/"+ os.Getenv("accountID"), true).Data.Native_Balance.Amount
// fmt.Printf("{Comcast Price: $%.2f, Portfolio: £%s, Comcast Profit: £%.2f}", cmcsaPrice, portfolio, UKcmcsaProfit)

// func createAccessSign(timestamp, method, requestPath, body string) string {

// 	key, err := base64.StdEncoding.DecodeString(os.Getenv("apiSec"))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	mac := hmac.New(sha256.New, key)
// 	mac.Write([]byte(timestamp + method + requestPath + body))
// 	sign := mac.Sum(nil)
// 	signBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(sign)))
// 	base64.StdEncoding.Encode(signBase64, sign)
// 	return string(signBase64)
// }
