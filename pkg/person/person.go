package person

import (
	"aws-lambda-golang/pkg/dyna"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type KillMe struct {
	Name string `json:"name"`
	Times int `json:"times"`
}

func (k *KillMe) increment() {
	k.Times++
}

type Name struct {
	Name string `json:"name"`
}

var (
    ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
    ErrorFailedToFetchRecord     = "failed to fetch record"
    ErrorInvalidPersonData         = "invalid person data"
    ErrorInvalidName           = "invalid name"
    ErrorCouldNotMarshalItem     = "could not marshal item"
    ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item error"
	ErrorCouldNotDynamoUpdateItem   = "could not dynamo update item error"
    ErrorPersonAlreadyExists       = "user.User already exists"
    ErrorPersonDoesNotExists       = "user.User does not exist"
  )

func FetchPerson(name string, table dyna.DynaTable) (*KillMe, error) {
	input := &dynamodb.GetItemInput{
        Key: map[string]*dynamodb.AttributeValue{
            "name": {
                S: aws.String(name),
            },
        },
        TableName: aws.String(table.TableName),
	}
	
	result, err := table.DynaClient.GetItem(input)
    if err != nil {
        return nil, errors.New(ErrorFailedToFetchRecord)
  
	}
	person := new(KillMe)
    err = dynamodbattribute.UnmarshalMap(result.Item, person)
    if err != nil {
        return nil, errors.New(ErrorFailedToUnmarshalRecord)
    }
    return person, nil
}

func CreatePerson(req events.APIGatewayProxyRequest, table dyna.DynaTable) (*KillMe, error) {
	var p KillMe
	if err := json.Unmarshal([]byte(req.Body), &p); err != nil {
		return nil, errors.New(ErrorInvalidPersonData)
	}
	currentPerson, _ := FetchPerson(p.Name, table)
	if currentPerson != nil && len(currentPerson.Name) != 0 {
		return nil , errors.New(ErrorPersonAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(table.TableName),
	}

	_, err = table.DynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &p, nil

}

func IncrementPerson(req events.APIGatewayProxyRequest, table dyna.DynaTable) (
    *KillMe,
    error,
  ) {
    var n Name
    if err := json.Unmarshal([]byte(req.Body), &n); err != nil {
        return nil, errors.New(ErrorInvalidName)
    }
  
	key, err := dynamodbattribute.MarshalMap(n)
    // Save user  
    input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{ ":inc": {N: aws.String("1")} },
		TableName: aws.String(table.TableName),
		Key: key,
		ReturnValues:     aws.String("UPDATED_NEW"),
    	UpdateExpression: aws.String("ADD times :inc"),
	}
	
	result, err := table.DynaClient.UpdateItem(input)
	if err != nil {
        return nil, err
	}
	log.Println(result)
	times, _ := strconv.Atoi(aws.StringValue(result.Attributes["times"].N))
	newChange := KillMe{Name: n.Name, Times: times}
    return &newChange, nil
}

func DeletePerson(req events.APIGatewayProxyRequest, table dyna.DynaTable) error {
	name := req.QueryStringParameters["name"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(name),
			},
		},
		TableName: aws.String(table.TableName),
	}
	_, err := table.DynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}

