# # base go image
# FROM golang:1.24-alpine AS builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# #CGO disabled for using just standard library, not needing any c libraries
# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

# NOTE! build stage will be handled via Make file

#build a tiny docker image

FROM alpine:latest

RUN mkdir /app

#copy the built binary from the builder stage
# COPY --from=builder /app/brokerApp /app 

COPY brokerApp /app 

CMD ["/app/brokerApp"]