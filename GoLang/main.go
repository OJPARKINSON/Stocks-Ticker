package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/joho/godotenv"
)

type CoinbaseResp struct {
	Data struct {
		Currency       string `json:"currency"`
		Amount         string `json:"amount"`
		Native_Balance struct {
			Amount string `json:"amount"`
		} `json:"native_balance"`
		Rates struct {
			GBP string `json:"GBP"`
		} `json:"rates"`
	} `json:"data"`
}

type iexapiResp struct {
	Quote struct {
		LatestPrice float64 `json:"latestPrice"`
	} `json:"quote"`
}

func getEnv() {
	env := os.Getenv("NODE_ENV")
	if env != "Production" {
		var err = godotenv.Load("../.env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}
}

func coinbaseRequest(url string) CoinbaseResp {
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	apiKey := os.Getenv("apiKey")
	apiSec := os.Getenv("apiSec")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "https://api.coinbase.com"+url, nil)
	if err != nil {
		log.Fatal(err)
	}

	h := hmac.New(sha256.New, []byte(apiSec))
	message := timestamp + req.Method + url
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))

	req.Header.Add("CB-ACCESS-KEY", apiKey)
	req.Header.Add("CB-ACCESS-SIGN", signature)
	req.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("CB-VERSION", "2015-07-22")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	respJSON := CoinbaseResp{}
	json.Unmarshal([]byte(body), &respJSON)

	return respJSON
}

func iexapi() float64 {
	token := os.Getenv("token")
	options := "&types=quote,chart"
	resp, err := http.Get("https://cloud.iexapis.com/stable/stock/cmcsa/batch?token=" + token + options)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	respJSO := iexapiResp{}
	json.Unmarshal([]byte(body), &respJSO)

	return respJSO.Quote.LatestPrice
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	getEnv()

	var sellPrice, amount float64 = 27.26, 330.0
	stringExchangeRate := coinbaseRequest("/v2/exchange-rates?currency=USD").Data.Rates.GBP
	exchangeRate, _ := strconv.ParseFloat(stringExchangeRate, 32)
	portfolio := coinbaseRequest("/v2/accounts/" + os.Getenv("acountID"))
	cmcsa := iexapi()
	UKcmcsa := exchangeRate * cmcsa

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Comcast Price: $%.2f \n 📈 Portfolio: £%s \n 📈 Comcast Profit: £%.2f \n", cmcsa, portfolio.Data.Native_Balance.Amount, UKcmcsa*amount-sellPrice*amount),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
