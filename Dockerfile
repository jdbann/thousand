FROM node:16-alpine3.12 AS build-assets
WORKDIR /build
COPY package.json yarn.lock ./
RUN yarn install
COPY bin/build-assets bin/
COPY css/ css/
COPY js/ js/
RUN bin/build-assets

FROM golang:1.16-alpine AS build-app
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download -x
COPY . ./
COPY --from=build-assets /build/static/css/ static/css/
COPY --from=build-assets /build/static/js/ static/js/
RUN CGO_ENABLED=0 go build -o bin/thousand -v .

FROM alpine
WORKDIR /app
COPY --from=build-app /build/bin/thousand bin/thousand
CMD ["bin/thousand"]
