package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"energyjournal/internal/server"
)

// lambdaHandler adapts the existing HTTP server to the AWS Lambda invocation model.
type lambdaHandler struct {
	adapter *httpadapter.HandlerAdapter
}

func newLambdaHandler() *lambdaHandler {
	srv := server.New(":0")

	return &lambdaHandler{
		adapter: httpadapter.New(srv.Handler),
	}
}

func (h *lambdaHandler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return h.adapter.ProxyWithContext(ctx, req)
}

func main() {
	handler := newLambdaHandler()
	log.Printf("Lambda HTTP adapter initialized")
	lambda.Start(handler.Handle)
}
