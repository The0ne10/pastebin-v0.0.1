FROM golang:1.22.3-bullseye


RUN mkdir -p /usr/local/src/app/pastebin
WORKDIR /usr/local/src/app/pastebin

# env
ENV CONFIG_PATH=./config/local.yaml

# dependecies
COPY app/go.mod app/go.sum ./
RUN go mod download

# build
COPY app ./
RUN go build -o ./bin/pastebin ./cmd/pastebin/main.go

# run
CMD ["./bin/pastebin"]
