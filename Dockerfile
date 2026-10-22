# Stage 1: Build binary using Alpine-Rust (Ultra Lightweight)
FROM rust:alpine as builder

# Required build dependencies for networking and compilation
RUN apk add --no-cache musl-dev openssl-dev openssl-libs-static pkgconfig

WORKDIR /app
COPY . .

# 1 Core execution to keep build RAM low & prevent timeout
ENV CARGO_BUILD_JOBS=1
RUN cargo build --release

# Stage 2: Minimal final runtime image (< 20MB)
FROM alpine:latest

# SSL certificates install in 2 seconds without hanging
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/target/release/siren_dispatch /app/siren_dispatch

CMD ["/app/siren_dispatch"]
