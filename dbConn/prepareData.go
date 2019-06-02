package main

import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Link struct {
	Original  string      `json:"original"`
	Processed string			`json:"processed"`
	Ttl      int			    `json:"ttl"`
}

func createTableLinks(svc *dynamodb.DynamoDB) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("processed"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("original"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("processed"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("original"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Links"),
	}

	createTable(svc, input)
}

func createTableClicks(svc *dynamodb.DynamoDB) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("processed"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("timestamp"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("processed"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("timestamp"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Clicks"),
	}

	createTable(svc, input)
}

func createTable(svc *dynamodb.DynamoDB, input *dynamodb.CreateTableInput) {
	result, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
	fmt.Println("Successfully created dynamoDB data model")
}

func insertData(svc *dynamodb.DynamoDB) {
	linksData, err := os.Open("linkdata.json")
	defer linksData.Close()
	if err != nil {
		fmt.Println("Could not open the moviedata.json file", err.Error())
		os.Exit(1)
  }
	var links []Link
	err = json.NewDecoder(linksData).Decode(&links)
	if err != nil {
		fmt.Println("Could not decode the linksData.json data", err.Error())
		os.Exit(1)
	}
	fmt.Println(links)

	for _, link := range links {

		info, err := dynamodbattribute.MarshalMap(link)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal the link, %v", err))
		}

		input := &dynamodb.PutItemInput{
			Item:      info,
			TableName: aws.String("Links"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

	fmt.Printf("We have processed %v records\n", len(links))
}

func queryData(svc *dynamodb.DynamoDB) {
	params := &dynamodb.ScanInput{
		TableName: aws.String("Links"),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Failed to query")
		return
	}

	links := []Link{}

	// Unmarshal the Items field in the result value to the Item Go type.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		fmt.Println("Failed to unmarshall query")
		return
	}

	// Print out the items returned
	fmt.Println("Query results:")
	for _, link := range links {
		fmt.Println(link)
	}
}

func main() {
	fmt.Println("Starting configuration of data model")
	config := &aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:" + os.Getenv("DYNAMO_PORT")),
	}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)

	createTableLinks(svc)
	createTableClicks(svc)
	insertData(svc)
	queryData(svc)

}
