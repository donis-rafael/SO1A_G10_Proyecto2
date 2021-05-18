FROM golang
WORKDIR /usr/local/pub
COPY . .
RUN go mod download
CMD ["go","run","sender.go"]
EXPOSE 80