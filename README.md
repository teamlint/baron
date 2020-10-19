# Baron
micro service development framework

Baron 根据 proto 文件快速生成 *.pb.go *.pb.baron.go(包含go-kit服务Endpoints,Transports)及微服务框架代码布局, 让你专注于业务逻辑处理.

## 安装

1. Install protoc 3 or newer. The easiest way is to
download a release from [github](https://github.com/google/protobuf/releases)
and add to `$PATH`.
Otherwise [install from source.](https://github.com/google/protobuf)
1. Install baron with

	```
	go get -u -d github.com/teamlint/baron
	cd $GOPATH/src/github.com/teamlint/baron
	make dependencies
	make
	```
	On Windows, do the following instead:
	```
	go get -u -d github.com/teamlint/baron
	cd %GOPATH%/src/github.com/teamlint/baron
	wininstall.bat
	```

## 使用

Using baron is easy. You define your service with [gRPC](http://www.grpc.io/)
and [protoc buffers](https://developers.google.com/protocol-buffers/docs/proto3),
and baron uses that definition to create an entire service. You can even
add [http annotations](
https://github.com/googleapis/googleapis/blob/928a151b2f871b4239b7707e1bb59258df3fe10a/google/api/http.proto#L36)
for HTTP 1.1/JSON transport!

Then you open the `service/service.go`,
add you business logic, and you're good to go.

Here is an example service definition: [Echo Service](./_example/echo.proto)

Try baron for yourself on Echo Service to see the service that is generated:

```
baron _example/echo.proto
```

See [USAGE.md](./docs/USAGE.md) and [TUTORIAL.md](./docs/TUTORIAL.md) for more details.

## 开发

See [DEVELOPING.md](./docs/DEVELOPING.md) for details.

## TODO
- NATS 传输协议支持
  
- 

## 参考
- https://github.com/solo726/bookinfo 使用go-kit实现微服务,truss自动生成go-kit代码
- https://github.com/OahcUil94/go-kit-training go-kit微服务套件使用
- https://github.com/phungvandat/clean-architecture Example about clean architecture in golang

