##
## Build
##
FROM golang:alpine as build
RUN apk add git
WORKDIR /app

COPY . .
RUN apk add --update --no-cache ca-certificates git
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /employee-service -ldflags "-X 'main.appVersionLDFlag=$(cat VERSION)' -X 'main.appBuildTimestampLDFlag=$(date -u '+%d-%m-%Y %H:%M:%S %Z')'"


##
## Deploy
##
FROM gcr.io/distroless/static-debian11

WORKDIR /app

COPY --from=build /employee-service /app/employee-service
COPY --from=build /app/swagger.json /app/swagger.json

ENTRYPOINT ["/app/employee-service"]
