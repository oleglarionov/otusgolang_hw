# Собираем в гошке
FROM golang:1.15.2 as build

ENV BIN_FILE /opt/calendar/sender-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в apline образе.
ARG LDFLAGS
RUN cd cmd/sender && CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} .

# На выходе тонкий образ
FROM alpine:3.9
RUN apk update && apk upgrade && apk add bash

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="sender"
LABEL MAINTAINERS="larionovog8@gmail.com"

ENV BIN_FILE "/opt/calendar/sender-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
