# Baron
micro service development framework

Baron 根据 proto 文件快速生成 [go-kit](https://github.com/go-kit/kit) 微服务框架, 让您专注于业务功能开发😉



## 功能特性

- [x] Service 差异化代码生成

- [x] Google 官方新版 [proton-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) 代码生成

- [x] Google 官方新版 [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) 代码生成

- [x] [gRPC](https://grpc.io/) 传输协议支持

- [x] HTTP 传输协议支持

- [x] [NATS](https://nats.io/) 传输协议支持

- [x] Proto3 `optional` 字段属性支持

- [x] `google/protobuf/any.proto` 字段类型支持

- [x] `google/protobuf/empty.proto` 字段类型支持

- [x] `google/protobuf/timestamp.proto` 字段类型支持

- [x] `google/protobuf/duration.proto` 字段类型支持

- [x] `google/protobuf/wrappers.proto` 字段类型支持

- [x] `google/protobuf/struct.proto` 字段类型支持

- [ ] gRPC Stream 

- [ ] server 初始化中间件

  

## 安装

### 安装 protoc 工具

[下载](https://github.com/protocolbuffers/protobuf) 并安装 protocol buffer 编译工具

### 安装 protoc 插件

```shell
$ export GO111MODULE=on  # Enable module mode
$ go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

###  安装 baron 命令行

```shell
$ go get -u -d github.com/teamlint/baron
cd $GOPATH/src/github.com/teamlint/baron
task install
```

*注: 使用`task`作为编译部署工具, 详细[官方文档](https://taskfile.dev/)*

## 命令行

### 简介

- 命令行名称: `baron`

- 命令行帮助: `baron -h`

  ```shell
  $ baron -h
  baron (version: v0.2.5 version date: 2020-11-05T15:36:20+08:00)
  
  Usage: baron [options] <protofile>...
  
  Generates baron(go-kit) services using proto3 and gRPC definitions.
  
  Options:
    -c, --client             Generate NAME-service client
    -h, --help               Print usage
    -s, --start              Output a 'start.proto' protobuf file in ./
    -d, --svcdef             Print service definition
    -o, --svcout string      Go package path where the generated Go service will be written. Trailing slash will create a NAME-service directory
    -t, --transport string   Service transport protocol: [grpc|nats] (default "all")
    -v, --verbose            Verbose output
    -V, --version            Print version
  ```

### 服务端使用

#### 1.初始化服务定义

**命令名称**:  `baron --start <pkgname> [outdir]`

- 省略`pkgname`参数则使用`start`作为默认包名, 并生成 `start.proto`文件
- 省略`outdir`参数则在当前目录生成 `<pkgname>.proto`文件

```shell
$ baron --start foo
INFO[0000] A "start" protobuf file named 'foo.proto' has been created in the
current directory. You can generate a service based on this new protobuf file
at any time using the following command:

    baron foo.proto

If you want to generate a protofile with a different name, use the
'--start' option with the name of your choice after '--start'. For
example, to generate a 'foo.proto', use the following command:

    baron --start foo
```

#### 2.编写服务定义

打开` <pkgname>.proto`文件, 使用 [proto3](https://developers.google.com/protocol-buffers/docs/proto3) 语法编写服务相关定义

#### 3.生成服务基础框架代码

```shell
$ baron foo.proto
```

查看生成文件

```shell
$ ls
foo-service     foo.pb.baron.go foo.pb.go       foo.proto       foo_grpc.pb.go
```

**外部引用**代码生成在` [outdir]`目录, 默认为当前目录, 实际开发建议使用独立目录

所有**服务端**代码放置在 <name>-service 子目录

查看服务端代码结构

```shell
$ tree -L 3 ./foo-service
./foo-service
├── cmd
│   └── foo
│       └── main.go
├── server
│   └── server.go
└── service
    ├── hooks.go
    ├── middlewares.go
    └── service.go

4 directories, 5 files
```

#### 4.编写业务逻辑代码

打开`service.go`文件编写业务逻辑代码

*注: 服务定义修改后, 使用`baron foo.proto`重新生成代码即可, 会对新老代码进行差异化修改*

### 客户端使用

#### 1.添加引用包 

```go
import pb "<gomod>/<proto-name>" // [outdir] 目录下所有代码
```
#### 2.调用服务

##### GRPC 客户端调用

```go
// Baron GRPC Client
log.Println("[GRPC.Start]")
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

// GRPC.Start.Status
{
  ctx := context.Background()
  var in pb.StatusRequest
  out, err := baronGRPCClient.Status(ctx, &in)
  if err != nil {
    log.Fatalf("[Baron.GRPCClient] Start.Status err=%v\n", err)
  }
  log.Printf("[Baron.GRPCClient] Start.Status result=%+v\n", *out)
  log.Println()
}
```

##### HTTP 客户端调用

```go
// Baron HTTP Client
log.Println("[HTTP][Start]")

// HTTP.Start.Status
{
  ctx := context.Background()
  var in pb.StatusRequest
  baronHTTPClient, err := pb.NewHTTPClient(httpAddr)
  if err != nil {
    log.Fatal(err)
  }
  out, err := baronHTTPClient.Status(ctx, &in)
  if err != nil {
    log.Fatalf("[Baron.HTTPClient] Start.Status err=%v\n", err)
  }
  log.Printf("[Baron.HTTPClient] Start.Status result=%+v\n", *out)
  log.Println()
}
```

##### NATS 客户端调用

```go
// Baron NATS Client
log.Println("[NATS.Start]")
nc, err := nats.Connect(natsAddr)
if err != nil {
  log.Fatal(err)
}
defer nc.Close()
baronNATSClient, err := pb.NewNATSClient(nc)
if err != nil {
  log.Fatal(err)
}

// NATS.Start.Status
{
  ctx := context.Background()
  var in pb.StatusRequest
  out, err := baronNATSClient.Status(ctx, &in)
  if err != nil {
    log.Fatalf("[Baron.NATSClient] Start.Status err=%v\n", err)
  }
  log.Printf("[Baron.NATSClient] Start.Status result=%+v\n", *out)
  log.Println()
}
```



## 注意事项
- HTTP 服务请求如果使用查询字符串传递复杂数据类型, 需要将字段值编码为JSON并做URL编码, 使用请求体传值可直接使用**原始值**

	**服务定义**
  
  ```protobuf
  message EchoRequest {
  	google.protobuf.StringValue json_str  = 6;
  }
```
  
  **使用 URL 查询参数传值**
  
  - 将参数 JSON 序列化
  `{"value":"Hello世界"}`
  
  - URL 编码
    `%7B%22value%22%3A%22Hello%E4%B8%96%E7%95%8C%22%7D`
    ```http
    http://localhost:5050/echo?json_str=%7B%22value%22%3A%22Hello%E4%B8%96%E7%95%8C%22%7D
    ```
  
  **使用 HTTP Body 传值**
  
  ```json
  {
  	"json_str": "Hello世界"
  }
  ```



## 已知问题

- `google/protobuf/struct.proto` 生成的字段必须设置值,设置 nil 报错



## 参考

- [truss](https://github.com/metaverse/truss) Truss helps you build go-kit microservices without having to worry about writing or maintaining boilerplate code.
- https://github.com/solo726/bookinfo 使用go-kit实现微服务,truss自动生成go-kit代码
- https://github.com/OahcUil94/go-kit-training go-kit微服务套件使用
- https://github.com/phungvandat/clean-architecture Example about clean architecture in golang
- https://github.com/nametake/protoc-gen-gohttp protoc plugin to generate to Go's net/http converter
- https://github.com/grpc-ecosystem/grpc-gateway/blob/4ba7ec0bc390cae4a2d03625ac122aa8a772ac3a/protoc-gen-grpc-gateway/httprule/parse.go

## 

