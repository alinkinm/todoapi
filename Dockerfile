FROM golang
USER root
SHELL ["/bin/bash", "-c"]

RUN go version
RUN mkdir -p /usr/applications/todoapi
WORKDIR /usr/applications/todoapi

COPY . /usr/applications/todoapi
COPY ./configs/config.yaml /usr/applications/todoapi/configs

RUN go mod download
RUN go build -o todoapi cmd/api/main.go

CMD ./todoapi