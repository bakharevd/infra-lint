FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y curl

COPY . /app
WORKDIR /app

CMD ["./start.sh"]