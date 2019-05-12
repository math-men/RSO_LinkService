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

func main() {
	fmt.Println("Starting configuration of data model")
	config := &aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:" + os.Getenv("PORT")),
	}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("year"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("year"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("title"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Movies"),
	}

	result, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
	fmt.Println("Successfully created dynamoDB data model")

	type Movie struct {
		Year  int         `json:"year"`
		Title string      `json:"title"`
	}

	moviesData, err := os.Open("moviedata.json")
	defer moviesData.Close()
	if err != nil {
		fmt.Println("Could not open the moviedata.json file", err.Error())
		os.Exit(1)
  }
	var movies []Movie
	err = json.NewDecoder(moviesData).Decode(&movies)
	if err != nil {
		fmt.Println("Could not decode the moviedata.json data", err.Error())
		os.Exit(1)
	}
	fmt.Println(movies)

	for _, movie := range movies {

		info, err := dynamodbattribute.MarshalMap(movie)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal the movie, %v", err))
		}

		input := &dynamodb.PutItemInput{
			Item:      info,
			TableName: aws.String("Movies"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

	fmt.Printf("We have processed %v records\n", len(movies))

}
