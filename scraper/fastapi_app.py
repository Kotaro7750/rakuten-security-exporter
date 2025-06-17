from datetime import datetime
from typing import Dict, Any, Optional
from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from playwright.sync_api import sync_playwright

from authorize import authorize, check_token_validity

import logging

app = FastAPI(title="Rakuten Security Scraper API", version="1.0.0")

logger = logging.getLogger("rakuten-security-scraper")

# Global variable to hold the servicer instance
_servicer = None


def set_servicer(servicer):
    """Set the servicer instance to be used by FastAPI endpoints"""
    global _servicer
    _servicer = servicer


# CORSミドルウェアの設定
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# グローバル変数で最新のスクレイピング結果を保持
latest_scraping_result: Optional[Dict[str, Any]] = None


# Response models for different endpoints
class BaseResponse(BaseModel):
    """Base response model with common fields."""
    status: str
    message: str
    timestamp: str


class StatusResponse(BaseResponse):
    """Response model for status endpoint."""
    scraping_time: Optional[str] = None
    scraping_result: Optional[str] = None  # success, failure, running, idle
    auth_token_status: Optional[str] = None  # valid, invalid
    auth_token_expiry: Optional[str] = None


class AuthResponse(BaseResponse):
    """Response model for auth endpoint."""
    auth_url: str  # Required for auth endpoint


class OAuthCallbackResponse(BaseResponse):
    """Response model for OAuth callback endpoint."""
    pass  # Uses only base fields


async def run_auth_job(auth_code: str):
    """Run the Google API authentication job asynchronously."""
    try:
        # 認証プロセスを実行（同期的な関数なのでawaitを使用しない）
        result = authorize(auth_code=auth_code)
        # 認証結果は別途管理されるため、グローバル変数に保存しない

    except Exception as e:
        # 認証中の例外を処理 (ログ出力のみ)
        print(f"認証実行中にエラーが発生しました: {str(e)}")


@app.get("/status", response_model=StatusResponse)
async def get_status():
    """Get the latest scraping and auth status."""
    # Get auth token status
    token_check = check_token_validity()

    # Determine scraping status and time
    scraping_result = "idle"
    scraping_time = None
    base_message = "システム状態を確認しました"

    if latest_scraping_result:
        scraping_result = latest_scraping_result.get("status", "unknown")
        scraping_time = latest_scraping_result.get("timestamp")
        base_message = latest_scraping_result.get("message", base_message)

    # Determine auth token status
    auth_status = "valid" if token_check.get("is_valid", False) else "invalid"
    auth_expiry = token_check.get("token_expiry", None)

    return StatusResponse(
        status="success",
        message=base_message,
        timestamp=datetime.now().strftime("%Y/%m/%d %H:%M:%S"),
        scraping_time=scraping_time,
        scraping_result=scraping_result,
        auth_token_status=auth_status,
        auth_token_expiry=auth_expiry
    )


@app.get("/auth", response_model=AuthResponse)
async def start_auth():
    """Start the OAuth flow and return the authorization URL."""
    try:
        # authorize関数は同期的なので、awaitを使用しない
        result = authorize()
        # auth_urlが存在する場合は、それをレスポンスに含める
        if "auth_url" in result:
            return AuthResponse(
                status=result["status"],
                message=result["message"],
                timestamp=result["timestamp"],
                auth_url=result["auth_url"]
            )
        raise HTTPException(status_code=500, detail="認証URLの取得に失敗しました")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/oauth2callback", response_model=OAuthCallbackResponse)
async def oauth_callback(code: str, background_tasks: BackgroundTasks):
    """Handle the OAuth callback with the authorization code."""
    if not code:
        raise HTTPException(status_code=400, detail="認証コードがありません")

    # run_auth_jobは非同期関数なので、awaitを使用
    background_tasks.add_task(run_auth_job, code)

    return OAuthCallbackResponse(
        status="accepted",
        message="認証コードを処理中...",
        timestamp=datetime.now().strftime("%Y/%m/%d %H:%M:%S")
    )
