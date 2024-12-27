package dec

type UploadStruct struct {
	LocalLocation  string // 本地目录
	Host           string // 远程主机
	Port           uint   // 远程主机端口
	HostUser       string // 远程主机用户名
	HostPass       string // 远程主机密码
	TargetLocation string // 目标目录
	CreateDir      string // 创建目录
}
