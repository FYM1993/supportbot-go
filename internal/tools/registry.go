package tools

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Registry 工具注册中心
type Registry struct {
	tools  map[string]*Tool
	mu     sync.RWMutex
	logger *zap.Logger
}

// NewRegistry 创建工具注册中心
func NewRegistry(logger *zap.Logger) *Registry {
	return &Registry{
		tools:  make(map[string]*Tool),
		logger: logger,
	}
}

// Register 注册工具
func (r *Registry) Register(tool *Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tool.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool already registered: %s", tool.Name)
	}

	r.tools[tool.Name] = tool
	r.logger.Info("工具已注册", zap.String("name", tool.Name))
	return nil
}

// Get 获取工具
func (r *Registry) Get(name string) (*Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, ok := r.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool, nil
}

// List 列出所有工具
func (r *Registry) List() []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]*Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// GetFunctionDefs 获取所有工具的 Function 定义（用于 LLM）
func (r *Registry) GetFunctionDefs() []map[string]interface{} {
	tools := r.List()
	defs := make([]map[string]interface{}, len(tools))
	for i, tool := range tools {
		defs[i] = tool.ToFunctionDef()
	}
	return defs
}

// Execute 执行工具调用
func (r *Registry) Execute(toolCall ToolCall) (interface{}, error) {
	r.logger.Info("执行工具调用",
		zap.String("tool", toolCall.Function.Name),
		zap.String("callId", toolCall.ID))

	// 获取工具
	tool, err := r.Get(toolCall.Function.Name)
	if err != nil {
		return nil, err
	}

	// 解析参数
	params, err := toolCall.ParseArguments()
	if err != nil {
		return nil, err
	}

	// 执行工具
	result, err := tool.Execute(params)
	if err != nil {
		r.logger.Error("工具执行失败",
			zap.String("tool", toolCall.Function.Name),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("工具执行成功",
		zap.String("tool", toolCall.Function.Name))

	return result, nil
}

// Count 获取注册的工具数量
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tools)
}
