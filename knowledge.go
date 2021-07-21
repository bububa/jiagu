package jiagu

import (
	"compress/gzip"
	"fmt"

	"github.com/bububa/jiagu/knowledge"
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
		gr, err := gzip.NewReader(modelR)
		defer gr.Close()
		knowledgeModel, err = knowledge.NewFromReader(gr)
		if err != nil {
			panic(err)
		}
	}
	return knowledgeModel
}

// Knowledge 知识图谱关系提取
func Knowledge(txt string) []knowledge.Entity {
	model := KnowledgeInstance()
	words := Seg(txt)
	return model.Entities(words)
}
