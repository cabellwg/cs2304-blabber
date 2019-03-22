# Build stage
FROM golang:alpine3.9 AS build

WORKDIR /go/src/app
COPY . .

RUN apk add git

RUN go get -d -v
RUN go install -v

# Production stage
FROM alpine:3.9

COPY --from=build /go/bin/app /usr/local/bin/

ENTRYPOINT ["app"]
