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
task install
```

## 使用

使用 [proto3](https://developers.google.com/protocol-buffers/docs/proto3) 定义服务 {NAME}.proto
使用 baron 生成基础框架代码 
baron {NAME}.proto
打开 `service/service.go`, 编写业务逻辑处理
启动服务端
```shell
$ go run {NAME}-service/cmd/{NAME}/main.go
```
客户端使用
添加引用包 
```go
import pb "{{MODULE}}/{NAME}"
```
调用服务
```go
conn, err := grpc.Dial(
    grpcAddr,
    grpc.WithInsecure(),
    grpc.WithBlock(),
)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
baronGRPCClient, err := pb.NewGRPCClient(conn)
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()
var in pb.EchoRequest
in.In = "hello"
out, err := baronGRPCClient.Echo(ctx, &in)
if err != nil {
    log.Fatalf("[Baron.GRPCClient] Echo.Echo err=%v\n", err)
}
log.Printf("[Baron.GRPCClient] Echo.Echo result=%+v\n", *out)

```



See [USAGE.md](./docs/USAGE.md) and [TUTORIAL.md](./docs/TUTORIAL.md) for more details.

## 开发

See [DEVELOPING.md](./docs/DEVELOPING.md) for details.

## TODO
- json 编码器需要同时支持普通对象和复杂对象(wrapper.proto), query 查询解码可以去掉复杂类型,使用自定义库
- stream 测试
- server 增加初始化方法

## 参考
- https://github.com/solo726/bookinfo 使用go-kit实现微服务,truss自动生成go-kit代码
- https://github.com/OahcUil94/go-kit-training go-kit微服务套件使用
- https://github.com/phungvandat/clean-architecture Example about clean architecture in golang
- https://github.com/nametake/protoc-gen-gohttp protoc plugin to generate to Go's net/http converter
- https://github.com/grpc-ecosystem/grpc-gateway/blob/4ba7ec0bc390cae4a2d03625ac122aa8a772ac3a/protoc-gen-grpc-gateway/httprule/parse.go

## 注意事项
- HTTP 服务请求支持如果使用查询字符串传递复杂数据类型, 需要将字段值编码为JSON并做URL编码
  例:
  ```proto
    message EchoRequest {
      google.protobuf.StringValue json_str  = 6;
    }
  ```
  使用 URL 查询参数传值
  ```
  http://localhost:5050/echo?json_str=%7B%22value%22%3A%22Hello%E4%B8%96%E7%95%8C%22%7D
  ```
  参数值:
  参数 JSON 序列化
  `{"value":"Hello世界"}`
  URL 编码
  `%7B%22value%22%3A%22Hello%E4%B8%96%E7%95%8C%22%7D`

## 已知问题
- HTTP 服务请求查询字符串(url query) 仅支持标量类型, 请求体(http body)无限制 
- `google/protobuf/struct.proto` 生成的字段必须设置值,设置 nil 报错

