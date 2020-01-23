FROM golang:latest
RUN go get -u github.com/sparrc/go-ping
RUN go get -u github.com/mediocregopher/radix
RUN go get -u github.com/gorilla/mux
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main ./src
EXPOSE 9999
EXPOSE 8080
CMD ["/app/main"]