package service

import (
	"fmt"
	"strings"

	"github.com/supportbot/supportbot-go/internal/client"
	"github.com/supportbot/supportbot-go/internal/vectorstore"
	"go.uber.org/zap"
)

// KnowledgeService 知识库服务
type KnowledgeService struct {
	embeddingClient *client.EmbeddingClient
	vectorStore     *vectorstore.MemoryVectorStore
	logger          *zap.Logger
}

// NewKnowledgeService 创建知识库服务
func NewKnowledgeService(embeddingClient *client.EmbeddingClient, vectorStore *vectorstore.MemoryVectorStore, logger *zap.Logger) *KnowledgeService {
	return &KnowledgeService{
		embeddingClient: embeddingClient,
		vectorStore:     vectorStore,
		logger:          logger,
	}
}

// AddKnowledge 添加知识（文本 → 向量化 → 存储）
func (s *KnowledgeService) AddKnowledge(id, content string, metadata map[string]string) error {
	s.logger.Info("添加知识", zap.String("id", id), zap.Int("length", len(content)))

	// 1. 获取文本向量
	vector, err := s.embeddingClient.GetEmbedding(content)
	if err != nil {
		return fmt.Errorf("向量化失败: %w", err)
	}

	// 2. 存储到向量数据库
	doc := vectorstore.Document{
		ID:       id,
		Content:  content,
		Vector:   vector,
		Metadata: metadata,
	}

	if err := s.vectorStore.AddDocument(doc); err != nil {
		return fmt.Errorf("存储失败: %w", err)
	}

	return nil
}

// AddKnowledgeBatch 批量添加知识
func (s *KnowledgeService) AddKnowledgeBatch(items []KnowledgeItem) error {
	s.logger.Info("批量添加知识", zap.Int("count", len(items)))

	// 1. 批量获取向量
	texts := make([]string, len(items))
	for i, item := range items {
		texts[i] = item.Content
	}

	vectors, err := s.embeddingClient.GetEmbeddings(texts)
	if err != nil {
		return fmt.Errorf("批量向量化失败: %w", err)
	}

	// 2. 批量存储
	docs := make([]vectorstore.Document, len(items))
	for i, item := range items {
		docs[i] = vectorstore.Document{
			ID:       item.ID,
			Content:  item.Content,
			Vector:   vectors[i],
			Metadata: item.Metadata,
		}
	}

	if err := s.vectorStore.AddDocuments(docs); err != nil {
		return fmt.Errorf("批量存储失败: %w", err)
	}

	return nil
}

// SearchKnowledge 检索知识（查询 → 向量化 → 相似度搜索）
func (s *KnowledgeService) SearchKnowledge(query string, topK int, minScore float64) ([]vectorstore.SearchResult, error) {
	s.logger.Info("检索知识", zap.String("query", query), zap.Int("topK", topK))

	// 1. 查询向量化
	queryVector, err := s.embeddingClient.GetQueryEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("查询向量化失败: %w", err)
	}

	// 2. 向量检索
	results, err := s.vectorStore.Search(queryVector, topK, minScore)
	if err != nil {
		return nil, fmt.Errorf("向量检索失败: %w", err)
	}

	return results, nil
}

// BuildContext 构建 RAG 上下文（将检索结果组合成文本）
func (s *KnowledgeService) BuildContext(results []vectorstore.SearchResult) string {
	if len(results) == 0 {
		return "未找到相关知识"
	}

	var builder strings.Builder
	builder.WriteString("参考知识库：\n\n")

	for i, result := range results {
		builder.WriteString(fmt.Sprintf("【知识片段 %d】(相似度: %.2f)\n", i+1, result.Score))
		builder.WriteString(result.Document.Content)
		builder.WriteString("\n\n")
	}

	return builder.String()
}

// KnowledgeItem 知识条目
type KnowledgeItem struct {
	ID       string
	Content  string
	Metadata map[string]string
}

// InitDefaultKnowledge 初始化默认知识库（电商场景）
func (s *KnowledgeService) InitDefaultKnowledge() error {
	s.logger.Info("初始化默认知识库...")

	knowledgeBase := []KnowledgeItem{
		{
			ID:      "product-return-policy",
			Content: "商品退货政策：购买后7天内，如商品未拆封、未使用，可无理由退货。需保持商品完好，包装齐全。退货运费由买家承担，特殊商品（如生鲜、定制品）不支持退货。",
			Metadata: map[string]string{
				"category": "退货政策",
				"source":   "官方政策",
			},
		},
		{
			ID:      "product-warranty",
			Content: "商品质保说明：电子产品享受1年免费质保服务，质保期内非人为损坏可免费维修或更换。质保期后提供付费维修服务。需提供购买凭证和产品序列号。",
			Metadata: map[string]string{
				"category": "质保服务",
				"source":   "售后手册",
			},
		},
		{
			ID:      "coupon-usage",
			Content: "优惠券使用规则：优惠券有使用期限，过期自动作废。单笔订单仅限使用一张优惠券，不与其他活动叠加。满减券需满足最低消费金额。部分特价商品不可使用优惠券。",
			Metadata: map[string]string{
				"category": "优惠活动",
				"source":   "活动规则",
			},
		},
		{
			ID:      "shipping-info",
			Content: "物流配送说明：订单支付成功后48小时内发货，节假日顺延。提供顺丰、圆通、韵达等多家物流选择。偏远地区可能需要额外3-5天配送时间。可在订单详情页查询物流信息。",
			Metadata: map[string]string{
				"category": "物流配送",
				"source":   "配送指南",
			},
		},
		{
			ID:      "payment-methods",
			Content: "支付方式说明：支持微信支付、支付宝、银行卡支付。支持花呗分期、信用卡分期付款。部分大额订单支持货到付款。支付完成后订单立即生效。",
			Metadata: map[string]string{
				"category": "支付方式",
				"source":   "支付指南",
			},
		},
		{
			ID:      "member-benefits",
			Content: "会员权益说明：普通会员享95折优惠，银卡会员9折，金卡会员85折，钻石会员8折。会员可积累积分，100积分可抵扣1元。会员生日当月享额外95折优惠。",
			Metadata: map[string]string{
				"category": "会员服务",
				"source":   "会员手册",
			},
		},
		{
			ID:      "product-maintenance",
			Content: "商品保养建议：电子产品避免潮湿环境，定期清洁灰尘。服装类建议干洗或手洗，避免暴晒。家具定期打蜡保养，避免阳光直射。具体保养方法请参考产品说明书。",
			Metadata: map[string]string{
				"category": "保养指南",
				"source":   "使用手册",
			},
		},
		{
			ID:      "promotion-double11",
			Content: "双11大促活动：11月1日-11日全场5折起，每日10点、20点整点秒杀。前1000名下单用户赠送50元无门槛券。购物满500元抽奖，最高可得iPhone 15。活动商品不支持退换货。",
			Metadata: map[string]string{
				"category": "促销活动",
				"source":   "活动页面",
				"valid_date": "2024-11-01 至 2024-11-11",
			},
		},
	}

	return s.AddKnowledgeBatch(knowledgeBase)
}

