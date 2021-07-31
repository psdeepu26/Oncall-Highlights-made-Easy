FROM golang:1.15.7
WORKDIR  $GOPATH/src/github.com/spatrayuni/tobs-oncall-highlights
COPY . .
CMD ["go","run","main.go"]