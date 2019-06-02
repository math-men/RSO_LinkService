package link

import (
	"fmt"
	"context"
	"math/rand"
	"time"
	"strings"
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

func (d *dynamodbLinkRepo) Create(ctx context.Context, l *models.Link) (string, error) {
  l.Processed = helperRandomLink()
	l.TTL = int64(time.Now().Unix() + l.TTL)
  info, err := dynamodbattribute.MarshalMap(l)
  if err != nil {
    return "", models.ErrMarshalling
  }

  input := &dynamodb.PutItemInput{
    Item:      info,
    TableName: aws.String("Links"),
  }

  _, err = d.Conn.PutItem(input)
  if err != nil {
		fmt.Println(err)

    return "", models.ErrInsert
  }


  return l.Processed, nil
}

func (d *dynamodbLinkRepo) Get(ctx context.Context, shortened string) ([]*models.Link, error) {
	fmt.Println(shortened)
	var queryInput = &dynamodb.QueryInput{
	    TableName: aws.String("Links"),
	    KeyConditions: map[string]*dynamodb.Condition{
	        "processed": {
	            ComparisonOperator: aws.String("EQ"),
	            AttributeValueList: []*dynamodb.AttributeValue{
	                {
	                    S: aws.String(shortened),
	                },
	            },
	        },
	    },
	}
	var result, err = d.Conn.Query(queryInput)
	if err != nil {
			fmt.Println(err)
	    return nil, models.ErrQuery
	}

	links := []*models.Link{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		return nil, models.ErrMarshalling
	}
	if len(links) == 0{
		return nil, models.ErrNotFound
	}

	return links, nil

}

func (d *dynamodbLinkRepo) Fetch(ctx context.Context) ([]*models.Link, error) {

		params := &dynamodb.ScanInput{
			TableName: aws.String("Links"),
		}
		fmt.Println("query")

		result, err := d.Conn.Scan(params)
		if err != nil {
			fmt.Println(err)
			return nil, models.ErrQuery
		}

		links := []*models.Link{}

		// Unmarshal the Items field in the result value to the Item Go type.
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &links)
		if err != nil {
			return nil, models.ErrMarshalling
		}

		return links, nil
}

func (d *dynamodbLinkRepo) RegisterClick(ctx context.Context, c *models.Click) (error) {

	info, err := dynamodbattribute.MarshalMap(c)
  if err != nil { //A=65 and Z = 65+25
    return models.ErrMarshalling
  }

  input := &dynamodb.PutItemInput{
    Item:      info,
    TableName: aws.String("Clicks"),
  }

  _, err = d.Conn.PutItem(input)
  if err != nil {
    return models.ErrInsert
  }

  return nil

}

func (d *dynamodbLinkRepo) GetClicks(ctx context.Context, shortened string) (int, error) {
	var queryInput = &dynamodb.QueryInput{
	    TableName: aws.String("Clicks"),
	    KeyConditions: map[string]*dynamodb.Condition{
	        "processed": {
	            ComparisonOperator: aws.String("EQ"),
	            AttributeValueList: []*dynamodb.AttributeValue{
	                {
	                    S: aws.String(shortened),
	                },
	            },
	        },
	    },
	}

	var result, err = d.Conn.Query(queryInput)
	if err != nil {
	    return -1, models.ErrQuery
	}
	links := []*models.Link{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		return -1, models.ErrMarshalling
	}
 	return len(links), nil
}

func validUrl(url string) bool {
	if strings.Contains(url, "http") {
		return true
	}
	return false;
}

func helperRandomLink() string {
	 len := 9
   bytes := make([]byte, len)
   for i := 0; i < len; i++ {
   	   bytes[i] = byte(65 + rand.Intn(25))
   }
   return string(bytes)

}
