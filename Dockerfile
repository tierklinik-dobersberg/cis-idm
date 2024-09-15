
# Build the frontend
FROM node:16 as builder

WORKDIR /app/ui

COPY ui/.npmrc ui/package.json ui/package-lock.json ./
RUN --mount=type=secret,id=github_token \
  export GITHUB_TOKEN="$(cat /run/secrets/github_token)" && npm install

RUN npx browserslist@latest --update-db

COPY ./ui .
RUN npm run build && rm -r .angular/cache node_modules

# Build the mails
FROM node:16 as mailbuild

WORKDIR /app/mails

COPY mails/package.json mails/package-lock.json ./
RUN npm install

COPY ./mails .
RUN npm run build

# Build the gobinary

FROM golang:1.23 as gobuild

RUN update-ca-certificates

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY ./ ./
COPY --from=builder /app/cmds/userd/static/ui /go/src/app/cmds/userd/static/ui
COPY --from=mailbuild /app/mails/dist /go/src/app/internal/tmpl/templates/mail/

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags "sqlite_foreign_keys" -v -ldflags "-s -w -linkmode external -extldflags -static" -o /go/bin/userd ./cmds/userd

FROM gcr.io/distroless/static

COPY --from=gobuild /go/bin/userd /go/bin/userd
EXPOSE 8080

ENTRYPOINT ["/go/bin/userd"]
