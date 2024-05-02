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
                'total_taxes': float(e['税額合計[円/現地通貨]']),
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
        asset_detail = split_asset_into_sections([row for row in csv.reader(f)])[1]

        jsonArray = []

        for i in range(2, len(asset_detail)):
            jsonArray.append({
                'type': asset_detail[i][0],
                'ticker': asset_detail[i][1],
                'name': asset_detail[i][2],
                'account': asset_detail[i][3],
                'count': float(asset_detail[i][4]),
                'average_acquisition_price': float(canonicalize_price(asset_detail[i][6])),
                'current_unit_price': float(canonicalize_price(asset_detail[i][8])),
                'current_price_yen': float(canonicalize_price(asset_detail[i][14])),
                'current_price': float(canonicalize_price(asset_detail[i][15])),
            })

        return jsonArray


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
