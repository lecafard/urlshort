FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/lecafard/urlshort/

COPY . .

# Fetch dependencies
# Using go get.
RUN go get -d -v # Build the binary.
RUN CGO_ENABLED=0 go build -o /go/bin/urlshort
RUN mkdir -p /build/data && mkdir -p /build/bin && cp /go/bin/urlshort /build/bin/urlshort && cp -r $GOPATH/src/github.com/lecafard/urlshort/static /build/bin/static

FROM scratch
# Copy our static executable.
COPY --from=builder /build /

# Set environment variables and cwd
WORKDIR /bin
ENV URLSHORT_LISTEN=":8000" URLSHORT_DATA="/data/data.db"
EXPOSE 8000/tcp

ENTRYPOINT ["/bin/urlshort"]
