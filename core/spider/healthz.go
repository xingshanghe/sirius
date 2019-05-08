package spider

// HealthZ 健康检测
type HealthZ struct {
	_ struct{} `type:"structure"`
	// Command 命令行参数
	command []*string
	// Interval 过期时间. 5-300s 默认 30s
	interval *int64
	// Retries 重试次数 1-10次 默认为3次.
	retries *int64
	// StartPeriod 开始计时时间.
	startPeriod *int64
	// Timeout 超时时间 2-60s 默认5.
	timeout *int64
}
