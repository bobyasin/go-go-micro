FROM alpine:latest

RUN mkdir /app

#copy the built binary from the builder stage
# COPY --from=builder /app/brokerApp /app 

COPY loggerServiceApp /app 

CMD ["/app/loggerServiceApp"]