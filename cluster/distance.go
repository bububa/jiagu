package cluster

import "math"

// elu 计算两点之间的欧氏距离并返回
func elu(p1 []float64, p2 []float64) float64 {
	var dist float64
	for i, v := range p1 {
		dist += math.Pow(v-p2[i], 2)
	}
	return math.Sqrt(dist)
}
