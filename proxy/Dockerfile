FROM golang:latest
RUN go get -u github.com/gorilla/mux
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main ./src
EXPOSE 9999
EXPOSE 9998
EXPOSE 8080
CMD ["/app/main"]