# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: rakuten-security-scraper.proto
# Protobuf Python Version: 4.25.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1erakuten-security-scraper.proto\" \n\x1eListWithdrawalHistoriesRequest\"F\n\x1fListWithdrawalHistoriesResponse\x12#\n\x07history\x18\x01 \x03(\x0b\x32\x12.WithdrawalHistory\"C\n\x11WithdrawalHistory\x12\x0c\n\x04\x64\x61te\x18\x01 \x01(\t\x12\x0e\n\x06\x61mount\x18\x02 \x01(\x03\x12\x10\n\x08\x63urrency\x18\x03 \x01(\t\"\x1e\n\x1cListDividendHistoriesRequest\"B\n\x1dListDividendHistoriesResponse\x12!\n\x07history\x18\x01 \x03(\x0b\x32\x10.DividendHistory\"\xeb\x01\n\x0f\x44ividendHistory\x12\x0c\n\x04\x64\x61te\x18\x01 \x01(\t\x12\x0f\n\x07\x61\x63\x63ount\x18\x02 \x01(\t\x12\x0c\n\x04type\x18\x03 \x01(\t\x12\x0e\n\x06ticker\x18\x04 \x01(\t\x12\x0c\n\x04name\x18\x05 \x01(\t\x12\x10\n\x08\x63urrency\x18\x06 \x01(\t\x12\r\n\x05\x63ount\x18\x07 \x01(\x01\x12\x1a\n\x12\x64ividend_unitprice\x18\x08 \x01(\x01\x12#\n\x1b\x64ividend_total_before_taxes\x18\t \x01(\x01\x12\x13\n\x0btotal_taxes\x18\n \x01(\x01\x12\x16\n\x0e\x64ividend_total\x18\x0b \x01(\x01\"\x13\n\x11TotalAssetRequest\"U\n\x12TotalAssetResponse\x12\x15\n\x05\x61sset\x18\x01 \x03(\x0b\x32\x06.Asset\x12(\n\rcurrency_rate\x18\x02 \x03(\x0b\x32\x11.CurrenyRateToJPY\"\xc4\x01\n\x05\x41sset\x12\x0c\n\x04type\x18\x01 \x01(\t\x12\x0e\n\x06ticker\x18\x02 \x01(\t\x12\x0c\n\x04name\x18\x03 \x01(\t\x12\x0f\n\x07\x61\x63\x63ount\x18\x04 \x01(\t\x12\r\n\x05\x63ount\x18\x05 \x01(\x01\x12!\n\x19\x61verage_acquisition_price\x18\x06 \x01(\x01\x12\x1a\n\x12\x63urrent_unit_price\x18\x07 \x01(\x01\x12\x15\n\rcurrent_price\x18\x08 \x01(\x01\x12\x19\n\x11\x63urrent_price_yen\x18\t \x01(\x01\"6\n\x10\x43urrenyRateToJPY\x12\x14\n\x0c\x63urrencyCode\x18\x01 \x01(\t\x12\x0c\n\x04rate\x18\x02 \x01(\x01\x32\x8c\x02\n\x16RakutenSecurityScraper\x12^\n\x17ListWithdrawalHistories\x12\x1f.ListWithdrawalHistoriesRequest\x1a .ListWithdrawalHistoriesResponse\"\x00\x12X\n\x15ListDividendHistories\x12\x1d.ListDividendHistoriesRequest\x1a\x1e.ListDividendHistoriesResponse\"\x00\x12\x38\n\x0bTotalAssets\x12\x12.TotalAssetRequest\x1a\x13.TotalAssetResponse\"\x00\x42\x37Z5github.com/Kotaro7750/rakuten-security-exporter/protob\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'rakuten_security_scraper_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  _globals['DESCRIPTOR']._options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z5github.com/Kotaro7750/rakuten-security-exporter/proto'
  _globals['_LISTWITHDRAWALHISTORIESREQUEST']._serialized_start=34
  _globals['_LISTWITHDRAWALHISTORIESREQUEST']._serialized_end=66
  _globals['_LISTWITHDRAWALHISTORIESRESPONSE']._serialized_start=68
  _globals['_LISTWITHDRAWALHISTORIESRESPONSE']._serialized_end=138
  _globals['_WITHDRAWALHISTORY']._serialized_start=140
  _globals['_WITHDRAWALHISTORY']._serialized_end=207
  _globals['_LISTDIVIDENDHISTORIESREQUEST']._serialized_start=209
  _globals['_LISTDIVIDENDHISTORIESREQUEST']._serialized_end=239
  _globals['_LISTDIVIDENDHISTORIESRESPONSE']._serialized_start=241
  _globals['_LISTDIVIDENDHISTORIESRESPONSE']._serialized_end=307
  _globals['_DIVIDENDHISTORY']._serialized_start=310
  _globals['_DIVIDENDHISTORY']._serialized_end=545
  _globals['_TOTALASSETREQUEST']._serialized_start=547
  _globals['_TOTALASSETREQUEST']._serialized_end=566
  _globals['_TOTALASSETRESPONSE']._serialized_start=568
  _globals['_TOTALASSETRESPONSE']._serialized_end=653
  _globals['_ASSET']._serialized_start=656
  _globals['_ASSET']._serialized_end=852
  _globals['_CURRENYRATETOJPY']._serialized_start=854
  _globals['_CURRENYRATETOJPY']._serialized_end=908
  _globals['_RAKUTENSECURITYSCRAPER']._serialized_start=911
  _globals['_RAKUTENSECURITYSCRAPER']._serialized_end=1179
# @@protoc_insertion_point(module_scope)
