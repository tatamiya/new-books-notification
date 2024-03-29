# original is https://github.com/GoogleCloudPlatform/cloud-code-samples/blob/v1/golang/go-hello-world/Dockerfile

# Use base golang image from Docker Hub
FROM golang:1.15 AS build

WORKDIR /crawler

# Install dependencies in go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the application source code
COPY ./src ./src

# Compile the application to /app.
# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN echo "Go gcflags: ${SKAFFOLD_GO_GCFLAGS}"
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -mod=readonly -v -o /app ./src

# Now create separate deployment image
FROM gcr.io/distroless/base

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

# Copy assets
WORKDIR /crawler
COPY --from=build /app ./app
COPY favorites.json favorites.json
COPY ./src/details/ccode.json ./src/subject/ccode.json

ENTRYPOINT ["./app"]
