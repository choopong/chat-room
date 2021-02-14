FROM golang:1.15 as build-stage
WORKDIR /app/
COPY ./app/ ./
RUN go mod download
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o webapp main.go

FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ Asia/Bangkok
WORKDIR /app/
COPY --from=build-stage app/webapp .
RUN chmod a+x webapp 
ENTRYPOINT ["./webapp"]