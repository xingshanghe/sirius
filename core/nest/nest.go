package nest

import (
	"github.com/xingshanghe/sirius/core/spider"
)

// Nest 爬虫蜘蛛的巢穴。
// 单个巢穴 蜘蛛种类相同
// 但是单个蜘蛛携带的收集器数量可能不同
type Nest struct {
	// 巢穴规格描述
	spec    interface{}
	spiders map[string]spider.Spider
}

func (n *Nest) AppendSpider(spider spider.Spider) {

}
