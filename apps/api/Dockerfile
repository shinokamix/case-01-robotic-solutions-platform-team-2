FROM golang:1.26.4-alpine AS dev
WORKDIR /workspace/apps/api
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api ./
COPY packages/contracts /workspace/packages/contracts
CMD ["go", "tool", "air"]

FROM golang:1.26.4-alpine AS build
WORKDIR /src/apps/api
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api ./
COPY packages/contracts /src/packages/contracts
RUN CGO_ENABLED=0 go build -trimpath -o /out/api ./cmd/api \
    && CGO_ENABLED=0 go build -trimpath -o /out/admin ./cmd/admin \
    && CGO_ENABLED=0 go build -trimpath -o /out/migrate ./cmd/migrate

FROM gcr.io/distroless/static-debian13:nonroot AS api
COPY --from=build /out/api /api
USER nonroot:nonroot
ENTRYPOINT ["/api"]

FROM gcr.io/distroless/static-debian13:nonroot AS migrate
COPY --from=build /out/migrate /migrate
USER nonroot:nonroot
ENTRYPOINT ["/migrate"]

FROM gcr.io/distroless/static-debian13:nonroot AS admin
COPY --from=build /out/admin /admin
USER nonroot:nonroot
ENTRYPOINT ["/admin"]
