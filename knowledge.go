package jiagu

import (
	"github.com/bububa/jiagu/knowledge"
	"github.com/bububa/jiagu/segment"
)

var knowledgeModel *knowledge.Knowledge

// Knowledge 知识图谱关系提取
func Knowledge(txt string) []knowledge.Entity {
	words := Seg(txt, segment.Default_SegMode)
	return knowledgeModel.Entities(words)
}
