
#build stage
FROM golang:alpine AS builder 
RUN mkdir /app 
ADD . /app
WORKDIR /app
RUN apk add --no-cache git
RUN go get -d
RUN go build -o logger ./log-server.go


FROM alpine:latest AS production
COPY --from=builder /app .
EXPOSE 6061
CMD [ "./logger" ]