# Baron
micro service development framework

Baron 根据 proto 文件快速生成 *.pb.go *.pb.baron.go(包含go-kit服务Endpoints,Transports)及微服务框架代码布局, 让你专注于业务逻辑处理.

## 安装

### 安装 protoc 工具

[下载](https://github.com/protocolbuffers/protobuf) 并安装 protocol buffer 编译工具

### 安装 protoc GO 语言代码生成插件

```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go
```

### 安装 protoc GRPC GO 语言代码生成插件

详细文档参考 [https://github.com/grpc/grpc-go](https://github.com/grpc/grpc-go)

```shell
$ go get -u google.golang.org/grpc
```

###  安装 baron 微服务框架代码生成工具

```
go get -u -d github.com/teamlint/baron
cd $GOPATH/src/github.com/teamlint/baron
make dependencies
make
```
Windows 操作系统使用下面安装方式:
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

枚举值为空且Map值为空,即为BaseType,表示
```go
// Indicates whether this field represents a basic protobuf type such as
// one of the ints, floats, strings, bools, etc. Since we can only create
// automatic marshaling of base types, if this is false a warning is given
// to the user.

if oneofType.Type.Enum == nil && oneofType.Type.Map == nil {
    option.IsBaseType = true
}

```

## TODO
- 支持标量切片类型
- server 增加初始化方法

## 参考
- https://github.com/solo726/bookinfo 使用go-kit实现微服务,truss自动生成go-kit代码
- https://github.com/OahcUil94/go-kit-training go-kit微服务套件使用
- https://github.com/phungvandat/clean-architecture Example about clean architecture in golang
- https://github.com/nametake/protoc-gen-gohttp protoc plugin to generate to Go's net/http converter
- https://github.com/grpc-ecosystem/grpc-gateway/blob/4ba7ec0bc390cae4a2d03625ac122aa8a772ac3a/protoc-gen-grpc-gateway/httprule/parse.go

## 问题
- 不支持结构体类型转化为非指针类型
