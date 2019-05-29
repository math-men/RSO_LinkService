package link

import (
	"fmt"
	"context"
	"math/rand"
	"time"
	"github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	models "../../models"
	lRepo "../../repository"
)

var BASE_URL = "sshort.me/"

func NewDynamoLinkRepo(Conn *dynamodb.DynamoDB) lRepo.LinkRepo {
	return &dynamodbLinkRepo{
		Conn: Conn,
	}
}

type dynamodbLinkRepo struct {
	Conn *dynamodb.DynamoDB
}

func (d *dynamodbLinkRepo) Create(ctx context.Context, l *models.Link) (int64, error) {
  l.Processed = helperRandomLink()
	l.TTL = int64(time.Now().Unix() + l.TTL)
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
    return 1, err
  }

  fmt.Println("Successful insert")

  return 0, nil
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
	    return nil, err
	}

	links := []*models.Link{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		return nil, err
	}
	if len(links) == 0{
		return nil, models.ErrNotFound
	}
	for _, link := range links {
		link.Processed = BASE_URL + link.Processed
	}
	return links, nil

}

func (d *dynamodbLinkRepo) Fetch(ctx context.Context) ([]*models.Link, error) {

		params := &dynamodb.ScanInput{
			TableName: aws.String("Links"),
		}
		result, err := d.Conn.Scan(params)
		fmt.Println(result)
		fmt.Println(err)
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

func (d *dynamodbLinkRepo) RegisterClick(ctx context.Context, c *models.Click) (error) {

	info, err := dynamodbattribute.MarshalMap(c)
  if err != nil {
    panic(fmt.Sprintf("failed to marshal the link, %v", err))
  }

  input := &dynamodb.PutItemInput{
    Item:      info,
    TableName: aws.String("Clicks"),
  }

  _, err = d.Conn.PutItem(input)
  if err != nil {
    fmt.Println(err.Error())
    return err
  }

  fmt.Println("Successful insert")

  return nil

}

func (d *dynamodbLinkRepo) GetClicks(ctx context.Context, owner string) (int, error) {
	var queryInput = &dynamodb.QueryInput{
	    TableName: aws.String("Clicks"),
	    KeyConditions: map[string]*dynamodb.Condition{
	        "owner": {
	            ComparisonOperator: aws.String("EQ"),
	            AttributeValueList: []*dynamodb.AttributeValue{
	                {
	                    S: aws.String(owner),
	                },
	            },
	        },
	    },
	}

	var result, err = d.Conn.Query(queryInput)
	if err != nil {
	    return -1, err
	}

	links := []*models.Link{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		return -1, err
	}
	return len(links), nil
}

func helperRandomLink() string {
	 len := 9
   bytes := make([]byte, len)
   for i := 0; i < len; i++ {
   	   bytes[i] = byte(65 + rand.Intn(25))  //A=65 and Z = 65+25
   }
   return string(bytes)

}
