package jiagu

import (
	"testing"

	"github.com/bububa/jiagu/segment"
)

// TestCluster 测试文本聚类
func TestCluster(t *testing.T) {
	corpus := []string{
		"百度深度学习中文情感分析工具Senta试用及在线测试",
		"情感分析是自然语言处理里面一个热门话题",
		"AI Challenger 2018 文本挖掘类竞赛相关解决方案及代码汇总",
		"深度学习实践：从零开始做电影评论文本情感分析",
		"BERT相关论文、文章和代码资源汇总",
		"将不同长度的句子用BERT预训练模型编码，映射到一个固定长度的向量上",
		"自然语言处理工具包spaCy介绍",
		"现在可以快速测试一下spaCy的相关功能，我们以英文数据为例，spaCy目前主要支持英文和德文",
	}
	tokenizer := func(txt string) []string {
		return Seg(txt, segment.Default_SegMode)
	}
	clusters, err := KmeansCluster(corpus, tokenizer, 4)
	if err != nil {
		t.Error(err)
		return
	}
	for _, cluster := range clusters {
		for _, feat := range cluster.Observations {
			t.Logf("%s\n", feat.GetID())
		}
		t.Log("======================")
	}
}
