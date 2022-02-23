# Teamway test task: Server

Please pay attention:
* The client is in another repository: https://github.com/oltur/teamway-client

Generate doc

```console
$ go get -u github.com/swaggo/swag/cmd/swag
$ swag init
```

Run app

```console
$ go run main.go
```

Run tests

```console
$ go test ./test/... -v -count=1
```

[open swagger](http://localhost:8081/swagger/index.html)

