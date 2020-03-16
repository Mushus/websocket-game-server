.PHONY: dep, dev, build

dep:
	GO111MODULE=off go get github.com/pilu/fresh

dev:
	fresh

build:
	GOOS=linux   GOARCH=amd64 go build -o ../pkg/gameserver-linux-amd64 .
	GOOS=linux   GOARCH=arm   go build -o ../pkg/gameserver-linux-arm .
	GOOS=windows GOARCH=amd64 go build -o ../pkg/gameserver-win-amd64 .
	GOOS=darwin  GOARCH=amd64 go build -o ../pkg/gameserver-mac-amd64 .