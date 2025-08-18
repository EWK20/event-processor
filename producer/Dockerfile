FROM golang:1.24 AS BUILD

WORKDIR /project

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN make build


FROM alpine:3.22

WORKDIR /app

COPY --from=BUILD /project/bin/ .

## Add the wait script to the image
COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

RUN ls /app

## Launch the wait tool and then application
CMD /wait && /app/producer