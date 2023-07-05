FROM heroiclabs/nakama-pluginbuilder:3.16.0 AS go-builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /backend
COPY go.mod .
COPY cmd cmd/
COPY internal internal/
COPY vendor vendor/

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get install -y --no-install-recommends gcc libc6-dev

RUN go build --trimpath --mod=vendor --buildmode=plugin -o ./backend.so ./cmd/fileHandler

FROM heroiclabs/nakama:3.16.0

COPY --from=go-builder /backend/backend.so /nakama/data/modules/
COPY local.yml /nakama/data/
COPY ./test/files/1.0.0.json /nakama/data/
COPY ./test/files/1.0.1.json /nakama/data/
