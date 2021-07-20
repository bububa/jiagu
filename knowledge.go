package jiagu

import (
	"fmt"

	"github.com/bububa/jiagu/knowledge"
	"github.com/bububa/jiagu/segment"
)

var knowledgeModel *knowledge.Knowledge

// KnowledgeInstance get knowledgeModel singleton
func KnowledgeInstance() *knowledge.Knowledge {
	if knowledgeModel == nil {
		modelR, err := modelFS.Open(fmt.Sprintf("model/%s", KG_MODEL))
		if err != nil {
			panic(err)
		}
		defer modelR.Close()
		knowledgeModel, err = knowledge.NewFromReader(modelR)
		if err != nil {
			panic(err)
		}
	}
	return knowledgeModel
}

// Knowledge 知识图谱关系提取
func Knowledge(txt string) []knowledge.Entity {
	model := KnowledgeInstance()
	words := Seg(txt, segment.Default_SegMode)
	return model.Entities(words)
}
