# Baron
micro service development framework

Baron handles the painful parts of services, freeing you to focus on the
business logic.

## 安装

Currently, there is no binary distribution of baron, you must install from
source.

To install this software, you must:

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

Then you open the `handlers/handlers.go`,
add you business logic, and you're good to go.

Here is an example service definition: [Echo Service](./_example/echo.proto)

Try baron for yourself on Echo Service to see the service that is generated:

```
baron _example/echo.proto
```

See [USAGE.md](./USAGE.md) and [TUTORIAL.md](./TUTORIAL.md) for more details.

## 开发

See [DEVELOPING.md](./DEVELOPING.md) for details.

## TODO
- 生成 api.pb.kit.go

- 命令行参数 
  * baron start [proto]
  * baron [option] api.proto ,option 增加transport, 默认all(http|grpc|nrpc等), 生成代码根据传输协议生成不同代码
  
- 应用结构调整

  每个微服务独立go.mod

- 

