package spider

import (
	"github.com/gocolly/colly"
)

// Spider 蜘蛛
// 单个蜘蛛可携带多个收集器
type Spider struct {
	// name 名称
	name string
	// legs 收集器
	collectors []colly.Collector
	// Max 收集器最大数量 默认为1
	maxLegs int64
	// HealthZ 健康检查
	healthZ HealthZ
}

func (s *Spider) Collectors() []colly.Collector {
	return s.collectors
}

func (s *Spider) SetCollectors(collectors []colly.Collector) {
	s.collectors = collectors
}

func (s *Spider) MaxLegs() int64 {
	return s.maxLegs
}

func (s *Spider) SetMaxLegs(maxLegs int64) {
	s.maxLegs = maxLegs
}

func (s *Spider) HealthZ() HealthZ {
	return s.healthZ
}

func (s *Spider) SetHealthZ(healthZ HealthZ) {
	s.healthZ = healthZ
}
