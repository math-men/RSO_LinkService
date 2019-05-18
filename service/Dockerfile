FROM golang:1.11-alpine
ENV SOURCES /go/src/
COPY . ${SOURCES}
WORKDIR ${SOURCES}
RUN apk add --no-cache git && \
    go get github.com/go-chi/chi && \
    go get github.com/go-chi/chi/middleware && \
    go get github.com/go-sql-driver/mysql && \
    go get github.com/dgrijalva/jwt-go && \
    go get github.com/go-chi/jwtauth

RUN CGO_ENABLED=0 go build -o main .
CMD ["./main"]
