import scraper
import pathlib
import datetime
import parse_csv
import threading
from concurrent.futures import ThreadPoolExecutor
from logging import getLogger
from authorize import check_token_validity

logger = getLogger("rakuten-security-scraper")


class CachedRakutenSecurityScraper():
    download_dir = ""
    id = ""
    password = ""
    ttl_second = 86400  # 1 day
    _scraping_lock = None
    _executor = None

    # When ttl_second is minus, cache never expires
    def __init__(self, id, password, download_dir, ttl_second):
        self.id = id
        self.password = password
        self.download_dir = download_dir
        self.ttl_second = ttl_second
        self._scraping_lock = threading.Lock()
        self._executor = ThreadPoolExecutor(max_workers=1)

    def _is_file_cached(self, path):
        logger.info("Check if %s exists", path)
        if not pathlib.Path(path).exists():
            logger.info("%s does not exist", path)
            return False

        if self.ttl_second < 0:
            return True

        timestamp = pathlib.Path(path).stat().st_mtime

        updated_time = datetime.datetime.fromtimestamp(timestamp)
        current_time = datetime.datetime.now()

        print((current_time - updated_time).total_seconds())

        return (current_time - updated_time).total_seconds() <= self.ttl_second

    def GetWithdrawalHistories(self):
        path = pathlib.Path(self.download_dir, "withdrawal.csv")

        if not self._is_file_cached(path):
            # キャッシュにない場合は失敗させ、非同期でスクレイピングを開始
            self._start_async_scraping_if_needed()
            raise Exception(
                "Data not cached. Scraping initiated in background.")

        return parse_csv.parse_withdrawal_history(self.download_dir)

    def GetDividendHistories(self):
        path = pathlib.Path(self.download_dir, "dividend.csv")

        if not self._is_file_cached(path):
            # キャッシュにない場合は失敗させ、非同期でスクレイピングを開始
            self._start_async_scraping_if_needed()
            raise Exception(
                "Data not cached. Scraping initiated in background.")

        return parse_csv.parse_dividend_history(self.download_dir)

    def GetTotalAsset(self):
        path = pathlib.Path(self.download_dir, "asset.csv")

        if not self._is_file_cached(path):
            # キャッシュにない場合は失敗させ、非同期でスクレイピングを開始
            self._start_async_scraping_if_needed()
            raise Exception(
                "Data not cached. Scraping initiated in background.")

        return parse_csv.parse_asset(self.download_dir)

    def _start_async_scraping_if_needed(self):
        """ロックを使用して1つのスクレイピングのみ実行する"""
        # ロックを獲得できない場合（すでにスクレイピング中）は何もしない
        if not self._scraping_lock.acquire(blocking=False):
            logger.info("Scraping already in progress, skipping")
            return

        try:
            # 認証トークンをチェック
            token_check = check_token_validity()
            if not token_check.get("is_valid", False):
                logger.error(
                    "Authentication token is invalid or expired. Cannot start scraping.")
                self._scraping_lock.release()
                return

            logger.info("Starting async scraping")
            # 非同期でスクレイピングを実行
            self._executor.submit(self._perform_scraping)
        except Exception as e:
            logger.error(f"Failed to start async scraping: {str(e)}")
        finally:
            # この時点ではロックを保持したまま
            pass

    def _perform_scraping(self):
        """実際のスクレイピング処理（別スレッドで実行）"""
        try:
            logger.info("Performing scraping")
            result = scraper.scrape(self.id, self.password, self.download_dir)

            if result["status"] == "success":
                logger.info("Scraping completed successfully")
                # CSVファイルのリネーム処理
                self._rename_downloaded_files()
            else:
                logger.error(f"Scraping failed: {result['message']}")

        except Exception as e:
            logger.error(f"Error during scraping: {str(e)}")
        finally:
            # スクレイピング完了後にロックを解放
            self._scraping_lock.release()
            logger.info("Scraping lock released")

    def _rename_downloaded_files(self):
        """ダウンロードされたファイルを適切な名前にリネーム"""
        try:
            # withdrawal.csv
            withdrawal_files = list(pathlib.Path(
                self.download_dir).glob("Withdrawal*.csv"))
            if withdrawal_files:
                withdrawal_files[0].rename(pathlib.Path(
                    self.download_dir, "withdrawal.csv"))

            # dividend.csv
            dividend_files = list(pathlib.Path(
                self.download_dir).glob("dividendlist_*.csv"))
            if dividend_files:
                dividend_files[0].rename(pathlib.Path(
                    self.download_dir, "dividend.csv"))

            # asset.csv
            asset_files = list(pathlib.Path(
                self.download_dir).glob("assetbalance*.csv"))
            if asset_files:
                asset_files[0].rename(pathlib.Path(
                    self.download_dir, "asset.csv"))

        except Exception as e:
            logger.error(f"Error renaming files: {str(e)}")
