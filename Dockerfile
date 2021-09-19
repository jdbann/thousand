FROM golang:1.16-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/thousand

FROM scratch
COPY --from=build /bin/thousand /bin/thousand
CMD ["/bin/thousand"]
