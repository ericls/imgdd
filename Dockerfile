ARG GIT_REV="docker-dev"
ARG GIT_HASH="unset"
# -----------------------------------------------------
# 1) FRONTEND BUILD STAGE
# -----------------------------------------------------
FROM --platform=$BUILDPLATFORM node:20-alpine AS ui-build
RUN npm install -g pnpm

RUN mkdir /code
WORKDIR /code

# Install dependencies first (leveraging Docker cache)
COPY ./web_client/package*.json ./
RUN pnpm import
RUN pnpm install --frozen-lockfile

# Copy the rest of the frontend source code and build it
COPY web_client ./
RUN npm run build

# -----------------------------------------------------
# 2) BACKEND BUILD STAGE (WITH EMBEDDED FRONTEND FILES)
# -----------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS server-build
ARG GIT_REV
ARG GIT_HASH
WORKDIR /code/imgdd

# Copy only Go dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend source code
COPY ./ ./

# Remove all files in the web_client directory
RUN rm -rf web_client
RUN mkdir -p web_client/dist

# Copy the built frontend files into the backend directory
COPY --from=ui-build /code/dist/ web_client/dist/

# Set cross-compilation environment variables
ARG TARGETPLATFORM
RUN echo "Building for $TARGETPLATFORM"

RUN GOOS=$(echo $TARGETPLATFORM | cut -d '/' -f1) \
  GOARCH=$(echo $TARGETPLATFORM | cut -d '/' -f2) \
  go build \
  -ldflags "-s -w \
    -X 'github.com/ericls/imgdd/buildflag.Debug=false' \
    -X 'github.com/ericls/imgdd/buildflag.Dev=false' \
    -X 'github.com/ericls/imgdd/buildflag.Docker=true' \
    -X github.com/ericls/imgdd/buildflag.Version=$GIT_REV \
    -X github.com/ericls/imgdd/buildflag.VersionHash=$GIT_HASH \
  " \
  -o /go/bin/imgdd .

# -----------------------------------------------------
# 3) FINAL IMAGE (Multi-Arch Support)
# -----------------------------------------------------
FROM alpine:3.21 AS final

ARG GIT_REV
ARG GIT_HASH
LABEL git_commit=$GIT_REV
# Create user and working directories
RUN addgroup -S imgdd && adduser -S imgdd -G imgdd

# Copy the compiled backend binary (which has the embedded frontend files)
COPY --from=server-build /go/bin/imgdd /usr/local/bin/imgdd

USER imgdd
EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/imgdd"]
CMD ["serve"]
