# build stage
FROM golang:1.10-alpine AS build-env
ARG APP_NAME=goapp

RUN apk add --no-cache curl bash git openssh
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/govinda-attal/cart-commerce
WORKDIR /go/src/github.com/govinda-attal/cart-commerce
RUN dep ensure -v
RUN go build -o cartcom

# final stage
FROM alpine:3.7
RUN apk -U add ca-certificates

WORKDIR /app
COPY --from=build-env /go/src/github.com/govinda-attal/cart-commerce/cartcom /app/
COPY --from=build-env /go/src/github.com/govinda-attal/cart-commerce/api /app/api
COPY --from=build-env /go/src/github.com/govinda-attal/cart-commerce/rules /app/rules
COPY --from=build-env /go/src/github.com/govinda-attal/cart-commerce/config /app/config

VOLUME [ "/app/config", "/app/api", "/app/rules"  ]
EXPOSE 9080

CMD ["./cartcom"]