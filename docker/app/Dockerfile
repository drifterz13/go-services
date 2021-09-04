FROM golang:1.16-alpine

RUN go get github.com/cespare/reflex@latest

COPY reflex.conf /
COPY start.sh /

CMD ["reflex", "-c", "/reflex.conf"]


