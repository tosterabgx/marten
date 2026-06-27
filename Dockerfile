FROM golang:1.26.4 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /marten-server ./cmd/marten-server


FROM scratch

COPY --from=build /marten-server /marten-server

EXPOSE 6472

ENTRYPOINT [ "/marten-server" ]
