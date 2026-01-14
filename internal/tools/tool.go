package tools

import (
	"encoding/json"
	"fmt"
)

// Tool 工具定义（类似 OpenAI Function Calling）
type Tool struct {
	Name        string          `json:"name"`                  // 工具名称
	Description string          `json:"description"`            // 工具描述
	Parameters  ParameterSchema `json:"parameters"`             // 参数定义
	Handler     ToolHandler     `json:"-"`                      // 工具处理函数（不序列化）
}

// ParameterSchema JSON Schema 格式的参数定义
type ParameterSchema struct {
	Type       string              `json:"type"`        // "object"
	Properties map[string]Property `json:"properties"`  // 参数属性
	Required   []string            `json:"required"`    // 必需参数
}

// Property 参数属性
type Property struct {
	Type        string   `json:"type"`                   // string, number, boolean, array, object
	Description string   `json:"description"`             // 参数描述
	Enum        []string `json:"enum,omitempty"`          // 枚举值
}

// ToolHandler 工具处理函数
type ToolHandler func(params map[string]interface{}) (interface{}, error)

// ToolCall LLM 返回的工具调用请求
type ToolCall struct {
	ID       string                 `json:"id"`        // 调用 ID
	Type     string                 `json:"type"`      // "function"
	Function ToolCallFunction       `json:"function"`
}

// ToolCallFunction 函数调用详情
type ToolCallFunction struct {
	Name      string                 `json:"name"`
	Arguments string                 `json:"arguments"` // JSON 字符串
}

// ToolResult 工具执行结果
type ToolResult struct {
	ToolCallID string      `json:"tool_call_id"`
	Result     interface{} `json:"result"`
	Error      string      `json:"error,omitempty"`
}

// ParseArguments 解析参数
func (tc *ToolCall) ParseArguments() (map[string]interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &params); err != nil {
		return nil, fmt.Errorf("解析参数失败: %w", err)
	}
	return params, nil
}

// Execute 执行工具
func (t *Tool) Execute(params map[string]interface{}) (interface{}, error) {
	if t.Handler == nil {
		return nil, fmt.Errorf("tool handler not implemented: %s", t.Name)
	}
	return t.Handler(params)
}

// ToFunctionDef 转换为 LLM Function 定义格式
func (t *Tool) ToFunctionDef() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.Name,
		"description": t.Description,
		"parameters":  t.Parameters,
	}
}

