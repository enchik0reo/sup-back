FROM golang:1.21.6 as build
WORKDIR /src
COPY ./ /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/sup ./cmd/sup/main.go

FROM alpine:3.19.0
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/configs ./configs
COPY --from=build /src/.env ./
COPY --from=build /bin/sup /bin/sup

RUN apk add bash

CMD ["/bin/sup"]