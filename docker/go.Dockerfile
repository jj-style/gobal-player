FROM golang:1.22-bookworm

RUN apt update && apt install -y libvlc-dev vlc-plugin-base vlc-plugin-video-output curl && rm -rf /var/lib/apt/lists/*

RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /bin
