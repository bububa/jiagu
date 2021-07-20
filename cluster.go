package jiagu

import (
	"github.com/bububa/jiagu/cluster"
)

// KmeansCluster kmeans 文本聚类
func KmeansCluster(corpus []string, tokenizer cluster.Tokenizer, k int) (cluster.Clusters, error) {
	engine := cluster.NewKmeans()
	engine.SetK(k)
	featureGetter := cluster.NewTfIdfFeaturesGetter(tokenizer)
	return cluster.Clusterize(corpus, engine, featureGetter)
}

// DBScanCluster dbscan 文本聚类
func DBScanCluster(corpus []string, tokenizer cluster.Tokenizer, eps float64, minPts int) (cluster.Clusters, error) {
	if eps <= 1e-15 {
		eps = 0.5
	}
	if minPts == 0 {
		minPts = 2
	}
	engine := cluster.NewDbscan(eps, minPts)
	featureGetter := cluster.NewTfIdfFeaturesGetter(tokenizer)
	return cluster.Clusterize(corpus, engine, featureGetter)
}
