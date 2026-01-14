package tools

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// RegisterBuiltinTools 注册内置工具（模拟电商业务）
func RegisterBuiltinTools(registry *Registry, logger *zap.Logger) error {
	logger.Info("注册内置工具...")

	// 1. 商品详情查询
	productDetailTool := &Tool{
		Name:        "get_product_detail",
		Description: "查询商品详细信息，包括名称、价格、规格、库存等",
		Parameters: ParameterSchema{
			Type: "object",
			Properties: map[string]Property{
				"product_id": {
					Type:        "string",
					Description: "商品ID，例如：30001, 30002, 30003",
				},
			},
			Required: []string{"product_id"},
		},
		Handler: func(params map[string]interface{}) (interface{}, error) {
			productID, ok := params["product_id"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid product_id")
			}

			// 模拟数据库查询
			products := map[string]map[string]interface{}{
				"30001": {
					"id":          "30001",
					"name":        "iPhone 15 Pro Max",
					"price":       9999.00,
					"stock":       128,
					"color":       "钛金色",
					"storage":     "256GB",
					"description": "A17 Pro 芯片，钛金属设计，支持 Action 按钮",
				},
				"30002": {
					"id":          "30002",
					"name":        "MacBook Pro 16\"",
					"price":       19999.00,
					"stock":       45,
					"chip":        "M3 Max",
					"memory":      "32GB",
					"storage":     "1TB SSD",
					"description": "专业性能，适合创作者",
				},
				"30003": {
					"id":          "30003",
					"name":        "AirPods Pro 2",
					"price":       1999.00,
					"stock":       320,
					"features":    "主动降噪、自适应音频、空间音频",
					"description": "新一代降噪耳机",
				},
			}

			product, ok := products[productID]
			if !ok {
				return map[string]interface{}{
					"error": "商品不存在",
				}, nil
			}

			return product, nil
		},
	}

	// 2. 订单详情查询
	orderDetailTool := &Tool{
		Name:        "get_order_detail",
		Description: "查询订单详细信息，包括订单状态、商品、金额、物流等",
		Parameters: ParameterSchema{
			Type: "object",
			Properties: map[string]Property{
				"order_id": {
					Type:        "string",
					Description: "订单号，例如：20240101001",
				},
			},
			Required: []string{"order_id"},
		},
		Handler: func(params map[string]interface{}) (interface{}, error) {
			orderID, ok := params["order_id"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid order_id")
			}

			// 模拟数据库查询
			orders := map[string]map[string]interface{}{
				"20240101001": {
					"order_id":     "20240101001",
					"status":       "已发货",
					"product":      "iPhone 15 Pro Max",
					"quantity":     1,
					"total_amount": 9999.00,
					"create_time":  "2024-01-01 10:30:00",
					"ship_time":    "2024-01-01 15:00:00",
					"tracking_no":  "SF1234567890",
					"address":      "北京市朝阳区xxx路xxx号",
				},
				"20240101002": {
					"order_id":     "20240101002",
					"status":       "待发货",
					"product":      "AirPods Pro 2",
					"quantity":     2,
					"total_amount": 3998.00,
					"create_time":  "2024-01-02 14:20:00",
					"address":      "上海市浦东新区xxx路xxx号",
				},
			}

			order, ok := orders[orderID]
			if !ok {
				return map[string]interface{}{
					"error": "订单不存在",
				}, nil
			}

			return order, nil
		},
	}

	// 3. 物流查询
	shippingTrackingTool := &Tool{
		Name:        "get_shipping_tracking",
		Description: "查询订单物流信息，包括当前位置、配送进度、预计送达时间",
		Parameters: ParameterSchema{
			Type: "object",
			Properties: map[string]Property{
				"order_id": {
					Type:        "string",
					Description: "订单号",
				},
			},
			Required: []string{"order_id"},
		},
		Handler: func(params map[string]interface{}) (interface{}, error) {
			orderID, ok := params["order_id"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid order_id")
			}

			// 模拟物流信息
			tracking := map[string]interface{}{
				"order_id":      orderID,
				"tracking_no":   "SF1234567890",
				"carrier":       "顺丰速运",
				"current_location": "北京分拨中心",
				"status":        "运输中",
				"estimated_delivery": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				"tracking_info": []map[string]string{
					{
						"time":     "2024-01-02 08:00:00",
						"location": "北京分拨中心",
						"status":   "已签收",
					},
					{
						"time":     "2024-01-01 20:30:00",
						"location": "北京中转站",
						"status":   "运输中",
					},
					{
						"time":     "2024-01-01 15:00:00",
						"location": "北京发货仓",
						"status":   "已发货",
					},
				},
			}

			return tracking, nil
		},
	}

	// 4. 商品库存查询
	productAvailabilityTool := &Tool{
		Name:        "get_product_availability",
		Description: "查询商品库存和配送信息",
		Parameters: ParameterSchema{
			Type: "object",
			Properties: map[string]Property{
				"product_id": {
					Type:        "string",
					Description: "商品ID",
				},
				"region": {
					Type:        "string",
					Description: "配送地区，例如：北京、上海、广州",
				},
			},
			Required: []string{"product_id"},
		},
		Handler: func(params map[string]interface{}) (interface{}, error) {
			productID, ok := params["product_id"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid product_id")
			}

			region, _ := params["region"].(string)
			if region == "" {
				region = "北京"
			}

			// 模拟库存查询
			return map[string]interface{}{
				"product_id":         productID,
				"stock":              128,
				"available":          true,
				"region":             region,
				"estimated_delivery": "明天送达",
				"shipping_fee":       0, // 包邮
			}, nil
		},
	}

	// 注册所有工具
	tools := []*Tool{
		productDetailTool,
		orderDetailTool,
		shippingTrackingTool,
		productAvailabilityTool,
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}

	logger.Info("内置工具注册完成", zap.Int("count", len(tools)))
	return nil
}

