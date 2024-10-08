FROM golang:1.22.0-alpine

COPY . /app

WORKDIR /app

RUN go mod download

RUN mkdir -p ./bin
RUN go build -o /bin/keeneticToMqtt /app/cmd/

ENV CONFIG_PATH "/data/options.json"

CMD [ "/bin/keeneticToMqtt" ]
