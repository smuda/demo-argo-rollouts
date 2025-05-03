FROM docker.io/golang:alpine3.21 AS src

COPY . /app-src

RUN cd /app-src \
    && mkdir /app \
    && go mod tidy \
    && go build -o ./binary .

FROM scratch

COPY --from=src /app-src/binary /app

CMD ["/app"]
