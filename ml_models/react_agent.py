"""
ReAct Agent for Inventory Management ML Features

This agent implements the ReAct (Reason + Act) pattern using LangChain and Gemini:
1. Receive user query + product context
2. Use Gemini LLM to reason about which tool(s) to use
3. Execute tool(s) including web scraping for trends
4. Observe results
5. Synthesize final response with explainability

Tools:
- Demand Forecasting Engine
- Smart Reordering
- Pricelist Optimization
- Market Trends Scraper (web scraping for recent trends)
"""

import json
import os
import re
import requests
from typing import Optional, Callable, Any
from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path

# Load environment variables from .env file
try:
    from dotenv import load_dotenv
    env_path = Path(__file__).parent / ".env"
    load_dotenv(env_path)
except ImportError:
    pass

try:
    from bs4 import BeautifulSoup
    BS4_AVAILABLE = True
except ImportError:
    BS4_AVAILABLE = False

# LangChain imports - handle newer versions
try:
    from langchain_google_genai import ChatGoogleGenerativeAI
    from langchain_core.messages import HumanMessage, SystemMessage
    from langchain_core.tools import Tool
    LANGCHAIN_AVAILABLE = True
except ImportError:
    LANGCHAIN_AVAILABLE = False
    ChatGoogleGenerativeAI = None
    HumanMessage = None
    SystemMessage = None
    Tool = None


# =============================================================================
# Configuration
# =============================================================================

# Set your Gemini API key in environment variable: GOOGLE_API_KEY
GEMINI_MODEL = "gemini-1.5-flash"  # or "gemini-1.5-pro" for better quality


# =============================================================================
# Data Models
# =============================================================================

class AgentState(Enum):
    THINKING = "thinking"
    ACTING = "acting"
    OBSERVING = "observing"
    DONE = "done"
    ERROR = "error"


@dataclass
class ProductContext:
    """Product data passed from Go backend."""
    product_id: str
    product_name: Optional[str] = None
    product_description: Optional[str] = None  # Detailed product description for AI context
    category: Optional[str] = None
    shop_country: Optional[str] = None  # Country code (e.g., 'US', 'UK', 'BD') for regional trends
    # Demand forecasting
    sales_history: Optional[list[dict]] = None
    history_months: int = 6
    # Reordering
    current_stock: Optional[int] = None
    safety_stock: Optional[int] = None
    lead_time_days: Optional[int] = None
    # Pricing
    days_in_inventory: Optional[int] = None
    current_price: Optional[float] = None
    cost: Optional[float] = None


@dataclass
class ToolResult:
    """Result from a tool execution."""
    tool_name: str
    success: bool
    data: dict
    explanation: str
    error: Optional[str] = None


@dataclass
class AgentStep:
    """Single step in the ReAct loop."""
    thought: str
    action: Optional[str] = None
    action_input: Optional[dict] = None
    observation: Optional[str] = None


@dataclass
class AgentResponse:
    """Final response from the agent."""
    query: str
    intent: str
    steps: list[AgentStep] = field(default_factory=list)
    results: list[ToolResult] = field(default_factory=list)
    final_answer: str = ""
    market_trends: Optional[dict] = None
    errors: list[str] = field(default_factory=list)


# =============================================================================
# Web Scraping for Market Trends
# =============================================================================

class MarketTrendsScraper:
    """
    Scrapes web sources for recent market trends related to products.
    Uses multiple sources for comprehensive trend data.
    """

    USER_AGENT = (
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
        "(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    )

    TREND_SOURCES = [
        {
            "name": "Google Trends (via search)",
            "url": "https://trends.google.com/trends/explore?q={query}",
            "type": "reference",
        },
    ]

    def __init__(self):
        self.session = requests.Session()
        self.session.headers.update({"User-Agent": self.USER_AGENT})

    def scrape_news_trends(self, product_name: str, category: str = None, country: str = None) -> dict:
        """
        Scrape recent news and trends for a product.
        Returns structured trend data.

        Args:
            product_name: Name of the product
            category: Product category
            country: Country code for regional trends (e.g., 'US', 'UK', 'BD')
        """
        # Build search query with regional context
        country_name = self._get_country_name(country) if country else ""
        search_query = f"{product_name} {category or ''} {country_name} market trends 2026".strip()
        trends_data = {
            "product": product_name,
            "category": category,
            "country": country,
            "search_query": search_query,
            "trends": [],
            "sentiment": "neutral",
            "trend_direction": "stable",
            "sources_checked": [],
            "error": None,
        }

        try:
            # Scrape from DuckDuckGo HTML (no API key needed)
            ddg_trends = self._scrape_duckduckgo(search_query)
            if ddg_trends:
                trends_data["trends"].extend(ddg_trends)
                trends_data["sources_checked"].append("DuckDuckGo")

            # Analyze sentiment from scraped headlines
            if trends_data["trends"]:
                sentiment, direction = self._analyze_sentiment(trends_data["trends"])
                trends_data["sentiment"] = sentiment
                trends_data["trend_direction"] = direction

        except Exception as e:
            trends_data["error"] = str(e)

        return trends_data

    def _get_country_name(self, country_code: str) -> str:
        """Convert country code to full name for better search results."""
        country_map = {
            "US": "United States",
            "UK": "United Kingdom",
            "GB": "United Kingdom",
            "BD": "Bangladesh",
            "IN": "India",
            "CA": "Canada",
            "AU": "Australia",
            "DE": "Germany",
            "FR": "France",
            "JP": "Japan",
            "CN": "China",
            "AE": "UAE",
            "SG": "Singapore",
            "MY": "Malaysia",
            "PK": "Pakistan",
            "NP": "Nepal",
        }
        return country_map.get(country_code.upper(), country_code) if country_code else ""

    def _scrape_duckduckgo(self, query: str) -> list[dict]:
        """Scrape DuckDuckGo HTML results for trend headlines."""
        trends = []
        if not BS4_AVAILABLE:
            return trends

        try:
            url = f"https://html.duckduckgo.com/html/?q={requests.utils.quote(query)}"
            response = self.session.get(url, timeout=10)

            if response.status_code == 200:
                soup = BeautifulSoup(response.text, "html.parser")
                results = soup.select(".result__title")[:5]  # Top 5 results

                for result in results:
                    link = result.select_one("a")
                    if link:
                        trends.append({
                            "headline": link.get_text(strip=True),
                            "url": link.get("href", ""),
                            "source": "web",
                        })
        except Exception as e:
            pass  # Gracefully handle scraping errors

        return trends

    def _scrape_google_news_rss(self, query: str) -> list[dict]:
        """Scrape Google News RSS feed for headlines."""
        trends = []
        if not BS4_AVAILABLE:
            return trends

        try:
            url = f"https://news.google.com/rss/search?q={requests.utils.quote(query)}&hl=en-US&gl=US&ceid=US:en"
            response = self.session.get(url, timeout=10)

            if response.status_code == 200:
                soup = BeautifulSoup(response.text, "xml")
                items = soup.find_all("item")[:5]

                for item in items:
                    title = item.find("title")
                    link = item.find("link")
                    pub_date = item.find("pubDate")

                    if title:
                        trends.append({
                            "headline": title.get_text(strip=True),
                            "url": link.get_text(strip=True) if link else "",
                            "date": pub_date.get_text(strip=True) if pub_date else "",
                            "source": "Google News",
                        })
        except Exception as e:
            pass

        return trends

    def _analyze_sentiment(self, trends: list[dict]) -> tuple[str, str]:
        """
        Simple sentiment analysis based on keyword matching.
        Returns (sentiment, trend_direction).
        """
        positive_keywords = [
            "growth", "surge", "rising", "increase", "boom", "demand",
            "popular", "trending", "hot", "best-selling", "record"
        ]
        negative_keywords = [
            "decline", "drop", "falling", "decrease", "slump", "slow",
            "weak", "struggling", "downturn", "shortage", "crisis"
        ]

        positive_count = 0
        negative_count = 0

        for trend in trends:
            headline = trend.get("headline", "").lower()
            for kw in positive_keywords:
                if kw in headline:
                    positive_count += 1
            for kw in negative_keywords:
                if kw in headline:
                    negative_count += 1

        if positive_count > negative_count:
            return "positive", "upward"
        elif negative_count > positive_count:
            return "negative", "downward"
        else:
            return "neutral", "stable"


# =============================================================================
# Tool Implementations
# =============================================================================

def demand_forecast_tool(ctx: ProductContext, market_trends: dict = None) -> ToolResult:
    """
    Demand Forecasting Engine
    Uses historical sales data + market trends to predict future demand.
    """
    history_months = ctx.history_months or 6
    sales_history = ctx.sales_history or []

    # Calculate forecast from historical data
    if sales_history:
        total_qty = sum(item.get("qty", 0) for item in sales_history)
        avg_daily = total_qty / max(len(sales_history), 1)
        forecast_weekly = avg_daily * 7

        if len(sales_history) >= 2:
            recent = sum(item.get("qty", 0) for item in sales_history[-7:])
            older = sum(item.get("qty", 0) for item in sales_history[:7])
            trend_pct = ((recent - older) / max(older, 1)) * 100
        else:
            trend_pct = 0
    else:
        forecast_weekly = 150
        trend_pct = 12.0

    # Adjust forecast based on market trends
    market_adjustment = 0
    market_insight = ""
    if market_trends and market_trends.get("trend_direction"):
        direction = market_trends["trend_direction"]
        if direction == "upward":
            market_adjustment = 10  # Increase forecast by 10%
            market_insight = "Market trends indicate rising demand."
        elif direction == "downward":
            market_adjustment = -10  # Decrease forecast by 10%
            market_insight = "Market trends indicate declining demand."
        else:
            market_insight = "Market trends are stable."

    adjusted_forecast = forecast_weekly * (1 + market_adjustment / 100)
    trend_direction = "upward" if trend_pct > 0 else "downward" if trend_pct < 0 else "stable"

    data = {
        "product_id": ctx.product_id,
        "forecast_units_per_week": round(adjusted_forecast, 1),
        "base_forecast": round(forecast_weekly, 1),
        "market_adjustment_pct": market_adjustment,
        "trend_percent": round(abs(trend_pct), 1),
        "trend_direction": trend_direction,
        "forecast_horizon_weeks": 4,
        "confidence": 0.85 if not market_trends else 0.88,
        "model_used": "moving_average_with_trends",
        "market_sentiment": market_trends.get("sentiment") if market_trends else None,
    }

    explanation = (
        f"Based on {history_months}-month sales history, base forecast is "
        f"{forecast_weekly:.1f} units/week with a {abs(trend_pct):.1f}% {trend_direction} trend. "
    )
    if market_insight:
        explanation += f"{market_insight} Adjusted forecast: {adjusted_forecast:.1f} units/week. "
    explanation += f"Confidence: {data['confidence']*100:.0f}%."

    return ToolResult(
        tool_name="demand_forecast",
        success=True,
        data=data,
        explanation=explanation
    )


def smart_reorder_tool(ctx: ProductContext, market_trends: dict = None) -> ToolResult:
    """
    Smart Reordering Engine
    Recommends reorder quantities based on stock levels, forecast, and market trends.
    """
    current_stock = ctx.current_stock or 0
    safety_stock = ctx.safety_stock or 50
    lead_time_days = ctx.lead_time_days or 14

    # Base daily demand (would come from forecast in production)
    daily_demand = 20

    # Adjust for market trends
    if market_trends and market_trends.get("trend_direction") == "upward":
        daily_demand *= 1.15  # Increase buffer for rising demand
        trend_note = "Increased buffer for rising market demand."
    elif market_trends and market_trends.get("trend_direction") == "downward":
        daily_demand *= 0.9  # Reduce for declining demand
        trend_note = "Reduced buffer for declining market demand."
    else:
        trend_note = ""

    stock_covers_days = current_stock / daily_demand if daily_demand > 0 else 999

    reorder_needed = (
        stock_covers_days < lead_time_days or
        current_stock <= safety_stock
    )

    if reorder_needed:
        reorder_qty = max(
            0,
            (lead_time_days * daily_demand) - current_stock + safety_stock
        )
    else:
        reorder_qty = 0

    # Determine urgency
    if stock_covers_days < 7:
        urgency = "critical"
    elif stock_covers_days < lead_time_days:
        urgency = "high"
    elif reorder_needed:
        urgency = "medium"
    else:
        urgency = "low"

    # Calculate confidence based on data availability
    confidence = 0.75  # Base confidence
    if market_trends and market_trends.get("trends"):
        confidence += 0.10  # Boost for market data
    if ctx.current_stock is not None and ctx.safety_stock is not None:
        confidence += 0.05  # Boost for complete stock data
    if ctx.lead_time_days is not None:
        confidence += 0.05  # Boost for lead time data
    confidence = min(confidence, 0.95)  # Cap at 95%

    data = {
        "product_id": ctx.product_id,
        "current_stock": current_stock,
        "safety_stock": safety_stock,
        "lead_time_days": lead_time_days,
        "adjusted_daily_demand": round(daily_demand, 1),
        "stock_covers_days": round(stock_covers_days, 1),
        "reorder_recommended": reorder_needed,
        "reorder_quantity": int(reorder_qty),
        "urgency": urgency,
        "market_adjusted": bool(trend_note),
        "confidence": round(confidence, 2),
        "model_used": "rule_based_reorder_with_trends",
    }

    if reorder_needed:
        explanation = (
            f"Current stock ({current_stock} units) covers only {stock_covers_days:.0f} days. "
            f"With {lead_time_days}-day lead time and {daily_demand:.0f} units/day demand, "
            f"recommend ordering {int(reorder_qty)} units. Urgency: {urgency}. "
            f"Confidence: {confidence*100:.0f}%. "
        )
        if trend_note:
            explanation += trend_note
    else:
        explanation = (
            f"Stock levels healthy at {current_stock} units, covering {stock_covers_days:.0f} days. "
            f"No reorder needed. Confidence: {confidence*100:.0f}%."
        )

    return ToolResult(
        tool_name="smart_reorder",
        success=True,
        data=data,
        explanation=explanation
    )


def pricelist_optimize_tool(ctx: ProductContext, market_trends: dict = None) -> ToolResult:
    """
    Pricelist Optimization Engine
    Suggests price adjustments or bundles based on inventory aging and market trends.
    """
    days_in_inventory = ctx.days_in_inventory or 0
    current_price = ctx.current_price or 0.0
    cost = ctx.cost or 0.0

    AGING_WARNING = 90
    AGING_CRITICAL = 120
    AGING_SEVERE = 180

    is_aging = days_in_inventory > AGING_CRITICAL
    # Determine if item is slow-moving based on aging (longer in inventory = slower)
    is_slow = days_in_inventory > AGING_WARNING

    # Determine markdown based on aging and market
    if days_in_inventory > AGING_SEVERE:
        markdown_pct = 25
        suggested_action = "aggressive_markdown"
    elif days_in_inventory > AGING_CRITICAL:
        markdown_pct = 15
        suggested_action = "markdown"
    elif is_slow and days_in_inventory > AGING_WARNING:
        markdown_pct = 10
        suggested_action = "light_markdown"
    elif is_slow:
        markdown_pct = 0
        suggested_action = "bundle"
    else:
        markdown_pct = 0
        suggested_action = "no_change"

    # Adjust based on market trends
    market_note = ""
    if market_trends:
        if market_trends.get("trend_direction") == "upward" and markdown_pct > 0:
            # Rising demand - less aggressive markdown
            markdown_pct = max(0, markdown_pct - 5)
            market_note = "Reduced markdown due to rising market demand."
        elif market_trends.get("trend_direction") == "downward" and suggested_action != "no_change":
            # Declining demand - more aggressive
            markdown_pct = min(35, markdown_pct + 5)
            market_note = "Increased markdown due to declining market trends."

    new_price = round(current_price * (1 - markdown_pct / 100), 2) if markdown_pct else current_price

    bundle_partner = None
    if suggested_action in ["bundle", "light_markdown"] and is_slow:
        bundle_partner = "SKU-HIGH-VELOCITY-001"

    margin_current = ((current_price - cost) / current_price * 100) if current_price > 0 else 0
    margin_new = ((new_price - cost) / new_price * 100) if new_price > 0 else 0

    # Calculate confidence based on data availability and market factors
    confidence = 0.70  # Base confidence
    if market_trends and market_trends.get("trends"):
        confidence += 0.12  # Boost for market data
    if ctx.days_in_inventory is not None:
        confidence += 0.05  # Boost for aging data
    if ctx.current_price is not None and ctx.cost is not None:
        confidence += 0.08  # Boost for pricing data
    confidence = min(confidence, 0.95)  # Cap at 95%

    data = {
        "product_id": ctx.product_id,
        "product_description": ctx.product_description,
        "days_in_inventory": days_in_inventory,
        "current_price": current_price,
        "cost": cost,
        "suggested_action": suggested_action,
        "markdown_percent": markdown_pct,
        "suggested_price": new_price,
        "bundle_partner_sku": bundle_partner,
        "current_margin_pct": round(margin_current, 1),
        "projected_margin_pct": round(margin_new, 1),
        "market_adjusted": bool(market_note),
        "market_trends_summary": market_trends.get("trends", [])[:3] if market_trends else [],
        "confidence": round(confidence, 2),
        "model_used": "rule_based_pricing_with_trends",
    }

    if suggested_action == "no_change":
        explanation = (
            f"Item moving well with only {days_in_inventory} days in inventory. "
            f"Margin {margin_current:.1f}% healthy. No change needed. "
            f"Confidence: {confidence*100:.0f}%."
        )
    elif suggested_action == "bundle":
        explanation = (
            f"Item has been in inventory for {days_in_inventory} days. "
            f"Recommend bundling with {bundle_partner} to increase turnover. "
            f"Confidence: {confidence*100:.0f}%."
        )
    else:
        explanation = (
            f"Aged {days_in_inventory} days in inventory. "
            f"Recommend {markdown_pct}% markdown (${current_price:.2f} → ${new_price:.2f}). "
            f"Margin: {margin_current:.1f}% → {margin_new:.1f}%. "
            f"Confidence: {confidence*100:.0f}%. "
        )
        if bundle_partner:
            explanation += f"Or bundle with {bundle_partner}. "
        if market_note:
            explanation += market_note

    return ToolResult(
        tool_name="pricelist_optimize",
        success=True,
        data=data,
        explanation=explanation
    )


# =============================================================================
# LangChain ReAct Agent
# =============================================================================

class ReactAgent:
    """
    ReAct Agent using LangChain and Gemini for inventory management.
    Implements Reason → Act → Observe loop with explainability.
    """

    INTENT_QUERIES = {
        "demand_forecast": "What's the demand forecast?",
        "smart_reorder": "What's the smart reorder recommendation?",
        "pricelist_optimize": "What's the pricelist optimization?",
        "full_analysis": "What's the full analysis?",
    }

    def __init__(self, use_llm: bool = True, api_key: str = None):
        """
        Initialize the agent.

        Args:
            use_llm: Whether to use Gemini LLM (requires API key)
            api_key: Google API key (or set GOOGLE_API_KEY env var)
        """
        self.use_llm = use_llm
        self.api_key = api_key or os.getenv("GOOGLE_API_KEY")
        self.scraper = MarketTrendsScraper()
        self.llm = None

        if self.use_llm and self.api_key and LANGCHAIN_AVAILABLE:
            self._init_llm()

    def _init_llm(self):
        """Initialize Gemini LLM via LangChain."""
        if not LANGCHAIN_AVAILABLE or not ChatGoogleGenerativeAI:
            print("Warning: LangChain not available")
            self.llm = None
            return

        try:
            self.llm = ChatGoogleGenerativeAI(
                model=GEMINI_MODEL,
                google_api_key=self.api_key,
                temperature=0.3,
                convert_system_message_to_human=True,
            )
        except Exception as e:
            print(f"Warning: Could not initialize Gemini LLM: {e}")
            self.llm = None

    def _get_langchain_tools(self, ctx: ProductContext, market_trends: dict) -> list:
        """Create LangChain Tool objects for the agent."""
        if not LANGCHAIN_AVAILABLE or not Tool:
            return []

        return [
            Tool(
                name="demand_forecast",
                func=lambda x: json.dumps(demand_forecast_tool(ctx, market_trends).__dict__),
                description=(
                    "Predict future product demand using historical sales data and market trends. "
                    "Use when asked about forecasting, predictions, or future demand."
                ),
            ),
            Tool(
                name="smart_reorder",
                func=lambda x: json.dumps(smart_reorder_tool(ctx, market_trends).__dict__),
                description=(
                    "Calculate optimal reorder quantity based on stock levels and lead time. "
                    "Use when asked about reordering, replenishment, or stock levels."
                ),
            ),
            Tool(
                name="pricelist_optimize",
                func=lambda x: json.dumps(pricelist_optimize_tool(ctx, market_trends).__dict__),
                description=(
                    "Suggest price adjustments or bundles for aging inventory. "
                    "Use when asked about pricing, markdowns, discounts, or slow-moving items."
                ),
            ),
        ]

    def _classify_intent(self, query: str) -> list[str]:
        """Classify user query into tool intents using keyword matching."""
        query_clean = query.strip().lower()

        # Check for exact matches first
        for intent, expected_query in self.INTENT_QUERIES.items():
            if query_clean == expected_query.lower():
                if intent == "full_analysis":
                    return ["demand_forecast", "smart_reorder", "pricelist_optimize"]
                return [intent]

        # Keyword-based matching for flexibility
        if any(kw in query_clean for kw in ["full analysis", "all", "complete", "everything"]):
            return ["demand_forecast", "smart_reorder", "pricelist_optimize"]

        if any(kw in query_clean for kw in ["demand", "forecast", "predict", "future demand"]):
            return ["demand_forecast"]

        if any(kw in query_clean for kw in ["reorder", "stock", "replenish", "order"]):
            return ["smart_reorder"]

        if any(kw in query_clean for kw in ["price", "pricing", "markdown", "discount", "optimize"]):
            return ["pricelist_optimize"]

        # Default to demand forecast only (not all 3)
        return ["demand_forecast"]

    def _fetch_market_trends(self, ctx: ProductContext) -> dict:
        """Fetch market trends via web scraping."""
        # Use product description for better search if available
        product_name = ctx.product_name or ctx.product_id
        if ctx.product_description:
            # Extract key terms from description for better search
            product_name = f"{product_name} {ctx.product_description[:50]}"
        category = ctx.category
        country = ctx.shop_country

        try:
            trends = self.scraper.scrape_news_trends(product_name, category, country)
            return trends
        except Exception as e:
            return {"error": str(e), "trends": []}

    def _reason_with_llm(self, query: str, ctx: ProductContext, market_trends: dict) -> str:
        """Use Gemini to reason about the query."""
        if not self.llm:
            return "LLM not available, using rule-based reasoning."

        system_prompt = """You are an AI assistant for inventory management.
        Analyze the user's query and the product context to provide reasoning about
        which tools to use and why. Consider market trends in your analysis.

        Available tools:
        - demand_forecast: Predict future demand
        - smart_reorder: Calculate reorder quantities
        - pricelist_optimize: Suggest pricing changes

        Provide a brief thought process."""

        context_str = f"""
        Product: {ctx.product_id}
        Product Name: {ctx.product_name}
        Description: {ctx.product_description or 'N/A'}
        Category: {ctx.category}
        Shop Country: {ctx.shop_country or 'Global'}
        Current Stock: {ctx.current_stock}
        Days in Inventory: {ctx.days_in_inventory}
        Market Sentiment: {market_trends.get('sentiment', 'unknown')}
        Market Trend Direction: {market_trends.get('trend_direction', 'unknown')}
        Recent Headlines: {[t.get('headline', '')[:50] for t in market_trends.get('trends', [])[:3]]}
        """

        try:
            messages = [
                SystemMessage(content=system_prompt),
                HumanMessage(content=f"Query: {query}\n\nContext:{context_str}"),
            ]
            response = self.llm.invoke(messages)
            return response.content
        except Exception as e:
            return f"LLM reasoning failed: {e}"

    def _synthesize_with_llm(self, query: str, results: list[ToolResult], market_trends: dict) -> str:
        """Use Gemini to synthesize a final answer."""
        if not self.llm:
            return self._synthesize_rule_based(results)

        results_str = "\n\n".join([
            f"**{r.tool_name}**:\n{r.explanation}"
            for r in results if r.success
        ])

        trends_str = ""
        if market_trends and market_trends.get("trends"):
            trends_str = "\n".join([
                f"- {t.get('headline', '')}"
                for t in market_trends.get("trends", [])[:3]
            ])

        prompt = f"""Based on the following analysis results and market trends,
        provide a concise, actionable summary for the user.

        Query: {query}

        Analysis Results:
        {results_str}

        Market Trends:
        Sentiment: {market_trends.get('sentiment', 'N/A')}
        Direction: {market_trends.get('trend_direction', 'N/A')}
        Headlines:
        {trends_str or 'No recent trends found.'}

        Provide a clear, professional response with specific recommendations.
        Each recommendation must include a reason (AI Explainability)."""

        try:
            messages = [HumanMessage(content=prompt)]
            response = self.llm.invoke(messages)
            return response.content
        except Exception as e:
            return self._synthesize_rule_based(results)

    def _synthesize_rule_based(self, results: list[ToolResult]) -> str:
        """Fallback synthesis without LLM."""
        if not results:
            return "Unable to process the request."

        parts = []
        for result in results:
            if result.success:
                parts.append(f"**{result.tool_name.replace('_', ' ').title()}**")
                parts.append(result.explanation)
                parts.append("")

        return "\n".join(parts).strip()

    def run(self, query: str, context: ProductContext) -> AgentResponse:
        """
        Execute the ReAct loop with LangChain and Gemini.
        """
        response = AgentResponse(query=query, intent="")

        try:
            # Step 1: Fetch market trends via web scraping
            market_trends = self._fetch_market_trends(context)
            response.market_trends = market_trends

            # Step 2: Reason (with LLM if available)
            intents = self._classify_intent(query)
            response.intent = intents[0] if len(intents) == 1 else "multi"

            if self.use_llm and self.llm:
                thought = self._reason_with_llm(query, context, market_trends)
            else:
                thought = f"Processing query for tools: {', '.join(intents)}"

            step = AgentStep(
                thought=thought,
                action="execute_tools",
                action_input={"tools": intents}
            )
            response.steps.append(step)

            # Step 3: Execute tools
            results = []
            for intent in intents:
                if intent == "demand_forecast":
                    results.append(demand_forecast_tool(context, market_trends))
                elif intent == "smart_reorder":
                    results.append(smart_reorder_tool(context, market_trends))
                elif intent == "pricelist_optimize":
                    results.append(pricelist_optimize_tool(context, market_trends))

            response.results = results

            # Step 4: Observe
            observation = "\n".join([
                f"[{r.tool_name}] {r.explanation}" for r in results if r.success
            ])
            step.observation = observation

            # Step 5: Synthesize final answer (with LLM if available)
            if self.use_llm and self.llm:
                response.final_answer = self._synthesize_with_llm(query, results, market_trends)
            else:
                response.final_answer = self._synthesize_rule_based(results)

        except Exception as e:
            response.errors.append(str(e))
            response.final_answer = f"Error: {str(e)}"

        return response

    def run_single_tool(self, tool_name: str, context: ProductContext) -> ToolResult:
        """Run a single tool directly."""
        market_trends = self._fetch_market_trends(context)

        if tool_name == "demand_forecast":
            return demand_forecast_tool(context, market_trends)
        elif tool_name == "smart_reorder":
            return smart_reorder_tool(context, market_trends)
        elif tool_name == "pricelist_optimize":
            return pricelist_optimize_tool(context, market_trends)
        else:
            return ToolResult(
                tool_name=tool_name,
                success=False,
                data={},
                explanation="",
                error=f"Unknown tool: {tool_name}"
            )


# =============================================================================
# Convenience Functions
# =============================================================================

def create_agent(use_llm: bool = True, api_key: str = None) -> ReactAgent:
    """Factory function to create an agent instance."""
    return ReactAgent(use_llm=use_llm, api_key=api_key)


def process_request(query: str, product_data: dict) -> dict:
    """
    Process a request from Go backend.
    """
    agent = create_agent()

    context = ProductContext(
        product_id=product_data.get("product_id", ""),
        product_name=product_data.get("product_name"),
        product_description=product_data.get("product_description"),
        category=product_data.get("category"),
        shop_country=product_data.get("shop_country"),
        sales_history=product_data.get("sales_history"),
        history_months=product_data.get("history_months", 6),
        current_stock=product_data.get("current_stock"),
        safety_stock=product_data.get("safety_stock"),
        lead_time_days=product_data.get("lead_time_days"),
        days_in_inventory=product_data.get("days_in_inventory"),
        current_price=product_data.get("current_price"),
        cost=product_data.get("cost"),
    )

    response = agent.run(query, context)

    return {
        "query": response.query,
        "intent": response.intent,
        "results": [
            {
                "tool": r.tool_name,
                "success": r.success,
                "data": r.data,
                "explanation": r.explanation,
                "error": r.error,
            }
            for r in response.results
        ],
        "market_trends": response.market_trends,
        "final_answer": response.final_answer,
        "errors": response.errors,
    }


# =============================================================================
# CLI for Testing
# =============================================================================

if __name__ == "__main__":
    print("=" * 60)
    print("ReAct Agent with LangChain, Gemini, and Web Scraping")
    print("=" * 60)

    # Check for API key
    api_key = os.getenv("GOOGLE_API_KEY")
    if api_key:
        print(f"✓ Gemini API key found")
    else:
        print("⚠ No GOOGLE_API_KEY found - running without LLM")

    agent = create_agent(use_llm=bool(api_key))

    test_context = ProductContext(
        product_id="SKU-123",
        product_name="Wireless Bluetooth Headphones",
        product_description="Premium noise-cancelling wireless headphones with 40-hour battery life, Bluetooth 5.0, and foldable design",
        category="Electronics",
        shop_country="BD",  # Bangladesh
        sales_history=[
            {"date": "2026-01-01", "qty": 15},
            {"date": "2026-01-02", "qty": 18},
            {"date": "2026-01-03", "qty": 12},
        ],
        history_months=3,
        current_stock=100,
        safety_stock=50,
        lead_time_days=14,
        days_in_inventory=145,
        current_price=29.99,
        cost=15.00,
    )

    print("\n" + "=" * 60)
    print("Testing: What's the demand forecast?")
    print("=" * 60)
    response = agent.run("What's the demand forecast?", test_context)
    print(f"\nIntent: {response.intent}")
    print(f"\nMarket Trends: {response.market_trends.get('sentiment', 'N/A')} ({response.market_trends.get('trend_direction', 'N/A')})")
    if response.market_trends.get("trends"):
        print("Headlines:")
        for t in response.market_trends["trends"][:3]:
            print(f"  - {t.get('headline', '')[:60]}...")
    print(f"\nFinal Answer:\n{response.final_answer}")

    print("\n" + "=" * 60)
    print("Testing: What's the full analysis?")
    print("=" * 60)
    response = agent.run("What's the full analysis?", test_context)
    print(f"\nIntent: {response.intent}")
    print(f"\nFinal Answer:\n{response.final_answer}")
