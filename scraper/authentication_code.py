from googleapiclient.discovery import build
from datetime import datetime
from google.auth.transport.requests import Request
import pickle
import base64
import re
import dateutil.parser
import logging
from authorize import get_auth_file_path

logger = logging.getLogger("rakuten-security-scraper")


def get_gmail_service():
    """Gmail APIのサービスを初期化して返す"""
    try:
        # 保存済みのトークンを読み込む
        token_path = get_auth_file_path('token.pickle')
        with open(token_path, 'rb') as token:
            creds = pickle.load(token)

        # 必要ならリフレッシュ
        if creds.expired and creds.refresh_token:
            creds.refresh(Request())
            with open(token_path, 'wb') as token_file:
                pickle.dump(creds, token_file)

        # Gmail API クライアント作成
        return build('gmail', 'v1', credentials=creds)
    except Exception:
        return None


def fetch_newest_message(service, query):
    """指定したクエリに一致する最も新しいメールを取得する"""
    try:
        # メールを検索
        results = service.users().messages().list(
            userId='me', q=query).execute()
        messages = results.get('messages', [])

        if not messages:
            return None

        # メール一覧を日付でソートするために日付情報を取得
        message_dates = []
        for message in messages:
            # メッセージの詳細情報を取得
            msg = service.users().messages().get(
                userId='me', id=message['id']).execute()

            # 日付を取得して解析
            headers = msg['payload']['headers']
            date_str = next(
                (h['value'] for h in headers if h['name'] == 'Date'),
                'No date')
            try:
                parsed_date = dateutil.parser.parse(date_str)
                message_dates.append((parsed_date, message['id'], msg))
            except Exception:
                # 日付の解析に失敗した場合は無視
                continue

        # 日付でソートし、最も新しいメールを取得
        if not message_dates:
            return None

        message_dates.sort(reverse=True)  # 降順ソート（最新が先頭）
        _, _, newest_msg = message_dates[0]
        return newest_msg
    except Exception:
        return None


def extract_authentication_codes(message):
    """メッセージから認証コードを抽出する"""
    extracted_values = []

    try:
        if 'parts' not in message['payload']:
            return []

        for part in message['payload']['parts']:
            if part['mimeType'] == 'text/html':
                body = part.get('body', {}).get('data', '')
                if not body:
                    continue

                html = base64.urlsafe_b64decode(body).decode('utf-8')
                html_lines = html.split('\n')
                pattern = re.compile(r'絵文字.の内容')
                extracted_lines = [
                    line for line in html_lines
                    if pattern.search(line)
                ]

                for line in extracted_lines:
                    # 行を空白で分割し、2つ目の値（インデックス1）を取得
                    parts = line.strip().split()
                    if len(parts) > 1:
                        extracted_values.append(parts[1])
    except Exception:
        pass

    return extracted_values


def get_authentication_codes(timestamp):
    """
    楽天証券のログイン追加認証コードを含むメールから認証コードを抽出する関数

    Args:
        timestamp (str): 検索するタイムスタンプ（YYYY/MM/DD HH:MM:SS形式）。
                        これより後のメールを検索。タイムスタンプには必ず日付と時刻の両方を含める必要があります。

    Returns:
        list: 抽出した認証コードの配列。エラー発生時や該当するデータがない場合は空の配列を返す
    """
    try:
        # Gmail APIサービスの初期化
        service = get_gmail_service()
        if not service:
            return []

        # タイムスタンプのフォーマット確認
        if " " not in timestamp:
            # 時刻部分がない場合はエラーとして空のリストを返す
            return []

        search_date, search_time = timestamp.split(" ", 1)

        search_datetime = datetime.strptime(timestamp, "%Y/%m/%d %H:%M:%S")
        unix_timestamp = int(search_datetime.timestamp())

        # 検索クエリの指定（after条件にはPST日付を使用）
        query = (f"from:service@rakuten-sec.co.jp after:{unix_timestamp} "
                 "subject:楽天証券よりログイン追加認証コードを送付いたします")

        # 最新のメッセージを取得
        newest_msg = fetch_newest_message(service, query)
        if not newest_msg:
            return []

        # 時刻指定がある場合は、取得したメッセージのタイムスタンプを確認
        if search_time and newest_msg:
            # メッセージのヘッダーから日時を取得
            headers = newest_msg['payload']['headers']
            date_str = next(
                (h['value'] for h in headers if h['name'] == 'Date'),
                None
            )

            if date_str:
                try:
                    # (JST)などのタイムゾーン表記を削除してからパース
                    cleaned_date_str = re.sub(
                        r'\s*\([A-Z]+\)\s*$', '', date_str)
                    msg_date = dateutil.parser.parse(cleaned_date_str)
                    search_full_dt = dateutil.parser.parse(
                        f"{search_date} {search_time}")

                    # 両方のdatetimeオブジェクトをタイムゾーン情報なし(naive)に統一
                    if msg_date.tzinfo is not None:
                        msg_date = msg_date.replace(tzinfo=None)
                    if search_full_dt.tzinfo is not None:
                        search_full_dt = search_full_dt.replace(tzinfo=None)

                    # 指定したタイムスタンプより前なら除外
                    if msg_date <= search_full_dt:
                        return []
                except Exception as e:
                    logger.error(f"日付解析エラー: {date_str}, エラーの詳細: {str(e)}")
                    # 日付解析エラーの場合は処理続行
                    pass

        # 認証コードを抽出
        return extract_authentication_codes(newest_msg)

    except Exception:
        # 何らかの例外が発生した場合は空の配列を返す
        return []
