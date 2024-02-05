# Собираем в гошке
FROM golang:1.19 as build

ENV BIN_FILE /opt/sender/sender-app
ENV CODE_DIR /go/src/
WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
#COPY ./go.mod .
#COPY ./go.sum .
#RUN go mod download
COPY ./ ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 GOOS=linux go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ${CODE_DIR}cmd/sender/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="sender"
LABEL MAINTAINERS="ids79@otus.ru"

ENV BIN_FILE "/opt/sender/sender-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/sender/config.toml
COPY ./configs/sender_config.toml ${CONFIG_FILE}

CMD ${BIN_FILE}