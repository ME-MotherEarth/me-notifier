FROM golang:1.17.6 as builder

MAINTAINER MotherEarthNetwork

WORKDIR /motherearth
COPY . .

WORKDIR /motherearth/cmd/notifier

RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:20.04
COPY --from=builder /motherearth/cmd/notifier /motherearth

EXPOSE 8080

WORKDIR /motherearth

ENTRYPOINT ["./notifier"]
CMD ["--api-type", "rabbit-api"]
