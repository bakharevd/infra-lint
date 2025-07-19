FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install -y    curl

COPY . /app
WORKDIR /app

CMD ["./start.sh"]

HEALTHCHECK CMD curl -f http://localhost/ || exit 1
USER nobody
