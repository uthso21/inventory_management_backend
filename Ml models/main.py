from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional
import os

from react_agent import (
    ReactAgent,
    ProductContext as AgentProductContext,
    create_agent,
)


app = FastAPI(title="ML Models API")


# =============================================================================
# Request / Response Models
# =============================================================================

class ProductContext(BaseModel):
    product_id: str
    product_name: Optional[str] = None  # For market trends search
    product_description: Optional[str] = None  # Detailed description for AI context
    category: Optional[str] = None  # For market trends search
    shop_country: Optional[str] = None  # Country code (e.g., 'US', 'UK', 'BD') for regional trends
    # Demand forecasting context
    sales_history: Optional[list[dict]] = None  # [{"date": "2025-01-01", "qty": 10}, ...]
    history_months: Optional[int] = 6  # 3-12 months
    # Reordering context
    current_stock: Optional[int] = None
    safety_stock: Optional[int] = None
    lead_time_days: Optional[int] = None
    # Pricing context
    days_in_inventory: Optional[int] = None
    current_price: Optional[float] = None
    cost: Optional[float] = None


class AgentRequest(BaseModel):
    query: str
    context: ProductContext


class ToolResultResponse(BaseModel):
    tool: str
    success: bool
    data: dict
    explanation: str
    confidence: float  # Confidence score (0.0 - 1.0) for this prediction
    model_used: str  # Model/algorithm used for this prediction
    error: Optional[str] = None


class MarketTrendsResponse(BaseModel):
    product: Optional[str] = None
    category: Optional[str] = None
    country: Optional[str] = None
    sentiment: Optional[str] = None
    trend_direction: Optional[str] = None
    trends: Optional[list[dict]] = None
    error: Optional[str] = None


class AgentResponse(BaseModel):
    intent: str
    results: list[ToolResultResponse]
    market_trends: Optional[MarketTrendsResponse] = None
    final_answer: str
    errors: list[str]


# =============================================================================
# ReAct Agent Instance (with Gemini LLM if API key available)
# =============================================================================

# Set GOOGLE_API_KEY environment variable to enable Gemini LLM
api_key = os.getenv("GOOGLE_API_KEY")
agent = create_agent(use_llm=bool(api_key), api_key=api_key)


# =============================================================================
# Agent Endpoint
# =============================================================================

@app.post("/agent", response_model=AgentResponse)
def agent_endpoint(request: AgentRequest):
    """
    Single entry point for all ML queries.
    Uses ReAct agent with LangChain + Gemini + Web Scraping for market trends.
    """
    # Convert Pydantic model to dataclass
    ctx = AgentProductContext(
        product_id=request.context.product_id,
        product_name=request.context.product_name,
        product_description=request.context.product_description,
        category=request.context.category,
        shop_country=request.context.shop_country,
        sales_history=request.context.sales_history,
        history_months=request.context.history_months or 6,
        current_stock=request.context.current_stock,
        safety_stock=request.context.safety_stock,
        lead_time_days=request.context.lead_time_days,
        days_in_inventory=request.context.days_in_inventory,
        current_price=request.context.current_price,
        cost=request.context.cost,
    )

    # Run the agent
    response = agent.run(request.query, ctx)

    # Convert market trends if present
    market_trends = None
    if response.market_trends:
        market_trends = MarketTrendsResponse(
            product=response.market_trends.get("product"),
            category=response.market_trends.get("category"),
            country=response.market_trends.get("country"),
            sentiment=response.market_trends.get("sentiment"),
            trend_direction=response.market_trends.get("trend_direction"),
            trends=response.market_trends.get("trends"),
            error=response.market_trends.get("error"),
        )

    # Convert to API response
    return AgentResponse(
        intent=response.intent,
        results=[
            ToolResultResponse(
                tool=r.tool_name,
                success=r.success,
                data=r.data,
                explanation=r.explanation,
                confidence=r.data.get("confidence", 0.0),
                model_used=r.data.get("model_used", "unknown"),
                error=r.error,
            )
            for r in response.results
        ],
        market_trends=market_trends,
        final_answer=response.final_answer,
        errors=response.errors,
    )


# =============================================================================
# Health Check
# =============================================================================

@app.get("/health")
def health_check():
    return {
        "status": "healthy",
        "service": "ML Models API",
        "llm_enabled": bool(api_key),
    }
