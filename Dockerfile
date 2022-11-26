FROM golang:1.19-alpine as build

# Install Alpine Dependencies
RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk add --no-cache bash git openssh make cmake

# Install dependencies
RUN mkdir /app
COPY ./go.mod /app
COPY ./go.sum /app
WORKDIR /app
RUN go mod download

# Build app
COPY Makefile /app
COPY . /app
RUN make

# Make runtime image
FROM golang:1.19-alpine as runtime
COPY --from=build /app/gomphotherium /app/gomphotherium
ENTRYPOINT [ "/app/gomphotherium", "tui" ]
