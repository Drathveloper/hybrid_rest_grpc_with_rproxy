FROM alpine:latest AS build
RUN apk update
RUN apk upgrade
RUN apk add --no-cache go
WORKDIR /app
ADD . /app
RUN go get -d -v /app/...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o poc

FROM alpine:latest
WORKDIR /app
RUN apk update
RUN apk upgrade
RUN apk add --no-cache libc6-compat
COPY --from=build /app/poc /app/poc
CMD ["./poc"]