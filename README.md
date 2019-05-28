# Microservices for RSO project - Link Service

![GOGO](https://ih0.redbubble.net/image.520470450.9907/flat,550x550,075,f.u4.jpg)


Run dynamodb locally and create tables
cd dbConn
./run.sh

Run service
cd service
EXPORT AWS_ACCESS_KEY=<sth>
EXPORT AWS_SECRET_KEY=<sth>
EXPORT AWS_REGION=us-west-2

EXPORT HOST=http://localhost:
EXPORT PORT=<SERVER_PORT>
EXPORT REGION=us-west-2
EXPORT DYNAMO_PORT=8000

go get github.com/aws/aws-sdk-go/aws
go get github.com/aws/aws-sdk-go/service/dynamodb
go get github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute
go get github.com/go-chi/chi
go get github.com/go-chi/chi/middleware
go build
go run main.go


