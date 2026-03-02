from pydantic import BaseModel
from typing import Optional, Literal


class WeeklyPrediction(BaseModel):
    week: int
    year: int
    month: int
    day: int
    predicted_units: float


class HistoricalSales(BaseModel):
    year: int
    month: int
    day: int
    units_sold: float


class DemandForecastRequest(BaseModel):
    product_name: str
    warehouse_name: str
    historical_data: list[HistoricalSales]
    forecast_horizon: int = 4


class DemandForecastResponse(BaseModel):
    product_name: str
    warehouse_name: str
    weekly_predictions: list[WeeklyPrediction]
    total_units_needed: float
    change_from_previous_percent: float
    trend_direction: str
    volatility: str
    confidence: float
    ai_explanation: list[str]


class SmartReorderRequest(BaseModel):
    product_id: str
    product_name: Optional[str] = None
    current_stock: int
    safety_stock: int
    lead_time_days: int
    daily_demand: float


class SmartReorderResponse(BaseModel):
    stock_covers_days: float
    reorder_recommended: bool
    reorder_quantity: int
    urgency: str
    confidence: float


class PriceOptimizationRequest(BaseModel):
    product_name: str
    current_price: float
    cogs_weighted_avg: float
    margin_percent: float
    sales_velocity: float


class PriceOptimizationResponse(BaseModel):
    product_name: str
    suggested_price: float
    suggested_action: str
    projected_margin_pct: float
    confidence: float
    ai_explanation: list[str]
