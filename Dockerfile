
# Build the frontend
FROM node:16 as builder

WORKDIR /app/ui

COPY ui/package.json ui/package-lock.json ./
RUN npm install

RUN npx browserslist@latest --update-db

COPY ./ui .
RUN npm run build

# Build the mails
FROM node:16 as mailbuild

WORKDIR /app/mails

COPY mails/package.json mails/package-lock.json ./
RUN npm install

COPY ./mails .
RUN npm run build

# Build the gobinary

FROM golang:1.19 as gobuild

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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/userd ./cmds/userd

FROM gcr.io/distroless/static

COPY --from=gobuild /go/bin/userd /go/bin/userd
EXPOSE 8080

ENTRYPOINT ["/go/bin/userd"]
