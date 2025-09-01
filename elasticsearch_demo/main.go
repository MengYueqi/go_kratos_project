package main

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"strconv"
	"time"
)

// createIndex 创建索引
func createIndex(client *elasticsearch.TypedClient) {
	resp, err := client.Indices.
		Create("my-review-1").
		Do(context.Background())
	if err != nil {
		fmt.Printf("create index failed, err:%v\n", err)
		return
	}
	fmt.Printf("index:%#v\n", resp.Index)
	fmt.Println(resp.Acknowledged)
}

// Review 评价数据
type Review struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userID"`
	Score       uint8     `json:"score"`
	Content     string    `json:"content"`
	Tags        []Tag     `json:"tags"`
	Status      int       `json:"status"`
	PublishTime time.Time `json:"publishDate"`
}

// Tag 评价标签
type Tag struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
}

// getDocument 获取文档
func getDocument(client *elasticsearch.TypedClient, id string) {
	resp, err := client.Get("my-review-1", id).
		Do(context.Background())
	if err != nil {
		fmt.Printf("get document by id failed, err:%v\n", err)
		return
	}
	fmt.Printf("fileds:%s\n", resp.Source_)
}

// indexDocument 索引文档
func indexDocument(client *elasticsearch.TypedClient) {
	// 定义 document 结构体对象
	d1 := Review{
		ID:      4,
		UserID:  147982601,
		Score:   2,
		Content: "这是一个好评！",
		Tags: []Tag{
			{1000, "好评"},
			{1100, "物超所值"},
			{9000, "有图"},
		},
		Status:      2,
		PublishTime: time.Now(),
	}

	// 添加文档
	resp, err := client.Index("my-review-1").
		Id(strconv.FormatInt(d1.ID, 10)).
		Document(d1).
		Do(context.Background())
	if err != nil {
		fmt.Printf("indexing document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%#v\n", resp.Result)
}

// aggregationDemo 聚合
func aggregationDemo(client *elasticsearch.TypedClient) {
	avgScoreAgg, err := client.Search().
		Index("my-review-1").
		Request(
			&search.Request{
				Size: some.Int(0),
				Aggregations: map[string]types.Aggregations{
					"avg_score": { // 将所有文档的 score 的平均值聚合为 avg_score
						Avg: &types.AverageAggregation{
							Field: some.String("Score"),
						},
					},
				},
			},
		).Do(context.Background())
	if err != nil {
		fmt.Printf("aggregation failed, err:%v\n", err)
		return
	}
	fmt.Printf("avgScore:%#v\n", avgScoreAgg.Aggregations["avg_score"])
}

func main() {
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
		return
	}

	//indexDocument(client)
	//getDocument(client, "4")

	aggregationDemo(client)
}
