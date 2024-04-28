from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListWithdrawalHistoriesRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ListWithdrawalHistoriesResponse(_message.Message):
    __slots__ = ("history",)
    HISTORY_FIELD_NUMBER: _ClassVar[int]
    history: _containers.RepeatedCompositeFieldContainer[WithdrawalHistory]
    def __init__(self, history: _Optional[_Iterable[_Union[WithdrawalHistory, _Mapping]]] = ...) -> None: ...

class WithdrawalHistory(_message.Message):
    __slots__ = ("date", "amount", "type", "currency")
    DATE_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    CURRENCY_FIELD_NUMBER: _ClassVar[int]
    date: str
    amount: int
    type: str
    currency: str
    def __init__(self, date: _Optional[str] = ..., amount: _Optional[int] = ..., type: _Optional[str] = ..., currency: _Optional[str] = ...) -> None: ...

class ListDividendHistoriesRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ListDividendHistoriesResponse(_message.Message):
    __slots__ = ("history",)
    HISTORY_FIELD_NUMBER: _ClassVar[int]
    history: _containers.RepeatedCompositeFieldContainer[DividendHistory]
    def __init__(self, history: _Optional[_Iterable[_Union[DividendHistory, _Mapping]]] = ...) -> None: ...

class DividendHistory(_message.Message):
    __slots__ = ("date", "account", "type", "ticker", "name", "currency", "count", "dividend_unitprice", "dividend_total_before_taxes", "total_taxes", "dividend_total")
    DATE_FIELD_NUMBER: _ClassVar[int]
    ACCOUNT_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    TICKER_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    CURRENCY_FIELD_NUMBER: _ClassVar[int]
    COUNT_FIELD_NUMBER: _ClassVar[int]
    DIVIDEND_UNITPRICE_FIELD_NUMBER: _ClassVar[int]
    DIVIDEND_TOTAL_BEFORE_TAXES_FIELD_NUMBER: _ClassVar[int]
    TOTAL_TAXES_FIELD_NUMBER: _ClassVar[int]
    DIVIDEND_TOTAL_FIELD_NUMBER: _ClassVar[int]
    date: str
    account: str
    type: str
    ticker: str
    name: str
    currency: str
    count: float
    dividend_unitprice: float
    dividend_total_before_taxes: float
    total_taxes: float
    dividend_total: float
    def __init__(self, date: _Optional[str] = ..., account: _Optional[str] = ..., type: _Optional[str] = ..., ticker: _Optional[str] = ..., name: _Optional[str] = ..., currency: _Optional[str] = ..., count: _Optional[float] = ..., dividend_unitprice: _Optional[float] = ..., dividend_total_before_taxes: _Optional[float] = ..., total_taxes: _Optional[float] = ..., dividend_total: _Optional[float] = ...) -> None: ...
