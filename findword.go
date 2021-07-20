package jiagu

import (
	"io"

	"github.com/bububa/jiagu/findword"
	"github.com/bububa/jiagu/utils"
)

// Findword  根据文本利用信息熵新词发现
func Findword(input io.Reader, minFreq int, minMtro float64, minEntro float64) ([]utils.StringCounterItem, error) {
	model := findword.New()
	if minFreq > 0 {
		model.SetMinFreq(minFreq)
	}
	if minMtro > 1e-15 {
		model.SetMinMtro(minMtro)
	}
	if minEntro > 0 {
		model.SetMinEntro(minEntro)
	}
	return model.Find(input)
}
