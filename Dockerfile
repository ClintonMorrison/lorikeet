FROM ubuntu 

RUN \
  apt-get update && \
  apt-get install -y golang-go && \
  apt-get install -y nginx && \
  apt-get install -y ca-certificates

ENV GOPATH /opt/go

COPY ./nginx.conf /etc/nginx

COPY ./ui/build /opt/go/src/github.com/ClintonMorrison/lorikeet/ui/build
COPY ./cmd /opt/go/src/github.com/ClintonMorrison/lorikeet/cmd
COPY ./internal /opt/go/src/github.com/ClintonMorrison/lorikeet/internal
COPY ./start.sh /opt/go/src/github.com/ClintonMorrison/lorikeet/start.sh

WORKDIR /opt/go/src/github.com/ClintonMorrison/lorikeet
RUN go mod init
RUN go build -o server cmd/server/main.go
RUN go build -o doBackup cmd/backup/main.go
RUN go build -o doRestore cmd/restore/main.go

CMD ["sh", "./start.sh"]
