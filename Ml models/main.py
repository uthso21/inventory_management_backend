from fastapi import FastAPI, HTTPException

from models import (
    DemandForecastRequest, DemandForecastResponse,
    SmartReorderRequest, SmartReorderResponse,
    PriceOptimizationRequest, PriceOptimizationResponse
)
from services import (
    calculate_demand_forecast,
    calculate_smart_reorder,
    calculate_price_optimization
)

app = FastAPI(title="ML Models API")


@app.get("/health")
def health_check():
    return {"status": "healthy"}


@app.post("/demand-forecast", response_model=DemandForecastResponse)
def demand_forecast(request: DemandForecastRequest):
    try:
        return calculate_demand_forecast(request)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/smart-reorder", response_model=SmartReorderResponse)
def smart_reorder(request: SmartReorderRequest):
    try:
        return calculate_smart_reorder(request)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/price-optimization", response_model=PriceOptimizationResponse)
def price_optimization(request: PriceOptimizationRequest):
    try:
        return calculate_price_optimization(request)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
