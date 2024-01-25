FROM alpine:latest

RUN mkdir /app

# The loggerApp binary is built with "make build_logger"
COPY loggerApp /app

CMD [ "/app/loggerApp" ]