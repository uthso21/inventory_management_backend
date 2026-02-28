package entities

// ProductContext represents the product data for ML analysis
type ProductContext struct {
	ProductID          string        `json:"product_id"`
	ProductName        string        `json:"product_name,omitempty"`
	ProductDescription string        `json:"product_description,omitempty"`
	Category           string        `json:"category,omitempty"`
	ShopCountry        string        `json:"shop_country,omitempty"`
	SalesHistory       []SalesRecord `json:"sales_history,omitempty"`
	HistoryMonths      int           `json:"history_months,omitempty"`
	CurrentStock       *int          `json:"current_stock,omitempty"`
	SafetyStock        *int          `json:"safety_stock,omitempty"`
	LeadTimeDays       *int          `json:"lead_time_days,omitempty"`
	DaysInInventory    *int          `json:"days_in_inventory,omitempty"`
	CurrentPrice       *float64      `json:"current_price,omitempty"`
	Cost               *float64      `json:"cost,omitempty"`
}

// SalesRecord represents a single sales data point
type SalesRecord struct {
	Date string `json:"date"`
	Qty  int    `json:"qty"`
}

// MLAgentRequest represents the request to the ML agent
type MLAgentRequest struct {
	Query   string         `json:"query"`
	Context ProductContext `json:"context"`
}

// MLAgentResponse represents the response from the ML agent
type MLAgentResponse struct {
	Intent       string            `json:"intent"`
	Results      []ToolResult      `json:"results"`
	MarketTrends *MarketTrends     `json:"market_trends,omitempty"`
	FinalAnswer  string            `json:"final_answer"`
	Errors       []string          `json:"errors"`
}

// ToolResult represents the result from a single ML tool
type ToolResult struct {
	Tool        string                 `json:"tool"`
	Success     bool                   `json:"success"`
	Data        map[string]interface{} `json:"data"`
	Explanation string                 `json:"explanation"`
	Confidence  float64                `json:"confidence"`   // Confidence score (0.0 - 1.0) for this prediction
	ModelUsed   string                 `json:"model_used"`   // Model/algorithm used for this prediction
	Error       *string                `json:"error,omitempty"`
}

// MarketTrends represents market trend data from web scraping
type MarketTrends struct {
	Product        string       `json:"product,omitempty"`
	Category       string       `json:"category,omitempty"`
	Country        string       `json:"country,omitempty"`
	Sentiment      string       `json:"sentiment,omitempty"`
	TrendDirection string       `json:"trend_direction,omitempty"`
	Trends         []TrendItem  `json:"trends,omitempty"`
	Error          *string      `json:"error,omitempty"`
}

// TrendItem represents a single trend headline
type TrendItem struct {
	Headline string `json:"headline"`
	URL      string `json:"url,omitempty"`
	Source   string `json:"source,omitempty"`
}

// MLQueryType defines the type of ML query
type MLQueryType string

const (
	QueryDemandForecast     MLQueryType = "demand_forecast"
	QuerySmartReorder       MLQueryType = "smart_reorder"
	QueryPricelistOptimize  MLQueryType = "pricelist_optimize"
	QueryFullAnalysis       MLQueryType = "full_analysis"
)

// Predefined query strings for frontend buttons
const (
	QueryStringDemandForecast    = "What's the demand forecast?"
	QueryStringSmartReorder      = "What's the smart reorder recommendation?"
	QueryStringPricelistOptimize = "What's the pricelist optimization?"
	QueryStringFullAnalysis      = "What's the full analysis?"
)
