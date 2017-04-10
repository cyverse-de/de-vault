FROM golang:1.6

ARG git_commit=unknown
LABEL org.cyverse.git-ref="$git_commit"

COPY . /go/src/github.com/cyverse-de/de-vault
RUN go install github.com/cyverse-de/de-vault

RUN apt-get update && apt-get install -y git ssh

ENTRYPOINT ["de-vault"]
CMD ["--help"]
