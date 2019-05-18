package link

import (
	"fmt"
	"context"
	"github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	models "../../models"
	lRepo "../../repository"
)

func NewDynamoLinkRepo(Conn *dynamodb.DynamoDB) lRepo.LinkRepo {
	return &dynamodbLinkRepo{
		Conn: Conn,
	}
}

type dynamodbLinkRepo struct {
	Conn *dynamodb.DynamoDB
}


func (d *dynamodbLinkRepo) Create(ctx context.Context, l *models.Link) (int64, error) {


  info, err := dynamodbattribute.MarshalMap(l)
  if err != nil {
    panic(fmt.Sprintf("failed to marshal the link, %v", err))
  }

  input := &dynamodb.PutItemInput{
    Item:      info,
    TableName: aws.String("Links"),
  }

  _, err = d.Conn.PutItem(input)
  if err != nil {
    fmt.Println(err.Error())
    return 1,nil
  }

  fmt.Println("Successful insert")

  return 123,nil
}

func (d *dynamodbLinkRepo) Fetch(ctx context.Context) ([]*models.Link, error) {

	params := &dynamodb.ScanInput{
			TableName: aws.String("Links"),
		}
		result, err := d.Conn.Scan(params)
		if err != nil {
			fmt.Println("Failed to query")
			return nil, err
		}

		links := []*models.Link{}

		// Unmarshal the Items field in the result value to the Item Go type.
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &links)
		if err != nil {
			fmt.Println("Failed to unmarshall query")
			return nil, err
		}

		// Print out the items returned
		fmt.Println("Query results:")
		for _, link := range links {
			fmt.Println(link)
	}

	return links, nil
}
