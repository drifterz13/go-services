FROM --platform=linux/amd64 golang:1.16-alpine

ARG SERVICE

RUN apk add curl
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

RUN mkdir /internal
COPY . /internal
WORKDIR /internal/$SERVICE
RUN go mod tidy

CMD ["air"]


