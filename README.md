# Microservices for RSO project - Link Service

![GOGO](https://ih0.redbubble.net/image.520470450.9907/flat,550x550,075,f.u4.jpg)


Run dynamodb locally and create tables<br/>
cd dbConn <br/>
./run.sh <br/>

Run service <br/>
cd service <br/>
EXPORT AWS_ACCESS_KEY=<sth> <br/>
EXPORT AWS_SECRET_KEY=<sth> <br/>
EXPORT AWS_REGION=us-west-2 <br/>

EXPORT HOST=http://localhost: <br/>
EXPORT PORT=<SERVER_PORT> <br/>
EXPORT REGION=us-west-2 <br/>
EXPORT DYNAMO_PORT=8000 <br/>

go get github.com/aws/aws-sdk-go/aws <br/>
go get github.com/aws/aws-sdk-go/service/dynamodb <br/>
go get github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute <br/>
go get github.com/go-chi/chi <br/> 
go get github.com/go-chi/chi/middleware <br/>
go build <br/>
go run main.go <br/>


