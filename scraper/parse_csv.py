import csv
import json
import pathlib
import re


# 日本円の「,」や「123 USD」などを数字のみに直す
def canonicalize_price(price):
    return re.sub('[^0-9.]', '', price) or '0'


def parse_dividend_history(download_dir):
    with open(pathlib.Path(download_dir, 'dividend.csv'), encoding='sjis') as f:
        jsonArray = []
        for row in csv.DictReader(f):
            jsonArray.append(row)

        convertedJsonArray = []

        for e in jsonArray:
            if e['受取通貨'] == 'USドル':
                currency = 'USD'
            else:
                currency = 'YEN'

            convertedJsonArray.append({
                'date': e['入金日'],
                'account': e['口座'],
                'type': e['商品'],
                'ticker': e['銘柄コード'],
                'name': e['銘柄'],
                'currency': currency,
                'count': int(e['数量[株/口]']),
                'dividend_unitprice': float(e['単価[円/現地通貨]']),
                'dividend_total_before_taxes': float(e['配当・分配金合計（税引前）[円/現地通貨]']),
                # 「税額合計[円/現地通貨]」は国内での課税分のみ表示され海外税額は含まれず、円以外の通貨で受け取る場合には「-」で表示される
                # そのためこの項目ではなく税引き前から受け取り金額を引くことで計算する
                'total_taxes': float(e['配当・分配金合計（税引前）[円/現地通貨]']) - float(e['受取金額[円/現地通貨]']),
                'dividend_total': float(e['受取金額[円/現地通貨]'])
            })

        return convertedJsonArray


def parse_withdrawal_history(download_dir):
    # [ { 'date': '1970/01/01', 'amount': 100, 'type': 'in'}, ... ] という形式に変換する
    with open(pathlib.Path(download_dir, 'withdrawal.csv'), encoding='sjis') as f:
        reader = csv.reader(f)

        next(reader)
        next(reader)
        next(reader)
        next(reader)

        jsonArray = []

        for row in csv.DictReader(f):
            jsonArray.append(row)

        return [
            {'date': e['入出金日'], 'amount': int(e['入金額[円]']), 'currency': 'JPY'} if e['入金額[円]']
            else {'date': e['入出金日'], 'amount': -1 * int(e['出金額[円]']), 'currency': 'JPY'}
            for e in jsonArray
        ]


def parse_asset(download_dir):
    with open(pathlib.Path(download_dir, 'asset.csv'), encoding='sjis') as f:
        sections = split_asset_into_sections([row for row in csv.reader(f)])

        asset_detail = sections[1]
        asset_array = []

        for i in range(2, len(asset_detail)):
            asset_type = asset_detail[i][0]

            currency = asset_detail[i][7]
            # 外貨預り金の場合「円/通貨名」という記載になるため
            if currency == "円" or asset_type == "外貨預り金":
                currency = 'JPY'

            count = float(canonicalize_price(asset_detail[i][4]))
            average_acquisition_price = float(
                canonicalize_price(asset_detail[i][6]))
            current_unit_price = float(canonicalize_price(asset_detail[i][8]))
            current_price_yen = float(canonicalize_price(asset_detail[i][14]))
            # 円建て資産の場合外貨建ての時価評価額は使えない
            current_price = float(canonicalize_price(
                asset_detail[i][14 if currency == 'JPY' else 15]))

            if asset_type == '投資信託':
                # 投資信託では平均取得額・現在値が1口あたりの金額ではなく基準価額となっているため1口あたりの金額に修正
                average_acquisition_price = current_price / current_unit_price * average_acquisition_price / count
                current_unit_price = current_price / count
            elif asset_type == '外貨預り金':
                # 外貨預り金は平均取得金額が出ていないため現在の金額で代用する
                # 正確に計算しようにも、外貨建て資産を売ったときなど計算不可能
                average_acquisition_price = current_unit_price

            asset_array.append({
                'type': asset_type,
                'ticker': asset_detail[i][1],
                'name': asset_detail[i][2],
                'currency': currency,
                'account': asset_detail[i][3],
                'count': count,
                'average_acquisition_price': average_acquisition_price,
                'current_unit_price': current_unit_price,
                'current_price_yen': current_price_yen,
                'current_price': current_price
            })

        currency_rate = sections[2]
        currency_rate_array = []

        for i in range(1, len(currency_rate)):
            rate = currency_rate[i][1]
            currency_code = currency_rate[i][2].split('/')[1]

            currency_rate_array.append({
                "currency_code": currency_code,
                "rate": float(rate)
            })

        return asset_array, currency_rate_array


# 「資産合計」「保有資産詳細」「参考為替レート」の3つに分解する
def split_asset_into_sections(asset_lines):
    sections = []

    section = []
    for line in asset_lines:
        if len(line) == 0:
            continue
        else:
            if '■' in line[0] and len(section) != 0:
                sections.append(section)
                section = []
            section.append(line)

    if len(section) != 0:
        sections.append(section)

    return sections
