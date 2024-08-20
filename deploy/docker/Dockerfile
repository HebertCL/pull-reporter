FROM golang:1.22 as build
ENV CGO_ENABLED 0

COPY . /build
WORKDIR /build/app/prp
RUN go build -o reporter

FROM alpine:latest
COPY --from=build /build/app/prp/. /app/prp/.
WORKDIR /app/prp
CMD [ "./reporter" ]