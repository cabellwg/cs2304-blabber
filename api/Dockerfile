# Build/development stage
FROM golang:alpine3.9 AS dev

ARG mode
ENV MODE $mode

RUN apk add git

WORKDIR /go/src/api
COPY . .

RUN go get -d -v ./...
RUN go install -v

CMD if [ ${MODE} = dev ]; \
    then \
  	  go get github.com/pilu/fresh && \
      echo $GOPATH; \
  	  fresh; \
  fi

# Production stage
FROM alpine:3.9 AS prod

COPY --from=dev /go/bin/api /usr/local/bin/

ENTRYPOINT ["api"]
