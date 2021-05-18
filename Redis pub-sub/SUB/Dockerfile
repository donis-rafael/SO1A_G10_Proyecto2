FROM golang
WORKDIR /
COPY . .
RUN go mod download
EXPOSE 8081
CMD ["go","run","receive.go"]