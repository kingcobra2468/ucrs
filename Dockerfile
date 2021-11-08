FROM golang:1.17

ARG ucrs_hostname=127.0.0.1
ARG ucrs_port=8080
ARG ucrs_redis_hostname=127.0.0.1
ARG ucrs_redis_port=6379
ARG ucrs_fcm_topic=un
ARG google_application_creds=./

ENV UCRS_HOSTNAME=${ucrs_hostname}
ENV UCRS_PORT=${ucrs_port}
ENV UCRS_REDIS_HOSTNAME=${ucrs_redis_hostname}
ENV UCRS_REDIS_PORT=${ucrs_redis_port}
ENV UCRS_FCM_TOPIC=${ucrs_fcm_topic}
ENV GOOGLE_APPLICATION_CREDENTIALS=${google_application_creds}

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build

EXPOSE ${ucrs_port}

CMD ["ucrs"]