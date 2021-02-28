FROM golang:1.16
# create a working directory
WORKDIR /root/data/

# downloading packages
COPY go.mod .
COPY go.sum .

RUN go mod download

# add source code
COPY . /root/data/

RUN go build

RUN go test -v ./usecase

# run main.go
CMD ["go","run","main.go"]