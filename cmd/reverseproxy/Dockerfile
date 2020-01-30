
#build stage
FROM golang:alpine AS builder 
RUN mkdir /app 
ADD . /app
WORKDIR /app
RUN go build -o proxy ./proxy.go


FROM alpine:latest AS production
COPY --from=builder /app .
EXPOSE 6060
CMD [ "./proxy" ]
