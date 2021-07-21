# jiagu自然语言处理工具
[![Go Reference](https://pkg.go.dev/badge/github.com/bububa/jiagu.svg)](https://pkg.go.dev/github.com/bububa/jiagu)
[![Go](https://github.com/bububa/jiagu/actions/workflows/go.yml/badge.svg)](https://github.com/bububa/jiagu/actions/workflows/go.yml)
[![goreleaser](https://github.com/bububa/jiagu/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/bububa/jiagu/actions/workflows/goreleaser.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/bububa/jiagu.svg)](https://github.com/bububa/jiagu)
[![GoReportCard](https://goreportcard.com/badge/github.com/bububa/jiagu)](https://goreportcard.com/report/github.com/bububa/jiagu)
[![GitHub license](https://img.shields.io/github/license/bububa/jiagu.svg)](https://github.com/bububa/jiagu/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/bububa/jiagu.svg)](https://GitHub.com/bububa/jiagu/releases/)

an NLP toolkit for Chinese.
This is a golang port from the [original python3 liberary](https://github.com/ownthink/Jiagu)
>>> Jiagu使用大规模语料训练而成。将提供中文分词、词性标注、命名实体识别、情感分析、知识图谱关系抽取、关键词抽取、文本摘要、新词发现、情感分析、文本聚类等常用自然语言处理功能。参考了各大工具优缺点制作，将Jiagu回馈给大家。

---

提供的功能有：
* 中文分词
* 词性标注
* 命名实体识别
* 知识图谱关系抽取
* 关键词提取
* 文本摘要
* 新词发现
* 情感分析
* 文本聚类
* 等等。。。。

---

## 使用方式
1. 快速上手：分词、词性标注、命名实体识别
```golang
import "github.com/bububa/jiagu"

func main() {
    // jiagu.init() // 可手动初始化，也可以动态初始化

    text := "厦门明天会不会下雨"

    words := jiagu.Seg(text) // 分词

    pos := jiagu.Pos(words) // 词性标注

    ner := jiagu.Ner(words) // 命名实体识别
}
```

2. 中文分词
```golang
import "github.com/bububa/jiagu"

func main() {
    text := "汉服和服装、维基图谱"

    words := jiagu.Seg(text)
    

    // fd, err := os.Open("user.dict")
    // defer fd.Close()
    // jiagu.LoadUserDict(fd) # 加载自定义字典，支持字典路径、字典列表形式。

    jiagu.AddVocabs([]string{"汉服和服装"})

    words := jiagu.Seg(text) // 自定义分词，字典分词模式有效
}
```

3. 知识图谱关系抽取
```golang
import "github.com/bububa/jiagu"

func main() {
    text := '姚明1980年9月12日出生于上海市徐汇区，祖籍江苏省苏州市吴江区震泽镇，前中国职业篮球运动员，司职中锋，现任中职联公司董事长兼总经理。'
    knowledge := jiagu.Knowledge(text)
}
```
训练数据：https://github.com/ownthink/KnowledgeGraphData

4. 关键词提取
```golang
import "github.com/bububa/jiagu"

func main() {
    text = `
    该研究主持者之一、波士顿大学地球与环境科学系博士陈池（音）表示，“尽管中国和印度国土面积仅占全球陆地的9%，但两国为这一绿化过程贡献超过三分之一。考虑到人口过多的国家一般存在对土地过度利用的问题，这个发现令人吃惊。”
    NASA埃姆斯研究中心的科学家拉玛·内曼尼（Rama Nemani）说，“这一长期数据能让我们深入分析地表绿化背后的影响因素。我们一开始以为，植被增加是由于更多二氧化碳排放，导致气候更加温暖、潮湿，适宜生长。”
    “MODIS的数据让我们能在非常小的尺度上理解这一现象，我们发现人类活动也作出了贡献。”
    NASA文章介绍，在中国为全球绿化进程做出的贡献中，有42%来源于植树造林工程，对于减少土壤侵蚀、空气污染与气候变化发挥了作用。
    据观察者网过往报道，2017年我国全国共完成造林736.2万公顷、森林抚育830.2万公顷。其中，天然林资源保护工程完成造林26万公顷，退耕还林工程完成造林91.2万公顷。京津风沙源治理工程完成造林18.5万公顷。三北及长江流域等重点防护林体系工程完成造林99.1万公顷。完成国家储备林建设任务68万公顷。
    `

    keywords := jiagu.keywords(text, 5) 
}
```

5. 文本摘要
```golang
import "github.com/bububa/jiagu"

func main() {
    summarize := jiagu.Summarize(text, 3) # 摘要
}
```

6. 新词发现
```golang
import "github.com/bububa/jiagu"

func main() {
    fd, err := os.Open("input.txt")
    defer fd.Close()

    words, err := jiagu.Findword(fd, 0, 0, 0) // 根据文本，利用信息熵做新词发现。
}
```

7. 情感分析
```golang
import "github.com/bububa/jiagu"

func main() {
    text := "很讨厌还是个懒鬼"
    words := jiagu.Seg(text)
    sentiment, probe := jiagu.Sentiment(words)
}
```

8. 文本聚类
```golang
import "github.com/bububa/jiagu"

func main() {
    docs := []string {
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
		return jiagu.Seg(txt)
	}
    cluster := jiagu.KmeansCluster(docs, tokenizer, 4)	

}
print(cluster)
```


## 附录
1. 词性标注说明
```text
n　　　普通名词
nt　 　时间名词
nd　 　方位名词
nl　 　处所名词
nh　 　人名
nhf　　姓
nhs　　名
ns　 　地名
nn 　　族名
ni 　　机构名
nz 　　其他专名
v　　 动词
vd　　趋向动词
vl　　联系动词
vu　　能愿动词
a　 　形容词
f　 　区别词
m　 　数词　　
q　 　量词
d　 　副词
r　 　代词
p　　 介词
c　 　连词
u　　 助词
e　 　叹词
o　 　拟声词
i　 　习用语
j　　 缩略语
h　　 前接成分
k　　 后接成分
g　 　语素字
x　 　非语素字
w　 　标点符号
ws　　非汉字字符串
wu　　其他未知的符号
```

2. 命名实体说明（采用BIO标记方式）
```text
B-PER、I-PER   人名
B-LOC、I-LOC   地名
B-ORG、I-ORG   机构名
```

## Reference：
- [Jiagu Python version](https://github.com/ownthink/Jiagu)
- [知识图谱训练数据](https://github.com/ownthink/KnowledgeGraphData)

