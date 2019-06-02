
package driver

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DB ...
type DB struct {
	SQL *sql.DB
	Dynamko *dynamodb.DynamoDB
}

// DBConn ...
var dbConn = &DB{}

func ConnectDynamo(host, port, region string) (*DB) {

	var config *aws.Config
	if(os.Getenv("ISLOCAL") == "true") {
		config = &aws.Config{
			Region:   aws.String(region),
			Endpoint: aws.String(host + port),
		}
	} else {
		config = &aws.Config{
			Region:   aws.String(region),
		}
	}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)
	dbConn.Dynamko = svc
	return dbConn
}

// ConnectSQL ...
func ConnectSQL(host, port, uname, pass, dbname string) (*DB, error) {
	dbSource := fmt.Sprintf(
		"root:%s@tcp(%s:%s)/%s?charset=utf8",
		pass,
		host,
		port,
		dbname,
	)
	d, err := sql.Open("mysql", dbSource)
	if err != nil {
		panic(err)
	}
	dbConn.SQL = d
	return dbConn, err
}
