package main

import (
    "context"
    "fmt"

    "github.com/aws/aws-lambda-go/lambda"
)

type SampleEvent struct {
    ID   string `json:"id"`
    Val  int    `json:"val"`
    Flag bool   `json:"flag"`
}

func HandleRequest(ctx context.Context, event SampleEvent) (string, error) {
    return fmt.Sprintf("%+v", event), nil
}

func main() {
    lambda.Start(HandleRequest)
}
