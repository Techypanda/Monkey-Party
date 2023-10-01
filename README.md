# Monkey Party
[![Go Coverage](https://github.com/Techypanda/Monkey-Party/wiki/coverage.svg)](https://raw.githack.com/wiki/Techypanda/Monkey-Party/coverage.html)

## Contents
- [How To Test This Project](#testing)
- [Roomagent](#roomagent)

## Testing
```sh
make test

or

go test

or if you want coverage:

go test -v -coverprofile cover.out
go tool cover -html cover.out -o cover.html 
```

## RoomAgent

### Running Locally
```sh
make roomagent
```

### Testing
```sh
make test-roomagent
```

#### Local Manual Testing
```sh
make grpc-gui
```