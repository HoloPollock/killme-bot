package dyna

import "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

type DynaTable struct {
	DynaClient dynamodbiface.DynamoDBAPI
	TableName string
}

func NewDynaTable(client dynamodbiface.DynamoDBAPI, name string) DynaTable {
	return DynaTable{client, name}
}