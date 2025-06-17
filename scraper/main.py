import os
import pathlib
import rakuten_security_scraper_pb2_grpc
import rakuten_security_scraper_pb2
import grpc
import sys
import threading
import uvicorn
from concurrent import futures
from logging import getLogger, StreamHandler
import cache
from distutils.util import strtobool
from fastapi_app import app, set_servicer


class RakutenSecurityScraperServicer(rakuten_security_scraper_pb2_grpc.RakutenSecurityScraperServicer):
    download_dir = ""
    id = ""
    password = ""
    cache = {}

    def __init__(self, id, password, download_dir, cache_ttl_second, cache_clear_on_start):
        self.id = id
        self.password = password

        self.download_dir = download_dir
        if cache_clear_on_start:
            for f in [f for f in pathlib.Path(self.download_dir).glob("*.csv")]:
                f.unlink()

        self.cache = cache.CachedRakutenSecurityScraper(
            id, password, download_dir, cache_ttl_second)

    def ListWithdrawalHistories(self, request, context):
        try:
            history = self.cache.GetWithdrawalHistories()
            response = rakuten_security_scraper_pb2.ListWithdrawalHistoriesResponse()
            response.history.extend(
                [self._covertToWithdrawalHisoty(e) for e in history])
            return response
        except Exception as e:
            logger.error(f"Failed to get withdrawal histories: {str(e)}")
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details("Data not available. Please try again later.")
            return rakuten_security_scraper_pb2.ListWithdrawalHistoriesResponse()

    def ListDividendHistories(self, request, context):
        try:
            history = self.cache.GetDividendHistories()
            response = rakuten_security_scraper_pb2.ListDividendHistoriesResponse()
            response.history.extend(
                [self._convertToDividendHistory(e) for e in history])
            return response
        except Exception as e:
            logger.error(f"Failed to get dividend histories: {str(e)}")
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details("Data not available. Please try again later.")
            return rakuten_security_scraper_pb2.ListDividendHistoriesResponse()

    def TotalAssets(self, request, context):
        try:
            asset, currency_rate = self.cache.GetTotalAsset()
            response = rakuten_security_scraper_pb2.TotalAssetResponse()
            response.asset.extend([self._convertToAsset(e) for e in asset])
            response.currency_rate.extend(
                [self._convertToCurrencyRate(e) for e in currency_rate])
            return response
        except Exception as e:
            logger.error(f"Failed to get total assets: {str(e)}")
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details("Data not available. Please try again later.")
            return rakuten_security_scraper_pb2.TotalAssetResponse()

    def _covertToWithdrawalHisoty(self, elem):
        return rakuten_security_scraper_pb2.WithdrawalHistory(
            date=elem["date"],
            amount=elem["amount"],
            currency=elem["currency"]
        )

    def _convertToDividendHistory(self, elem):
        return rakuten_security_scraper_pb2.DividendHistory(
            date=elem["date"],
            account=elem["account"],
            type=elem["type"],
            ticker=elem["ticker"],
            name=elem["name"],
            currency=elem["currency"],
            count=elem["count"],
            dividend_unitprice=elem["dividend_unitprice"],
            dividend_total_before_taxes=elem["dividend_total_before_taxes"],
            total_taxes=elem["total_taxes"],
            dividend_total=elem["dividend_total"]
        )

    def _convertToAsset(self, elem):
        return rakuten_security_scraper_pb2.Asset(
            type=elem["type"],
            ticker=elem["ticker"],
            name=elem["name"],
            currency=elem["currency"],
            account=elem["account"],
            count=elem["count"],
            average_acquisition_price=elem["average_acquisition_price"],
            current_unit_price=elem["current_unit_price"],
            current_price=elem["current_price"],
            current_price_yen=elem["current_price_yen"]
        )

    def _convertToCurrencyRate(self, elem):
        return rakuten_security_scraper_pb2.CurrenyRateToJPY(
            currencyCode=elem["currency_code"],
            rate=elem["rate"]
        )


id = os.environ['RAKUTEN_SEC_ID']
password = os.environ['RAKUTEN_SEC_PASSWORD']
download_dir = os.environ['DOWNLOAD_DIR']
cache_ttl_second = int(os.getenv('CACHE_TTL_SECOND', '86400'))
cache_clear_on_start = bool(
    strtobool(os.getenv('CACHE_CLEAR_ON_START', 'True')))


logger = getLogger("rakuten-security-scraper")
logger.setLevel('INFO')
stream_handler = StreamHandler(stream=sys.stdout)
logger.addHandler(stream_handler)


def run_grpc_server():
    """Run gRPC server"""
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    servicer = RakutenSecurityScraperServicer(
        id, password, download_dir, cache_ttl_second, cache_clear_on_start)
    rakuten_security_scraper_pb2_grpc.add_RakutenSecurityScraperServicer_to_server(
        servicer, server)

    # Set the servicer for FastAPI
    set_servicer(servicer)

    server.add_insecure_port("0.0.0.0:50051")
    logger.info("Starting gRPC server on port 50051")
    server.start()
    server.wait_for_termination()


def run_fastapi_server():
    """Run FastAPI server"""
    logger.info("Starting FastAPI server on port 8000")
    uvicorn.run(app, host="0.0.0.0", port=8000, log_level="info")


if __name__ == "__main__":
    # Start gRPC server in a separate thread
    grpc_thread = threading.Thread(target=run_grpc_server)
    grpc_thread.daemon = True
    grpc_thread.start()

    # Run FastAPI server in the main thread
    run_fastapi_server()
