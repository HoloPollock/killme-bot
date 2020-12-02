package main

import (
	"aws-lambda-golang/pkg/dyna"
	"aws-lambda-golang/pkg/handlers"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)


var (
	dynatable dyna.DynaTable
)

const TABLENAME = "KillMe"

func main() {
    region := os.Getenv("AWS_REGION")
    awsSession, err := session.NewSession(&aws.Config{
        Region: aws.String(region)},
    )
    if err != nil {
        return
    }
	dynaClient := dynamodb.New(awsSession)
	dynatable = dyna.NewDynaTable(dynaClient, TABLENAME)
    lambda.Start(handler)
 }

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetPerson(req, dynatable)
	case "POST":
		return handlers.CreatePerson(req, dynatable)
	case "PUT":
		return handlers.UpdatePerson(req, dynatable)
	case "DELETE":
		return handlers.DeletePerson(req, dynatable)
	default:
		return handlers.UnhandledMethod()
	}
}

