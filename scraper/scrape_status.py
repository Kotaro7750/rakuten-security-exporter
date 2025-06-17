"""Enums and data classes for Rakuten Security Scraper."""

from dataclasses import dataclass
from datetime import datetime
from enum import Enum
from typing import Optional


class ScrapeResultEnum(str, Enum):
    """Enum for scraping result status."""

    SUCCESS = "success"
    FAILURE = "failure"
    RUNNING = "running"
    IDLE = "idle"
    UNKNOWN = "unknown"

    def __str__(self) -> str:
        """Return the string value of the enum."""
        return self.value


@dataclass
class ScrapeResult:
    """Data class for scraping result."""

    status: ScrapeResultEnum
    message: str
    timestamp: str
    auth_url: Optional[str] = None

    @classmethod
    def create_success(
        cls, message: str = "処理が正常に完了しました"
    ) -> "ScrapeResult":
        """Create a success result."""
        return cls(
            status=ScrapeResultEnum.SUCCESS,
            message=message,
            timestamp=datetime.now().strftime("%Y/%m/%d %H:%M:%S")
        )

    @classmethod
    def create_failure(
        cls, message: str = "処理に失敗しました"
    ) -> "ScrapeResult":
        """Create a failure result."""
        return cls(
            status=ScrapeResultEnum.FAILURE,
            message=message,
            timestamp=datetime.now().strftime("%Y/%m/%d %H:%M:%S")
        )
