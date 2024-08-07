import os
import pathlib
import rakuten_security_scraper_pb2_grpc
import rakuten_security_scraper_pb2
import grpc
import sys
from concurrent import futures
from logging import getLogger, StreamHandler
import cache


class RakutenSecurityScraperServicer(rakuten_security_scraper_pb2_grpc.RakutenSecurityScraperServicer):
    download_dir = ""
    id = ""
    password = ""
    cache = {}

    def __init__(self, id, password, download_dir):
        self.id = id
        self.password = password

        self.download_dir = download_dir
        for f in [f for f in pathlib.Path(self.download_dir).glob("*.csv")]:
            f.unlink()

        self.cache = cache.CachedRakutenSecurityScraper(
            id, password, download_dir, 86400)

    def ListWithdrawalHistories(self, request, context):

        history = self.cache.GetWithdrawalHistories()

        response = rakuten_security_scraper_pb2.ListWithdrawalHistoriesResponse()
        response.history.extend(
            [self._covertToWithdrawalHisoty(e) for e in history])
        return response

    def ListDividendHistories(self, request, context):
        history = self.cache.GetDividendHistories()

        response = rakuten_security_scraper_pb2.ListDividendHistoriesResponse()
        response.history.extend([self._convertToDividendHistory(e) for e in history])

        return response

    def TotalAssets(self, request, context):
        asset, currency_rate = self.cache.GetTotalAsset()

        response = rakuten_security_scraper_pb2.TotalAssetResponse()
        response.asset.extend([self._convertToAsset(e) for e in asset])
        response.currency_rate.extend([self._convertToCurrencyRate(e) for e in currency_rate])

        return response

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

logger = getLogger("rakuten-security-scraper")
logger.setLevel('INFO')
stream_handler = StreamHandler(stream=sys.stdout)
logger.addHandler(stream_handler)

server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
rakuten_security_scraper_pb2_grpc.add_RakutenSecurityScraperServicer_to_server(
    RakutenSecurityScraperServicer(id, password, download_dir), server)

server.add_insecure_port("0.0.0.0:50051")
server.start()
server.wait_for_termination()
