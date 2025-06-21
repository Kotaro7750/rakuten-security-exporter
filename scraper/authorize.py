from google_auth_oauthlib.flow import Flow
from google.auth.transport.requests import Request
from typing import Dict, Any, Optional
import pickle
import os
import os.path
import json
import logging
from datetime import datetime

# Gmail API のスコープ（読み取り専用）
SCOPES = ['https://www.googleapis.com/auth/gmail.readonly']
# OAuth callback URI
REDIRECT_URI = "http://localhost:8000/oauth2callback"

# ログ設定
logger = logging.getLogger("rakuten-security-scraper")


def get_auth_file_path(filename: str) -> str:
    """
    環境変数GMAIL_API_AUTH_DIRを使用してファイルパスを構築する関数

    Args:
        filename (str): ファイル名

    Returns:
        str: 完全なファイルパス

    Raises:
        FileNotFoundError: 指定されたディレクトリが存在しない場合
    """
    auth_dir = os.environ.get('GMAIL_API_AUTH_DIR', '.')
    # ディレクトリが存在しない場合はエラーを発生
    if not os.path.exists(auth_dir):
        error_msg = f"認証ファイル用ディレクトリが存在しません: {auth_dir}"
        logger.error(error_msg)
        raise FileNotFoundError(error_msg)
    return os.path.join(auth_dir, filename)


def get_credentials():
    """
    保存されたトークンからcredentialsを取得する関数

    Returns:
        Credentials or None: 保存されたトークン情報、ない場合はNone
    """
    creds = None
    token_path = get_auth_file_path('token.pickle')

    # トークンが保存されていればそれを使う
    if os.path.exists(token_path):
        with open(token_path, 'rb') as token:
            creds = pickle.load(token)

    # リフレッシュが必要かチェック
    if creds and creds.expired and creds.refresh_token:
        creds.refresh(Request())
        # リフレッシュしたトークンを保存
        with open(token_path, 'wb') as token:
            pickle.dump(creds, token)

    return creds


def generate_auth_url() -> Dict[str, Any]:
    """
    クライアント側で使用する認証URLを生成する関数

    Returns:
        Dict[str, Any]: 認証URL情報を含む辞書
    """
    result = {
        "status": "success",
        "message": "認証URLを生成しました。",
        "timestamp": datetime.now().strftime("%Y/%m/%d %H:%M:%S")
    }

    try:
        flow = Flow.from_client_secrets_file(
            get_auth_file_path('credentials.json'),
            scopes=SCOPES,
            redirect_uri=REDIRECT_URI
        )

        auth_url, state = flow.authorization_url(
            access_type='offline',
            include_granted_scopes='true'
        )

        # stateを保存
        oauth_state_path = get_auth_file_path('oauth_state.json')
        with open(oauth_state_path, 'w') as f:
            json.dump({'state': state}, f)

        result["auth_url"] = auth_url

    except Exception as e:
        result["status"] = "error"
        result["message"] = f"認証URL生成中にエラーが発生しました: {str(e)}"

    return result


def authorize_with_code(auth_code: str) -> Dict[str, Any]:
    """
    認証コードを使用してOAuth認証を完了する関数

    Args:
        auth_code (str): 認証コールバックで受け取った認可コード

    Returns:
        Dict[str, Any]: 認証処理の結果を含む辞書
    """
    result = {
        "status": "success",
        "message": "認証完了。トークンが保存されました。",
        "timestamp": datetime.now().strftime("%Y/%m/%d %H:%M:%S")
    }

    try:
        # 保存したstateを読み込む
        oauth_state_path = get_auth_file_path('oauth_state.json')
        if not os.path.exists(oauth_state_path):
            result["status"] = "error"
            result["message"] = "state情報が見つかりません。再度認証を開始してください。"
            return result

        with open(oauth_state_path, 'r') as f:
            json.load(f)['state']

        # Flowを再構築
        flow = Flow.from_client_secrets_file(
            get_auth_file_path('credentials.json'),
            scopes=SCOPES,
            redirect_uri=REDIRECT_URI
        )
        # stateは自動的に内部で使われるので明示的に設定は不要

        # 認証コードを使ってトークンを取得
        flow.fetch_token(code=auth_code)

        creds = flow.credentials

        # トークンを保存
        token_path = get_auth_file_path('token.pickle')
        with open(token_path, 'wb') as token:
            pickle.dump(creds, token)

        # 認証情報をJSONに変換（デバッグ用）
        creds_json = json.loads(creds.to_json())
        result["credentials_info"] = {
            "token_expiry": creds_json.get("token_expiry", ""),
            "scopes": creds_json.get("scopes", [])
        }

        # 一時的なstate情報を削除
        if os.path.exists(oauth_state_path):
            os.remove(oauth_state_path)

    except Exception as e:
        result["status"] = "error"
        result["message"] = f"認証コード処理中にエラーが発生しました: {str(e)}"

    return result


def check_token_validity() -> Dict[str, Any]:
    """
    現在のトークンの有効性をチェックする関数

    Returns:
        Dict[str, Any]: トークンの有効性チェック結果を含む辞書
    """
    result = {
        "status": "success",
        "message": "トークンは有効です。",
        "is_valid": True,
        "timestamp": datetime.now().strftime("%Y/%m/%d %H:%M:%S")
    }

    try:
        creds = get_credentials()

        if not creds:
            result["status"] = "error"
            result["message"] = "認証トークンが見つかりません。認証が必要です。"
            result["is_valid"] = False
            return result

        if not creds.valid:
            result["status"] = "error"
            result["message"] = "認証トークンが無効または期限切れです。再認証が必要です。"
            result["is_valid"] = False
            return result

        # トークンの詳細情報を追加
        if creds.expiry:
            result["token_expiry"] = creds.expiry.strftime("%Y/%m/%d %H:%M:%S")
        else:
            result["token_expiry"] = "期限情報なし"

    except Exception as e:
        result["status"] = "error"
        result["message"] = f"トークン有効性チェック中にエラーが発生しました: {str(e)}"
        result["is_valid"] = False

    return result


def authorize(auth_code: Optional[str] = None) -> Dict[str, Any]:
    """
    Google API 認証を実行し、結果を返す関数

    Args:
        auth_code (Optional[str]): 認証コード。Noneの場合は認証フローを開始

    Returns:
        Dict[str, Any]: 認証処理の結果を含む辞書
    """
    # すでに有効なトークンがあるか確認
    creds = get_credentials()

    if creds and creds.valid:
        # 有効なトークンがある場合はそれを返す
        result = {
            "status": "success",
            "message": "有効な認証トークンがあります。",
            "timestamp": datetime.now().strftime("%Y/%m/%d %H:%M:%S")
        }

        creds_json = json.loads(creds.to_json())
        result["credentials_info"] = {
            "token_expiry": creds_json.get("token_expiry", ""),
            "scopes": creds_json.get("scopes", [])
        }
        return result

    elif auth_code:
        # 認証コードがある場合は、コードを使ってトークンを取得
        return authorize_with_code(auth_code)
    else:
        # それ以外の場合は、認証URLを生成
        return generate_auth_url()
