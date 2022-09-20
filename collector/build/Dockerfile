FROM ubuntu:20.04

RUN apt-get update && apt-get install -y \
curl

EXPOSE 8081

COPY build/api ./

ENTRYPOINT ./api
