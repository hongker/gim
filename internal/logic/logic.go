package logic

// logic服务由定时任务和grpc服务组成
// 定时任务负责将消息推送到客户端
// grpc服务负责处理网关转发的用户tcp请求
import (
	"gim/internal/logic/application"
	"gim/internal/logic/config"
	"gim/internal/logic/infrastructure"
	"gim/internal/logic/interfaces"
	"gim/pkg/system"
	"go.uber.org/dig"
	"log"
)

func Run() {
	// 通过DI初始化依赖
	container := dig.New()
	infrastructure.Inject(container)
	application.Inject(container)
	interfaces.Inject(container)

	// 启动
	system.SecurePanic(container.Invoke(serve))

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

// serve 启动服务
func serve(srv *interfaces.Server, job *interfaces.Job) {
	// 初始化配置项
	conf := config.Init()

	// 启动定时任务
	job.Start()

	// 启动grpc服务
	system.SecurePanic(srv.Start(conf.RPC))
}
