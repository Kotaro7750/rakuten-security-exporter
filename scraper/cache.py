import scraper
import pathlib
import datetime
import parse_csv


class CachedRakutenSecurityScraper():
    download_dir = ""
    id = ""
    password = ""
    ttl_second = 86400  # 1 day

    def __init__(self, id, password, download_dir, ttl_second):
        self.id = id
        self.password = password
        self.download_dir = download_dir
        self.ttl_second = ttl_second

    def _is_file_cached(self, path):
        print(path)
        if not pathlib.Path(path).exists():
            print("not exist")
            return False

        timestamp = pathlib.Path(path).stat().st_mtime

        updated_time = datetime.datetime.fromtimestamp(timestamp)
        current_time = datetime.datetime.now()

        print((current_time - updated_time).total_seconds())

        return (current_time - updated_time).total_seconds() <= self.ttl_second

    def GetWithdrawalHistories(self):
        path = pathlib.Path(self.download_dir, "withdrawal.csv")

        if not self._is_file_cached(path):
            scraper.download_withdrawal_history(
                self.id, self.password, self.download_dir)
            list(pathlib.Path(self.download_dir).glob(
                "Withdrawal*.csv"))[0].rename(path)

        return parse_csv.parse_withdrawal_history(self.download_dir)
