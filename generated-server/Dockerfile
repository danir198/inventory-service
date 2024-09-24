FROM golang:1.17-alpine AS build
WORKDIR /go/src
COPY . .
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -a -installsuffix cgo -o swagger .

FROM scratch AS runtime
COPY --from=build /go/src/swagger ./
EXPOSE 8080/tcp
ENTRYPOINT ["./swagger"]
