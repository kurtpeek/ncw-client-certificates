This demonstrates how to make client side certificates with go

First generate the certificates with

```
./makecert.sh test@test.com
```

Run the server in one terminal

```
go run main.go server
```

and the client in the other

```
go run main.go client
```
