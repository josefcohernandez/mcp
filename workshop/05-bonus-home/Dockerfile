# Build from the repository root:
#   docker build -f workshop/05-bonus-home/Dockerfile -t devfest-mcp .
#   docker run --rm -p 8080:8080 devfest-mcp
FROM golang:1.25 AS build
WORKDIR /src
# Download modules first so this layer is cached while only code changes.
COPY workshop/05-bonus-home/go.mod workshop/05-bonus-home/go.sum ./
RUN go mod download
COPY workshop/05-bonus-home/*.go ./
# CGO_ENABLED=0 produces a fully static binary. Without it, a program that imports
# "net" is linked against glibc and does not start on distroless/static.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/devfest-mcp .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/devfest-mcp /devfest-mcp
COPY workshop/data/schedule.json /data/schedule.json
ENV SCHEDULE_PATH=/data/schedule.json
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/devfest-mcp", "-http=:8080"]
