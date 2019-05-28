
FROM golang:1.11-alpine
ENV SOURCES /go/src/
ENV HOST=http://localhost:
ENV PORT=9001
ENV DYNAMO_PORT=8000
ENV REGION=us-west-2
ENV AWS_ACCESS_KEY=test
ENV AWS_SECRET_KEY=test
ENV AWS_REGION=us-west-2

COPY . ${SOURCES}
WORKDIR ${SOURCES}
RUN apk add --no-cache git && \
    go get github.com/go-chi/chi && \
    go get github.com/go-chi/chi/middleware && \
    go get github.com/aws/aws-sdk-go/aws && \
    go get github.com/aws/aws-sdk-go/service/dynamodb  && \
    go get github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute && \
    go get github.com/dgrijalva/jwt-go && \
    go get github.com/go-chi/jwtauth

RUN CGO_ENABLED=0 go build -o main .
CMD ["./main"]
