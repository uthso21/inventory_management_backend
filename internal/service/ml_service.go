package usecases

import (
"bytes"
"context"
"encoding/json"
"errors"
"fmt"
"io"
"net/http"
"os"
"time"

entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

var (
ErrMLServiceUnavailable = errors.New("ML service is unavailable")
ErrInvalidMLResponse    = errors.New("invalid response from ML service")
ErrValidation           = errors.New("validation error")
)

type MLService interface {
HealthCheck(ctx context.Context) (*entities.MLHealthResponse, error)
GetDemandForecast(ctx context.Context, req *entities.DemandForecastRequest) (*entities.DemandForecastResponse, error)
GetSmartReorder(ctx context.Context, req *entities.SmartReorderRequest) (*entities.SmartReorderResponse, error)
GetPriceOptimization(ctx context.Context, req *entities.PriceOptimizationRequest) (*entities.PriceOptimizationResponse, error)
}

type mlService struct {
baseURL    string
httpClient *http.Client
}

func NewMLService() MLService {
baseURL := os.Getenv("ML_SERVICE_URL")
if baseURL == "" {
baseURL = "http://localhost:8000"
}
return &mlService{
baseURL: baseURL,
httpClient: &http.Client{
Timeout: 30 * time.Second,
},
}
}

func (s *mlService) doRequest(ctx context.Context, method, endpoint string, reqBody interface{}, result interface{}) error {
url := fmt.Sprintf("%s%s", s.baseURL, endpoint)

var body io.Reader
if reqBody != nil {
data, err := json.Marshal(reqBody)
if err != nil {
return fmt.Errorf("failed to marshal request: %w", err)
}
body = bytes.NewBuffer(data)
}

httpReq, err := http.NewRequestWithContext(ctx, method, url, body)
if err != nil {
return fmt.Errorf("failed to create request: %w", err)
}

if body != nil {
httpReq.Header.Set("Content-Type", "application/json")
}

resp, err := s.httpClient.Do(httpReq)
if err != nil {
return ErrMLServiceUnavailable
}
defer resp.Body.Close()

respBody, err := io.ReadAll(resp.Body)
if err != nil {
return fmt.Errorf("failed to read response: %w", err)
}

if resp.StatusCode != http.StatusOK {
return fmt.Errorf("ML service error: %s", string(respBody))
}

if err := json.Unmarshal(respBody, result); err != nil {
return ErrInvalidMLResponse
}

return nil
}

func (s *mlService) HealthCheck(ctx context.Context) (*entities.MLHealthResponse, error) {
var result entities.MLHealthResponse
if err := s.doRequest(ctx, http.MethodGet, "/health", nil, &result); err != nil {
return nil, err
}
return &result, nil
}

func (s *mlService) GetDemandForecast(ctx context.Context, req *entities.DemandForecastRequest) (*entities.DemandForecastResponse, error) {
// Validation
if req.ProductName == "" {
return nil, fmt.Errorf("%w: product_name is required", ErrValidation)
}
if req.WarehouseName == "" {
return nil, fmt.Errorf("%w: warehouse_name is required", ErrValidation)
}
if len(req.HistoricalData) < 2 {
return nil, fmt.Errorf("%w: at least 2 historical data points required", ErrValidation)
}

var result entities.DemandForecastResponse
if err := s.doRequest(ctx, http.MethodPost, "/demand-forecast", req, &result); err != nil {
return nil, err
}
return &result, nil
}

func (s *mlService) GetSmartReorder(ctx context.Context, req *entities.SmartReorderRequest) (*entities.SmartReorderResponse, error) {
// Validation
if req.ProductID == "" {
return nil, fmt.Errorf("%w: product_id is required", ErrValidation)
}

var result entities.SmartReorderResponse
if err := s.doRequest(ctx, http.MethodPost, "/smart-reorder", req, &result); err != nil {
return nil, err
}
return &result, nil
}

func (s *mlService) GetPriceOptimization(ctx context.Context, req *entities.PriceOptimizationRequest) (*entities.PriceOptimizationResponse, error) {
// Validation
if req.ProductName == "" {
return nil, fmt.Errorf("%w: product_name is required", ErrValidation)
}

var result entities.PriceOptimizationResponse
if err := s.doRequest(ctx, http.MethodPost, "/price-optimization", req, &result); err != nil {
return nil, err
}
return &result, nil
}
