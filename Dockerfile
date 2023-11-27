FROM golang:alpine
WORKDIR $GOPATH/src/

COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN go build -o ./out/linker .
CMD [". /out/linker"]