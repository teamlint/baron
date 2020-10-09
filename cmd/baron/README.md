# Baron

`baron` 读取 gRPC 服务定义文件, 生成下面内容:

1. 基于注释生成 Markdown 和 HTML 文档 
2. 基于 [Go Kit](http://gokit.io) 微服务工具生成服务端代码:
	- gRPC transport
	- HTTP/JSON transport (including all encoding/decoding)
	- no-op handler methods for each service RPC, ready for business logic to be added
3. 生成微服务客户端代码:
	- gRPC transport
	- HTTP/JSON transport (including all encoding/decoding)
	- no-op handler methods for each service RPC, ready for marshalling command line
      arguments into a request object and sending a request to a server
4. Web API 浏览器 (基于 Swagger 生成)

