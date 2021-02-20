# Собираем в гошке
FROM golang:1.15.2 as build

ENV BIN_FILE /opt/calendar/calendar-app
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
RUN cd cmd/calendar && CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} .

# На выходе тонкий образ
FROM alpine:3.9
RUN apk update && apk upgrade && apk add bash

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="larionovog8@gmail.com"

ENV BIN_FILE "/opt/calendar/calendar-app"
ENV MIGRATIONS_DIR "/opt/calendar/migrations"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

COPY ./migrations/ ${MIGRATIONS_DIR}

CMD ${BIN_FILE} -migrate ${MIGRATIONS_DIR}
