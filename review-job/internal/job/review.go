package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
	"review-job/internal/conf"
)

// 评价数据流处理任务

// 定义 kafka 中接收到的数据
type Msg struct {
	Type     string                   `json:"type"`
	Data     []map[string]interface{} `json:"data"`
	Table    string                   `json:"table"`
	IsDdl    bool                     `json:"isDdl"`
	Database string                   `json:"database"`
}

// 自定义执行 job，实现 transport.server
type JobWorker struct {
	kafkaReader *kafka.Reader
	esClient    *ESClient
	logger      *log.Helper
}

type ESClient struct {
	Client *elasticsearch.TypedClient
	index  string
}

func NewJobWorker(kafka *kafka.Reader, esClient *ESClient, logger log.Logger) *JobWorker {
	return &JobWorker{
		kafkaReader: kafka,
		esClient:    esClient,
		logger:      log.NewHelper(logger),
	}
}

func NewKafkaReader(cfg *conf.Kafka) (*kafka.Reader, error) {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupId,
		Topic:    cfg.Topic,
		MaxBytes: 10e6, // 10MB
	}), nil
}

func NewESClient(cfg_es *conf.Elasticsearch) (*ESClient, error) {
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: cfg_es.Addr,
	}

	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
		return &ESClient{}, err
	}
	return &ESClient{
		index:  cfg_es.Index,
		Client: client,
	}, nil
}

// Start 开始之后执行的程序
func (jw JobWorker) Start(ctx context.Context) error {
	jw.logger.Debugf("job worker starting")
	// 接收消息
	for {
		m, err := jw.kafkaReader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if err != nil {
			jw.logger.Error("read message error:", err)
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		// 将数据写入 es
		msg := new(Msg)
		err = json.Unmarshal(m.Value, msg)
		if err != nil {
			jw.logger.Error("json unmarshal error:", err)
			continue
		}

		if msg.Type == "INSERT" {
			// 往 ES 中新增文档
			for idx := range msg.Data {
				jw.indexDocument(msg.Data[idx])
			}
		} else {
			// 其他操作，在 ES 里更新文档
			for idx := range msg.Data {
				jw.updateDocument(msg.Data[idx])
			}
		}
	}
	return nil
}

// 退出方法
func (jw JobWorker) Stop(ctx context.Context) error {
	jw.logger.Debugf("job worker stopping")
	// 关闭 kafka 连接
	if err := jw.kafkaReader.Close(); err != nil {
		jw.logger.Error("failed to close reader:", err)
	}
	return nil
}

// indexDocument 索引文档
func (jw JobWorker) indexDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string)

	// 添加文档
	resp, err := jw.esClient.Client.Index(jw.esClient.index).
		Id(reviewID).
		Document(d).
		Do(context.Background())
	if err != nil {
		jw.logger.Debugf("indexing document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%#v\n", resp.Result)
}

// updateDocument 更新文档
func (jw JobWorker) updateDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string)

	resp, err := jw.esClient.Client.Update(jw.esClient.index, reviewID).
		Doc(d). // 使用结构体变量更新
		Do(context.Background())
	if err != nil {
		jw.logger.Errorf("update document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%v\n", resp.Result)
}
