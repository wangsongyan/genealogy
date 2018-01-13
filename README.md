# genealogy 家谱
## installation
1. install govendor
```
go get -u github.com/kardianos/govendor
```
2. get project source code
```
go get -u https://github.com/wangsongyan/genealogy
```
3. fetch dependent package
```
cd $GOPATH/src/github.com/wangsongyan/genealogy
govendor fetch
```
4. modify [/models/models.go](https://github.com/wangsongyan/genealogy/blob/master/models/models.go#L61) dsn
5. run
```
go run main.go
```
