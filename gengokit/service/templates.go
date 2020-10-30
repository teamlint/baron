package service

const (
	// baron
	BaronPath = "baron/NAME.pb.baron.go.tmpl" // baron 模板路径
	// cmd
	CmdClientPath = "cmd/NAME-client/main.go.tmpl" // 服务客户端CLI模板路径
	CmdServerPath = "cmd/NAME/main.go.tmpl"        // 服务端CLI模板路径
	// service
	HookPath           = "service/hooks.go.tmpl"           // 服务hooks模板路径
	MiddlewaresPath    = "service/middlewares.go.tmpl"     // 服务中间件模板路径
	ServicePath        = "service/service.go.tmpl"         // 服务实现模板路径
	ServiceMethodsPath = "service/service_methods.go.tmpl" // 服务方法模板路径
	// server
	ServerPath = "server/server.go.tmpl" // 服务服务端模板路径

)

var (
	ExcludedPath = []string{BaronPath, ServiceMethodsPath} // 排除的模板路径,由程序手动调用
)

// IsExcludedPath 是否是排除的模板路径
func IsExcludedPath(path string) bool {
	for _, p := range ExcludedPath {
		if p == path {
			return true
		}
	}
	return false
}
