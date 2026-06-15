FROM --platform=$BUILDPLATFORM golang:1.26.1-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/draftbot .

FROM scratch
COPY --from=build /out/draftbot /bin
ENTRYPOINT [ "/bin/draftbot" ]