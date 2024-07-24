FROM reporter:latest as builder


FROM golang:1.21
WORKDIR /app
COPY --from=builder ./app/bin/migrate/migrate ./bin/

ENTRYPOINT [ "./bin/migrate" ]