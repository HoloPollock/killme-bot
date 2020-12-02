package handlers

import (
	"aws-lambda-golang/pkg/dyna"
	"aws-lambda-golang/pkg/person"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
)

type ErrorBody struct {
	Msg *string `json:"error"`
}

func GetPerson(req events.APIGatewayProxyRequest, d dyna.DynaTable) (*events.APIGatewayProxyResponse, error) {
	name := req.QueryStringParameters["name"]
	if len(name) <= 0 {
		return apiResponse(http.StatusBadRequest,ErrorBody{Msg: aws.String("no user requested")}) 
	}
	result, err := person.FetchPerson(name, d)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	
	return apiResponse(http.StatusOK, result)
}

func CreatePerson(req events.APIGatewayProxyRequest, d dyna.DynaTable) (*events.APIGatewayProxyResponse, error) {
	result, err := person.CreatePerson(req, d)
      if err != nil {
          return apiResponse(http.StatusBadRequest, ErrorBody{
              aws.String(err.Error()),
          })
      }
      return apiResponse(http.StatusCreated, result)
}

func UpdatePerson(req events.APIGatewayProxyRequest, d dyna.DynaTable) (*events.APIGatewayProxyResponse, error) {
	result, err := person.IncrementPerson(req, d)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
  }

func DeletePerson(req events.APIGatewayProxyRequest, d dyna.DynaTable) (*events.APIGatewayProxyResponse, error) {
	err := person.DeletePerson(req, d)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, "method Not allowed")
  }