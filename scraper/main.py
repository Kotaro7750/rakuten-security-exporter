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
            id, password, download_dir, 30)

    def ListWithdrawalHistorys(self, request, context):

        history = self.cache.GetWithdrawalHistories()

        response = rakuten_security_scraper_pb2.Response()
        response.history.extend(
            [self._covertToWithdrawalHisoty(e) for e in history])
        return response

    def _covertToWithdrawalHisoty(self, elem):
        return rakuten_security_scraper_pb2.WithdrawalHistory(
            date=elem["date"],
            amount=elem["amount"],
            type=elem["type"],
            currency=elem["currency"]
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

server.add_insecure_port("localhost:50051")
server.start()
server.wait_for_termination()
