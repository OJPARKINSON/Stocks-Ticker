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

type Response struct {
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

func Authentication(req *http.Request, params string) {
	apiSec := os.Getenv("apiSec")
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	apiKey := os.Getenv("apiKey")
	h := hmac.New(sha256.New, []byte(apiSec))
	message := timestamp + req.Method + params
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))

	req.Header.Add("CB-ACCESS-KEY", apiKey)
	req.Header.Add("CB-ACCESS-SIGN", signature)
	req.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("CB-VERSION", "2015-07-22")
}

func Request(url string, route string, auth bool) Response {
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", url+route, nil)
	if err != nil {
		log.Fatal(err)
	}

	if auth {
		Authentication(req, route)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	respJSON := Response{}
	json.Unmarshal(body, &respJSON)
	return respJSON
}

func handler(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	getEnv()
	token := os.Getenv("token")
	acountID := os.Getenv("acountID")
	options := "&types=quote,chart"

	var sellPrice, amount = 27.26, 330.0
	stringExchangeRate := Request("https://api.coinbase.com", "/v2/exchange-rates?currency=USD", true).Data.Rates.GBP
	exchangeRate, _ := strconv.ParseFloat(stringExchangeRate, 32)
	portfolio := Request("https://api.coinbase.com", "/v2/accounts/"+acountID, true).Data.Native_Balance.Amount
	cmcsa := Request("https://cloud.iexapis.com/stable/stock/cmcsa/batch?token=", token+options, false)
	UKcmcsaPrice := exchangeRate * cmcsa.Quote.LatestPrice
	UKcmcsaProft := UKcmcsaPrice*amount - sellPrice*amount

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("{Comcast Price: $%.2f, Portfolio: £%s, Comcast Profit: £%.2f}", cmcsa.Quote.LatestPrice, portfolio, UKcmcsaProft),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
