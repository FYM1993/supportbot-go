package vectorstore

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"go.uber.org/zap"
)

// Document 文档结构
type Document struct {
	ID       string            // 文档唯一标识
	Content  string            // 文档内容
	Vector   []float64         // 文档向量
	Metadata map[string]string // 元数据（来源、标题等）
}

// SearchResult 搜索结果
type SearchResult struct {
	Document Document // 文档
	Score    float64  // 相似度得分（0-1，越高越相似）
	Distance float64  // 向量距离
}

// MemoryVectorStore 内存向量存储（简化版）
type MemoryVectorStore struct {
	documents map[string]*Document // 文档存储
	mu        sync.RWMutex         // 读写锁
	logger    *zap.Logger
}

// NewMemoryVectorStore 创建内存向量存储
func NewMemoryVectorStore(logger *zap.Logger) *MemoryVectorStore {
	return &MemoryVectorStore{
		documents: make(map[string]*Document),
		logger:    logger,
	}
}

// AddDocument 添加文档
func (s *MemoryVectorStore) AddDocument(doc Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if doc.ID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}

	if len(doc.Vector) == 0 {
		return fmt.Errorf("document vector cannot be empty")
	}

	s.documents[doc.ID] = &doc
	s.logger.Info("文档已添加", zap.String("id", doc.ID), zap.Int("dimension", len(doc.Vector)))
	return nil
}

// AddDocuments 批量添加文档
func (s *MemoryVectorStore) AddDocuments(docs []Document) error {
	for _, doc := range docs {
		if err := s.AddDocument(doc); err != nil {
			return err
		}
	}
	return nil
}

// Search 向量检索（返回 Top-K 最相似的文档）
func (s *MemoryVectorStore) Search(queryVector []float64, topK int, minScore float64) ([]SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(queryVector) == 0 {
		return nil, fmt.Errorf("query vector cannot be empty")
	}

	s.logger.Info("开始向量检索",
		zap.Int("docCount", len(s.documents)),
		zap.Int("topK", topK),
		zap.Float64("minScore", minScore))

	// 计算所有文档的相似度
	results := make([]SearchResult, 0, len(s.documents))
	for _, doc := range s.documents {
		score := cosineSimilarity(queryVector, doc.Vector)
		if score >= minScore {
			results = append(results, SearchResult{
				Document: *doc,
				Score:    score,
				Distance: 1 - score, // 余弦距离 = 1 - 余弦相似度
			})
		}
	}

	// 按相似度排序（降序）
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// 截取 Top-K
	if len(results) > topK {
		results = results[:topK]
	}

	s.logger.Info("检索完成",
		zap.Int("resultCount", len(results)),
		zap.Float64("topScore", getTopScore(results)))

	return results, nil
}

// GetDocument 获取文档
func (s *MemoryVectorStore) GetDocument(id string) (*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, ok := s.documents[id]
	if !ok {
		return nil, fmt.Errorf("document not found: %s", id)
	}
	return doc, nil
}

// DeleteDocument 删除文档
func (s *MemoryVectorStore) DeleteDocument(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.documents[id]; !ok {
		return fmt.Errorf("document not found: %s", id)
	}

	delete(s.documents, id)
	s.logger.Info("文档已删除", zap.String("id", id))
	return nil
}

// Count 获取文档数量
func (s *MemoryVectorStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.documents)
}

// Clear 清空所有文档
func (s *MemoryVectorStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.documents = make(map[string]*Document)
	s.logger.Info("向量存储已清空")
}

// cosineSimilarity 计算余弦相似度（0-1，越高越相似）
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	normA = math.Sqrt(normA)
	normB = math.Sqrt(normB)

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (normA * normB)
}

// getTopScore 获取最高得分（用于日志）
func getTopScore(results []SearchResult) float64 {
	if len(results) == 0 {
		return 0
	}
	return results[0].Score
}
