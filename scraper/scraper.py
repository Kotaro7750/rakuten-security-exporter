import time
from datetime import datetime
from datetime import timedelta
from typing import Dict
from playwright.sync_api import sync_playwright
from authentication_code import get_authentication_codes
from authorize import check_token_validity
from logging import getLogger

logger = getLogger("rakuten-security-scraper")


def scrape(id: str, password: str, download_dir: str) -> Dict[str, str]:
    """
    楽天証券の認証を実行し、結果を返す関数

    Args:
        id: 楽天証券のログインID
        password: 楽天証券のパスワード
        download_dir: ダウンロード先ディレクトリ

    Returns:
        Dict[str, str]: スクレイピング結果を含む辞書
    """
    result = {
        "status": "success",
        "message": "認証が完了しました",
        "timestamp": datetime.now().strftime("%Y/%m/%d %H:%M:%S")
    }

    # 認証トークンのチェック（スクレイピング前の早期終了）
    logger.info("Checking authentication token validity before scraping")
    token_check = check_token_validity()
    if not token_check.get("is_valid", False):
        logger.error(
            "Authentication token is invalid or expired. Terminating scraping early.")
        result["status"] = "error"
        result["message"] = "認証トークンが無効または期限切れです。スクレイピングを中止しました。"
        return result

    logger.info("Authentication token is valid, proceeding with scraping")

    try:
        with sync_playwright() as playwright:
            browser = playwright.chromium.launch(
                headless=True, downloads_path=download_dir)
            context = browser.new_context()
            page = context.new_page()
            page.goto("https://www.rakuten-sec.co.jp/ITS/V_ACT_Login.html")

            page.get_by_role("textbox", name="ログインID").click()
            page.get_by_role("textbox", name="ログインID").fill(id)
            page.get_by_role("textbox", name="パスワード").click()
            page.get_by_role("textbox", name="パスワード").fill(password)

            # 30秒前のタイムスタンプを生成
            current_time = datetime.now()
            timestamp = (current_time - timedelta(seconds=30)
                         ).strftime("%Y/%m/%d %H:%M:%S")

            page.get_by_role("button", name=" ログインする").click()

            auth_codes = get_and_display_authentication_codes(
                timestamp=timestamp)
            print(auth_codes)
            page.get_by_role("button", name=auth_codes[0], exact=True).click()
            page.get_by_role("button", name=auth_codes[1], exact=True).click()
            page.get_by_role("button", name="認証する").click()

            page.get_by_role("link", name="チャットで問合せる").hover()
            page.get_by_role("button", name="マイメニュー 口座管理・入出金など").click()

            page.locator("#megaMenu").get_by_role(
                "link", name="保有商品一覧").click()
            page.get_by_role("cell", name="保有商品の評価額合計").hover()

            with page.expect_download() as download_info:
                page.get_by_role("link", name="CSVで保存").click()
            download = download_info.value
            download.save_as(f'{download_dir}/asset.csv')

            page.get_by_role("link", name="楽天証券").click()

            page.get_by_role("link", name="チャットで問合せる").hover()
            page.get_by_role("button", name="マイメニュー 口座管理・入出金など").click()

            page.get_by_role("link", name="入出金履歴", exact=True).click()
            with page.expect_download() as download_info:
                page.get_by_role("link", name="CSVで保存").click()
            download = download_info.value
            download.save_as(f'{download_dir}/withdrawal.csv')

            page.get_by_role("link", name="楽天証券").click()

            page.get_by_role("link", name="チャットで問合せる").hover()
            page.get_by_role("button", name="マイメニュー 口座管理・入出金など").click()
            page.locator("#megaMenu").get_by_role(
                "link", name="配当・分配金").click()
            page.get_by_role("img", name="すべて").click()
            page.get_by_role("button", name="Submit").click()
            with page.expect_download() as download_info:
                page.get_by_role("link", name="CSVで保存").click()
            download = download_info.value
            download.save_as(f'{download_dir}/dividend.csv')

        # Navigation successful
        result["message"] = "認証が完了しました"

    except Exception as e:
        result["status"] = "error"
        result["message"] = f"エラーが発生しました: {str(e)}"
    finally:
        try:
            if 'context' in locals():
                context.close()
            if 'browser' in locals():
                browser.close()
        except:
            pass

    print(result)
    return result


def get_and_display_authentication_codes(timestamp=None):
    """
    認証コードを取得して表示する関数（exponential backoffリトライ処理付き）

    Args:
        timestamp (str, optional): 認証コードを検索する開始タイムスタンプ（YYYY/MM/DD HH:MM:SS形式）
                                  指定がない場合は現在の日時を使用

    Returns:
        list: 取得した認証コードのリスト
    """
    import time

    # タイムスタンプが指定されていない場合は現在の日時を使用
    if timestamp is None:
        current_time = datetime.now()
        timestamp = current_time.strftime("%Y/%m/%d %H:%M:%S")

    print(f"検索開始時刻: {timestamp}")
    print("認証コードを取得しています...")

    # 認証コードを取得（exponential backoffリトライ処理を実装）
    max_retries = 5
    retry_count = 0
    auth_codes = []

    # Exponential backoff設定
    initial_wait = 5  # 初期待機時間（秒）
    multiplier = 2    # 待機時間の倍率
    max_wait = 60     # 最大待機時間（秒）

    while retry_count < max_retries and not auth_codes:
        if retry_count > 0:
            # exponential backoffで待機時間を計算
            wait_time = min(initial_wait * (multiplier **
                            (retry_count - 1)), max_wait)
            print(f"リトライ {retry_count}/{max_retries}... {wait_time}秒待機します")
            time.sleep(wait_time)

        auth_codes = get_authentication_codes(timestamp=timestamp)
        retry_count += 1

        # 結果の表示
        if auth_codes:
            print(f"{timestamp} 以降の認証コード:")
            for i, code in enumerate(auth_codes, 1):
                print(f"コード {i}: {code}")
        elif retry_count < max_retries:
            next_wait_time = min(
                initial_wait * (multiplier ** retry_count), max_wait)
            print(f"{timestamp} 以降の認証コードを取得できませんでした。{next_wait_time}秒後にリトライします。")
        else:
            print(f"{timestamp} 以降の認証コードを取得できませんでした。最大リトライ回数に達しました。")

    return auth_codes
