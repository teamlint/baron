# How to develop baron

## Dependencies

1. Everything required to install `baron`
2. go-bindata for compiling templates into binary
```
$ go get github.com/jteeuwen/go-bindata/...
```

## Building

Whenever templates are modified, the templates must be recompiled to binary,
this is done with:

```
$ go generate github.com/teamlint/baron/...
```

Then to build baron and its protoc plugin to your $GOPATH/bin directory:

```
$ go install github.com/teamlint/baron/...
```

Both can be done from the Makefile in the root directory:

```
$ cd $GOPATH/github.com/teamlint/baron
$ make
```

## Testing

Before submitting a pull request always run tests that cover modified code.
Also build baron and run baron's integration test. This can be done by

```
$ cd $GOPATH/src/github.com/teamlint/baron
$ make
$ make test
# If the tests failed and you want to remove generated code
$ make testclean
```

## Structure

baron works as follows:

1. Read in a group of `.proto` files
2. Execute `protoc` with our `protoc-gen-protocast` protoc plugin, which
   outputs the protoc AST representation of the .proto files
3. Parse protoc's AST output and  the `.proto` file with the
   `grpc Service` definition for http annotations using `go-baron/deftree`
4. Use `protoc` and `protoc-gen-go` to generate `.pb.go` files containing
   protobuf structs and transport for golang
5. Use the constructed `deftree` with `gengokit` to template out basic gokit service with grpc
   and http/json transport and empty handlers
6. Generate documentation from comments with `gendocs`

If there was already generated code in the filesystem then baron will not
overwrite user code in the /NAME-service/handlers directory

Additional internal packages of note used by these programs are:

- `deftree`, located in `deftree/`, which makes sense of the protobuf file
  passed to it by `protoc`, and is used by `gengokit` and
  `gendoc`



## Go-Kit

## go-kit目录分析

```shell
.
├── auth 权限验证
├── circuitbreaker 熔断器
├── cmd 自动生成代码命令行工具
├── endpoint 
├── log 日志
├── metrics 监控指标
├── ratelimit 限流
├── sd 服务发现
├── tracing 追踪
├── transport
└── util 工具包
```

## gRPC 服务端选项

`ServerOption`为Serve设置可选的函数调用, 有以下几种: 

1. ServerBefore: 在调用decode函数之前执行，在HTTP请求对象上执行ServerBefore函数
2. ServerAfter: 在调用endpoint之后, encode函数之前执行
3. ServerErrorHandler: 收集decode, endpoint, encode中返回的第二个参数的错误对象信息, 简单的收集log 
4. ServerErrorEncoder: 收集decode, endpoint, encode中返回的第二个参数的错误对象信息, 并可以写入到http.ResponseWriter返回客户端
5. ServerFinalizer: 在每个HTTP请求结束时执行，在encode或者ServerErrorEncoder之后执行的函数

正常的请求流程: `ServerBefore -> decode -> endpoint -> service -> ServerAfter -> encode -> ServerFinalizer`
出现错误的请求流程: `ServerBefore -> 出现错误(decode -> endpoint -> encode) -> ServerErrHandler -> ServerErrorEncoder(可写httpResponse) -> ServerFinalizer`


