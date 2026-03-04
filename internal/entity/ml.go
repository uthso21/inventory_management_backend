package entities

// ============================================================================
// Demand Forecasting Models
// ============================================================================

// HistoricalSales represents a single historical sales data point
type HistoricalSales struct {
	Year      int     `json:"year"`
	Month     int     `json:"month"`
	Day       int     `json:"day"`
	UnitsSold float64 `json:"units_sold"`
}

// WeeklyPrediction represents a single weekly demand prediction
type WeeklyPrediction struct {
	Week           int     `json:"week"`
	Year           int     `json:"year"`
	Month          int     `json:"month"`
	Day            int     `json:"day"`
	PredictedUnits float64 `json:"predicted_units"`
}

// DemandForecastRequest represents the request for demand forecasting
type DemandForecastRequest struct {
	ProductName     string            `json:"product_name"`
	WarehouseName   string            `json:"warehouse_name"`
	HistoricalData  []HistoricalSales `json:"historical_data"`
	ForecastHorizon int               `json:"forecast_horizon,omitempty"`
}

// DemandForecastResponse represents the response from demand forecasting
type DemandForecastResponse struct {
	ProductName               string             `json:"product_name"`
	WarehouseName             string             `json:"warehouse_name"`
	WeeklyPredictions         []WeeklyPrediction `json:"weekly_predictions"`
	TotalUnitsNeeded          float64            `json:"total_units_needed"`
	ChangeFromPreviousPercent float64            `json:"change_from_previous_percent"`
	TrendDirection            string             `json:"trend_direction"`
	Volatility                string             `json:"volatility"`
	Confidence                float64            `json:"confidence"`
	AIExplanation             []string           `json:"ai_explanation"`
}

// ============================================================================
// Smart Reorder Models
// ============================================================================

// SmartReorderRequest represents the request for smart reorder
type SmartReorderRequest struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name,omitempty"`
	CurrentStock int     `json:"current_stock"`
	SafetyStock  int     `json:"safety_stock"`
	LeadTimeDays int     `json:"lead_time_days"`
	DailyDemand  float64 `json:"daily_demand"`
}

// SmartReorderResponse represents the response from smart reorder
type SmartReorderResponse struct {
	StockCoversDays    float64 `json:"stock_covers_days"`
	ReorderRecommended bool    `json:"reorder_recommended"`
	ReorderQuantity    int     `json:"reorder_quantity"`
	Urgency            string  `json:"urgency"`
	Confidence         float64 `json:"confidence"`
}

// ============================================================================
// Price Optimization Models
// ============================================================================

// PriceOptimizationRequest represents the request for price optimization
type PriceOptimizationRequest struct {
	ProductName     string  `json:"product_name"`
	CurrentPrice    float64 `json:"current_price"`
	COGSWeightedAvg float64 `json:"cogs_weighted_avg"`
	MarginPercent   float64 `json:"margin_percent"`
	SalesVelocity   float64 `json:"sales_velocity"`
}

// PriceOptimizationResponse represents the response from price optimization
type PriceOptimizationResponse struct {
	ProductName        string   `json:"product_name"`
	SuggestedPrice     float64  `json:"suggested_price"`
	SuggestedAction    string   `json:"suggested_action"`
	ProjectedMarginPct float64  `json:"projected_margin_pct"`
	Confidence         float64  `json:"confidence"`
	AIExplanation      []string `json:"ai_explanation"`
}

// ============================================================================
// Health Check
// ============================================================================

// MLHealthResponse represents health check response from ML service
type MLHealthResponse struct {
	Status string `json:"status"`
}
