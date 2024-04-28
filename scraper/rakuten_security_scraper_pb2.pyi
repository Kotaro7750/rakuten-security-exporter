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
