FROM --platform=linux/amd64 debian:stable-slim

ENV PORT="8080"

RUN apt-get update && apt-get install -y ca-certificates

ADD notely /usr/bin/notely

CMD ["notely"]
