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
	ErrMLServiceTimeout     = errors.New("ML service request timed out")
	ErrInvalidMLResponse    = errors.New("invalid response from ML service")
)

// MLAgentService defines the interface for ML agent operations
type MLAgentService interface {
	// ProcessQuery sends a query to the ML microservice and returns the response
	ProcessQuery(ctx context.Context, req *entities.MLAgentRequest) (*entities.MLAgentResponse, error)

	// GetDemandForecast is a convenience method for demand forecasting
	GetDemandForecast(ctx context.Context, productCtx *entities.ProductContext) (*entities.MLAgentResponse, error)

	// GetSmartReorder is a convenience method for smart reorder recommendations
	GetSmartReorder(ctx context.Context, productCtx *entities.ProductContext) (*entities.MLAgentResponse, error)

	// GetPricelistOptimization is a convenience method for pricelist optimization
	GetPricelistOptimization(ctx context.Context, productCtx *entities.ProductContext) (*entities.MLAgentResponse, error)

	// GetFullAnalysis runs all three tools
	GetFullAnalysis(ctx context.Context, productCtx *entities.ProductContext) (*entities.MLAgentResponse, error)

	// HealthCheck checks if the ML service is available
	HealthCheck(ctx context.Context) (bool, error)
}

// mlAgentService is the concrete implementation of MLAgentService
type mlAgentService struct {
	baseURL    string
	httpClient *http.Client
}

// MLAgentConfig holds configuration for the ML agent service
type MLAgentConfig struct {
	BaseURL string
	Timeout time.Duration
}

// DefaultMLAgentConfig returns default configuration
func DefaultMLAgentConfig() MLAgentConfig {
	baseURL := os.Getenv("ML_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000" // Default FastAPI URL
	}

	return MLAgentConfig{
		BaseURL: baseURL,
		Timeout: 60 * time.Second, // ML operations can take time
	}
}

// NewMLAgentService creates a new instance of MLAgentService
func NewMLAgentService(config MLAgentConfig) MLAgentService {
	return &mlAgentService{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// NewMLAgentServiceWithDefaults creates a new MLAgentService with default config
func NewMLAgentServiceWithDefaults() MLAgentService {
	return NewMLAgentService(DefaultMLAgentConfig())
}

func (s *mlAgentService) ProcessQuery(ctx context.Context, req *entities.MLAgentRequest) (*entities.MLAgentResponse, error) {
	// Validate request
	if req.Query == "" {
		return nil, ErrInvalidInput
	}
	if req.Context.ProductID == "" {
		return nil, errors.New("product_id is required")
	}

	// Prepare request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/agent", s.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, ErrMLServiceTimeout
		}
		return nil, fmt.Errorf("%w: %v", ErrMLServiceUnavailable, err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var mlResp entities.MLAgentResponse
	if err := json.Unmarshal(respBody, &mlResp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidMLResponse, err)
	}

	return &mlResp, nil
}

func (s *mlAgentService) GetDemandForecast(ctx context.Context, productCtx *entities.ProductContext) (*entities.MLAgentResponse, error) {
	req := &entities.MLAgentRequest{
		Query:   entities.QueryStringDemandForecast,
		Context: *productCtx,
	}
	return s.ProcessQuery(ctx, req)
}

func (s *mlAgentService) GetSmartReorder(ctx context.Context, productCtx *entities.ProductContext) (*entities.MLAgentResponse, error) {
	req := &entities.MLAgentRequest{
		Query:   entities.QueryStringSmartReorder,
		Context: *productCtx,
	}
	return s.ProcessQuery(ctx, req)
}

func (s *mlAgentService) GetPricelistOptimization(ctx context.Context, productCtx *entities.ProductContext) (*entities.MLAgentResponse, error) {
	req := &entities.MLAgentRequest{
		Query:   entities.QueryStringPricelistOptimize,
		Context: *productCtx,
	}
	return s.ProcessQuery(ctx, req)
}

func (s *mlAgentService) GetFullAnalysis(ctx context.Context, productCtx *entities.ProductContext) (*entities.MLAgentResponse, error) {
	req := &entities.MLAgentRequest{
		Query:   entities.QueryStringFullAnalysis,
		Context: *productCtx,
	}
	return s.ProcessQuery(ctx, req)
}

func (s *mlAgentService) HealthCheck(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/health", s.baseURL)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrMLServiceUnavailable, err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
