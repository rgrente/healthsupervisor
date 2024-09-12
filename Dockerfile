FROM ext-dockerio.artifactory.si.francetelecom.fr/golang:1.23-alpine AS builder

LABEL name="healthsupervisor" 
LABEL version="1.0"

WORKDIR /app
COPY ./src ./src
WORKDIR /app/src
RUN go mod download && \
    CGO_ENABLED=0 go build healthsupervisor
# ENTRYPOINT ["go", "run", "main.go"]      
# EXPOSE 8080

FROM scratch

COPY --from=builder /app/src/healthsupervisor /healthsupervisor

EXPOSE 80
ENTRYPOINT ["/healthsupervisor"]
