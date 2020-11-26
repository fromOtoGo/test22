FROM golang
COPY /solution/server/. /go/src/solution/server/
RUN go get "github.com/gorilla/mux"
RUN go install /go/src/solution/server/
ENTRYPOINT /go/bin/server
EXPOSE 8080

