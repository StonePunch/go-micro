# Build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY authApp /app

# Run executable created in the first image and then transfered to the second, smaller image
CMD [ "/app/authApp" ]