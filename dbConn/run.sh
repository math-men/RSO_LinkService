#!/bin/bash
IMAGE="amazon/dynamodb-local"
sudo docker ps -q --filter ancestor=$IMAGE | xargs -r sudo docker stop
export PORT=8000
sudo docker pull amazon/dynamodb-local
sudo docker-compose up &
COUNT=$(sudo docker ps | grep $IMAGE | wc -l)
while [ $COUNT -eq 0 ]
do
    COUNT=$(sudo docker ps | grep $IMAGE | wc -l)
    echo "Waiting for container to start"
    sleep 2
done
go get github.com/aws/aws-sdk-go/aws
go get github.com/aws/aws-sdk-go/aws/session
go build
go run prepareData.go
