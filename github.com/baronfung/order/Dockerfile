FROM golang:latest

# for install go packages RUN go get /path
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/baronfung/order


# build executable
RUN go install github.com/baronfung/order


# execute 
ENTRYPOINT /go/bin/order

# Document that the service listens on port 8080.
EXPOSE 8080

COPY ./init.sql /docker-entrypoint-initdb.d/init.sql
RUN chmod +x /docker-entrypoint-initdb.d/init.sql