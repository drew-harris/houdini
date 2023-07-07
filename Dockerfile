# To build the image:
#     docker build -t ghcr.io/go-rod/rod -f lib/docker/Dockerfile .
#

# build rod-manager
FROM golang as go


COPY . /houdini
WORKDIR /houdini
RUN go build

FROM ubuntu:jammy

RUN mkdir houdini
WORKDIR /houdini

COPY --from=go /houdini/houdini ./houdini

ARG apt_sources="http://archive.ubuntu.com"

RUN sed -i "s|http://archive.ubuntu.com|$apt_sources|g" /etc/apt/sources.list && \
    apt-get update > /dev/null && \
    apt-get install --no-install-recommends -y \
    # chromium dependencies
    libnss3 \
    libxss1 \
    libasound2 \
    libxtst6 \
    libgtk-3-0 \
    libgbm1 \
    ca-certificates \
    # fonts
    fonts-liberation fonts-noto-color-emoji fonts-noto-cjk \
    # timezone
    tzdata \
    # process reaper
    dumb-init \
    > /dev/null && \
    # cleanup
    rm -rf /var/lib/apt/lists/*

# process reaper

EXPOSE 8080

CMD ["./houdini"]
