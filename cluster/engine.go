package cluster

// ClusterEngine 聚类引擎接口
type ClusterEngine interface {
	Clusterize(Observations) (Clusters, error)
}

// Clusterize 文本聚类
func Clusterize(corpus []string, engine ClusterEngine, featureGetter FeaturesGetter) (Clusters, error) {
	features, _ := featureGetter.Features(corpus)
	return engine.Clusterize(features)
}
