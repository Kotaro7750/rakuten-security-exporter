// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: rakuten-security-scraper.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListWithdrawalHistoriesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListWithdrawalHistoriesRequest) Reset() {
	*x = ListWithdrawalHistoriesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWithdrawalHistoriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWithdrawalHistoriesRequest) ProtoMessage() {}

func (x *ListWithdrawalHistoriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWithdrawalHistoriesRequest.ProtoReflect.Descriptor instead.
func (*ListWithdrawalHistoriesRequest) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{0}
}

type ListWithdrawalHistoriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	History []*WithdrawalHistory `protobuf:"bytes,1,rep,name=history,proto3" json:"history,omitempty"`
}

func (x *ListWithdrawalHistoriesResponse) Reset() {
	*x = ListWithdrawalHistoriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWithdrawalHistoriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWithdrawalHistoriesResponse) ProtoMessage() {}

func (x *ListWithdrawalHistoriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWithdrawalHistoriesResponse.ProtoReflect.Descriptor instead.
func (*ListWithdrawalHistoriesResponse) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{1}
}

func (x *ListWithdrawalHistoriesResponse) GetHistory() []*WithdrawalHistory {
	if x != nil {
		return x.History
	}
	return nil
}

type WithdrawalHistory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date     string `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Amount   int64  `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Currency string `protobuf:"bytes,3,opt,name=currency,proto3" json:"currency,omitempty"`
}

func (x *WithdrawalHistory) Reset() {
	*x = WithdrawalHistory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WithdrawalHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WithdrawalHistory) ProtoMessage() {}

func (x *WithdrawalHistory) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WithdrawalHistory.ProtoReflect.Descriptor instead.
func (*WithdrawalHistory) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{2}
}

func (x *WithdrawalHistory) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *WithdrawalHistory) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *WithdrawalHistory) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

type ListDividendHistoriesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListDividendHistoriesRequest) Reset() {
	*x = ListDividendHistoriesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDividendHistoriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDividendHistoriesRequest) ProtoMessage() {}

func (x *ListDividendHistoriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDividendHistoriesRequest.ProtoReflect.Descriptor instead.
func (*ListDividendHistoriesRequest) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{3}
}

type ListDividendHistoriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	History []*DividendHistory `protobuf:"bytes,1,rep,name=history,proto3" json:"history,omitempty"`
}

func (x *ListDividendHistoriesResponse) Reset() {
	*x = ListDividendHistoriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDividendHistoriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDividendHistoriesResponse) ProtoMessage() {}

func (x *ListDividendHistoriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDividendHistoriesResponse.ProtoReflect.Descriptor instead.
func (*ListDividendHistoriesResponse) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{4}
}

func (x *ListDividendHistoriesResponse) GetHistory() []*DividendHistory {
	if x != nil {
		return x.History
	}
	return nil
}

type DividendHistory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date                     string  `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Account                  string  `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Type                     string  `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Ticker                   string  `protobuf:"bytes,4,opt,name=ticker,proto3" json:"ticker,omitempty"`
	Name                     string  `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Currency                 string  `protobuf:"bytes,6,opt,name=currency,proto3" json:"currency,omitempty"`
	Count                    float64 `protobuf:"fixed64,7,opt,name=count,proto3" json:"count,omitempty"`
	DividendUnitprice        float64 `protobuf:"fixed64,8,opt,name=dividend_unitprice,json=dividendUnitprice,proto3" json:"dividend_unitprice,omitempty"`
	DividendTotalBeforeTaxes float64 `protobuf:"fixed64,9,opt,name=dividend_total_before_taxes,json=dividendTotalBeforeTaxes,proto3" json:"dividend_total_before_taxes,omitempty"`
	TotalTaxes               float64 `protobuf:"fixed64,10,opt,name=total_taxes,json=totalTaxes,proto3" json:"total_taxes,omitempty"`
	DividendTotal            float64 `protobuf:"fixed64,11,opt,name=dividend_total,json=dividendTotal,proto3" json:"dividend_total,omitempty"`
}

func (x *DividendHistory) Reset() {
	*x = DividendHistory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DividendHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DividendHistory) ProtoMessage() {}

func (x *DividendHistory) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DividendHistory.ProtoReflect.Descriptor instead.
func (*DividendHistory) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{5}
}

func (x *DividendHistory) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *DividendHistory) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *DividendHistory) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *DividendHistory) GetTicker() string {
	if x != nil {
		return x.Ticker
	}
	return ""
}

func (x *DividendHistory) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DividendHistory) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *DividendHistory) GetCount() float64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *DividendHistory) GetDividendUnitprice() float64 {
	if x != nil {
		return x.DividendUnitprice
	}
	return 0
}

func (x *DividendHistory) GetDividendTotalBeforeTaxes() float64 {
	if x != nil {
		return x.DividendTotalBeforeTaxes
	}
	return 0
}

func (x *DividendHistory) GetTotalTaxes() float64 {
	if x != nil {
		return x.TotalTaxes
	}
	return 0
}

func (x *DividendHistory) GetDividendTotal() float64 {
	if x != nil {
		return x.DividendTotal
	}
	return 0
}

type TotalAssetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TotalAssetRequest) Reset() {
	*x = TotalAssetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TotalAssetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TotalAssetRequest) ProtoMessage() {}

func (x *TotalAssetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TotalAssetRequest.ProtoReflect.Descriptor instead.
func (*TotalAssetRequest) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{6}
}

type TotalAssetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Asset        []*Asset            `protobuf:"bytes,1,rep,name=asset,proto3" json:"asset,omitempty"`
	CurrencyRate []*CurrenyRateToJPY `protobuf:"bytes,2,rep,name=currency_rate,json=currencyRate,proto3" json:"currency_rate,omitempty"`
}

func (x *TotalAssetResponse) Reset() {
	*x = TotalAssetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TotalAssetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TotalAssetResponse) ProtoMessage() {}

func (x *TotalAssetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TotalAssetResponse.ProtoReflect.Descriptor instead.
func (*TotalAssetResponse) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{7}
}

func (x *TotalAssetResponse) GetAsset() []*Asset {
	if x != nil {
		return x.Asset
	}
	return nil
}

func (x *TotalAssetResponse) GetCurrencyRate() []*CurrenyRateToJPY {
	if x != nil {
		return x.CurrencyRate
	}
	return nil
}

type Asset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type                    string  `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Ticker                  string  `protobuf:"bytes,2,opt,name=ticker,proto3" json:"ticker,omitempty"`
	Name                    string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Account                 string  `protobuf:"bytes,4,opt,name=account,proto3" json:"account,omitempty"`
	Count                   float64 `protobuf:"fixed64,5,opt,name=count,proto3" json:"count,omitempty"`
	AverageAcquisitionPrice float64 `protobuf:"fixed64,6,opt,name=average_acquisition_price,json=averageAcquisitionPrice,proto3" json:"average_acquisition_price,omitempty"`
	CurrentUnitPrice        float64 `protobuf:"fixed64,7,opt,name=current_unit_price,json=currentUnitPrice,proto3" json:"current_unit_price,omitempty"`
	CurrentPrice            float64 `protobuf:"fixed64,8,opt,name=current_price,json=currentPrice,proto3" json:"current_price,omitempty"`
	CurrentPriceYen         float64 `protobuf:"fixed64,9,opt,name=current_price_yen,json=currentPriceYen,proto3" json:"current_price_yen,omitempty"`
	Currency                string  `protobuf:"bytes,10,opt,name=currency,proto3" json:"currency,omitempty"`
}

func (x *Asset) Reset() {
	*x = Asset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Asset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Asset) ProtoMessage() {}

func (x *Asset) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Asset.ProtoReflect.Descriptor instead.
func (*Asset) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{8}
}

func (x *Asset) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Asset) GetTicker() string {
	if x != nil {
		return x.Ticker
	}
	return ""
}

func (x *Asset) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Asset) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *Asset) GetCount() float64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *Asset) GetAverageAcquisitionPrice() float64 {
	if x != nil {
		return x.AverageAcquisitionPrice
	}
	return 0
}

func (x *Asset) GetCurrentUnitPrice() float64 {
	if x != nil {
		return x.CurrentUnitPrice
	}
	return 0
}

func (x *Asset) GetCurrentPrice() float64 {
	if x != nil {
		return x.CurrentPrice
	}
	return 0
}

func (x *Asset) GetCurrentPriceYen() float64 {
	if x != nil {
		return x.CurrentPriceYen
	}
	return 0
}

func (x *Asset) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

type CurrenyRateToJPY struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrencyCode string  `protobuf:"bytes,1,opt,name=currencyCode,proto3" json:"currencyCode,omitempty"`
	Rate         float64 `protobuf:"fixed64,2,opt,name=rate,proto3" json:"rate,omitempty"`
}

func (x *CurrenyRateToJPY) Reset() {
	*x = CurrenyRateToJPY{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rakuten_security_scraper_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CurrenyRateToJPY) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CurrenyRateToJPY) ProtoMessage() {}

func (x *CurrenyRateToJPY) ProtoReflect() protoreflect.Message {
	mi := &file_rakuten_security_scraper_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CurrenyRateToJPY.ProtoReflect.Descriptor instead.
func (*CurrenyRateToJPY) Descriptor() ([]byte, []int) {
	return file_rakuten_security_scraper_proto_rawDescGZIP(), []int{9}
}

func (x *CurrenyRateToJPY) GetCurrencyCode() string {
	if x != nil {
		return x.CurrencyCode
	}
	return ""
}

func (x *CurrenyRateToJPY) GetRate() float64 {
	if x != nil {
		return x.Rate
	}
	return 0
}

var File_rakuten_security_scraper_proto protoreflect.FileDescriptor

var file_rakuten_security_scraper_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x72, 0x61, 0x6b, 0x75, 0x74, 0x65, 0x6e, 0x2d, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69,
	0x74, 0x79, 0x2d, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x20, 0x0a, 0x1e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77,
	0x61, 0x6c, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x4f, 0x0a, 0x1f, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72,
	0x61, 0x77, 0x61, 0x6c, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61,
	0x77, 0x61, 0x6c, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x07, 0x68, 0x69, 0x73, 0x74,
	0x6f, 0x72, 0x79, 0x22, 0x5b, 0x0a, 0x11, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x61,
	0x6c, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79,
	0x22, 0x1e, 0x0a, 0x1c, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x4b, 0x0a, 0x1d, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2a, 0x0a, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x48, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x79, 0x52, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x22, 0xe7, 0x02,
	0x0a, 0x0f, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72,
	0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x2d, 0x0a, 0x12, 0x64, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x5f, 0x75, 0x6e,
	0x69, 0x74, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x11, 0x64,
	0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x55, 0x6e, 0x69, 0x74, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x3d, 0x0a, 0x1b, 0x64, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x5f, 0x74, 0x61, 0x78, 0x65, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x18, 0x64, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x54, 0x61, 0x78, 0x65, 0x73, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x74, 0x61, 0x78, 0x65, 0x73, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x54, 0x61, 0x78, 0x65, 0x73,
	0x12, 0x25, 0x0a, 0x0e, 0x64, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x64, 0x69, 0x76, 0x69, 0x64, 0x65,
	0x6e, 0x64, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0x13, 0x0a, 0x11, 0x54, 0x6f, 0x74, 0x61, 0x6c,
	0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x6a, 0x0a, 0x12,
	0x54, 0x6f, 0x74, 0x61, 0x6c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1c, 0x0a, 0x05, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x06, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x05, 0x61, 0x73, 0x73, 0x65, 0x74,
	0x12, 0x36, 0x0a, 0x0d, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x72, 0x61, 0x74,
	0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x79, 0x52, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x4a, 0x50, 0x59, 0x52, 0x0c, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x52, 0x61, 0x74, 0x65, 0x22, 0xce, 0x02, 0x0a, 0x05, 0x41, 0x73, 0x73,
	0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x3a, 0x0a, 0x19, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x61, 0x63,
	0x71, 0x75, 0x69, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x17, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x41, 0x63,
	0x71, 0x75, 0x69, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x2c,
	0x0a, 0x12, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x10, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x12, 0x2a, 0x0a, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x72, 0x69,
	0x63, 0x65, 0x5f, 0x79, 0x65, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x59, 0x65, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x22, 0x4a, 0x0a, 0x10, 0x43, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x79, 0x52, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x4a, 0x50, 0x59, 0x12, 0x22, 0x0a,
	0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x04, 0x72, 0x61, 0x74, 0x65, 0x32, 0x8c, 0x02, 0x0a, 0x16, 0x52, 0x61, 0x6b, 0x75, 0x74, 0x65,
	0x6e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72,
	0x12, 0x5e, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77,
	0x61, 0x6c, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1f, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x61, 0x6c, 0x48, 0x69, 0x73, 0x74,
	0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x61, 0x6c, 0x48, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x58, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1d, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x44,
	0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x0b, 0x54, 0x6f,
	0x74, 0x61, 0x6c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x73, 0x12, 0x12, 0x2e, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e,
	0x54, 0x6f, 0x74, 0x61, 0x6c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x4b, 0x6f, 0x74, 0x61, 0x72, 0x6f, 0x37, 0x37, 0x35, 0x30, 0x2f, 0x72, 0x61,
	0x6b, 0x75, 0x74, 0x65, 0x6e, 0x2d, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2d, 0x65,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rakuten_security_scraper_proto_rawDescOnce sync.Once
	file_rakuten_security_scraper_proto_rawDescData = file_rakuten_security_scraper_proto_rawDesc
)

func file_rakuten_security_scraper_proto_rawDescGZIP() []byte {
	file_rakuten_security_scraper_proto_rawDescOnce.Do(func() {
		file_rakuten_security_scraper_proto_rawDescData = protoimpl.X.CompressGZIP(file_rakuten_security_scraper_proto_rawDescData)
	})
	return file_rakuten_security_scraper_proto_rawDescData
}

var file_rakuten_security_scraper_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_rakuten_security_scraper_proto_goTypes = []interface{}{
	(*ListWithdrawalHistoriesRequest)(nil),  // 0: ListWithdrawalHistoriesRequest
	(*ListWithdrawalHistoriesResponse)(nil), // 1: ListWithdrawalHistoriesResponse
	(*WithdrawalHistory)(nil),               // 2: WithdrawalHistory
	(*ListDividendHistoriesRequest)(nil),    // 3: ListDividendHistoriesRequest
	(*ListDividendHistoriesResponse)(nil),   // 4: ListDividendHistoriesResponse
	(*DividendHistory)(nil),                 // 5: DividendHistory
	(*TotalAssetRequest)(nil),               // 6: TotalAssetRequest
	(*TotalAssetResponse)(nil),              // 7: TotalAssetResponse
	(*Asset)(nil),                           // 8: Asset
	(*CurrenyRateToJPY)(nil),                // 9: CurrenyRateToJPY
}
var file_rakuten_security_scraper_proto_depIdxs = []int32{
	2, // 0: ListWithdrawalHistoriesResponse.history:type_name -> WithdrawalHistory
	5, // 1: ListDividendHistoriesResponse.history:type_name -> DividendHistory
	8, // 2: TotalAssetResponse.asset:type_name -> Asset
	9, // 3: TotalAssetResponse.currency_rate:type_name -> CurrenyRateToJPY
	0, // 4: RakutenSecurityScraper.ListWithdrawalHistories:input_type -> ListWithdrawalHistoriesRequest
	3, // 5: RakutenSecurityScraper.ListDividendHistories:input_type -> ListDividendHistoriesRequest
	6, // 6: RakutenSecurityScraper.TotalAssets:input_type -> TotalAssetRequest
	1, // 7: RakutenSecurityScraper.ListWithdrawalHistories:output_type -> ListWithdrawalHistoriesResponse
	4, // 8: RakutenSecurityScraper.ListDividendHistories:output_type -> ListDividendHistoriesResponse
	7, // 9: RakutenSecurityScraper.TotalAssets:output_type -> TotalAssetResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_rakuten_security_scraper_proto_init() }
func file_rakuten_security_scraper_proto_init() {
	if File_rakuten_security_scraper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rakuten_security_scraper_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListWithdrawalHistoriesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListWithdrawalHistoriesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WithdrawalHistory); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDividendHistoriesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDividendHistoriesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DividendHistory); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TotalAssetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TotalAssetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Asset); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rakuten_security_scraper_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CurrenyRateToJPY); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rakuten_security_scraper_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rakuten_security_scraper_proto_goTypes,
		DependencyIndexes: file_rakuten_security_scraper_proto_depIdxs,
		MessageInfos:      file_rakuten_security_scraper_proto_msgTypes,
	}.Build()
	File_rakuten_security_scraper_proto = out.File
	file_rakuten_security_scraper_proto_rawDesc = nil
	file_rakuten_security_scraper_proto_goTypes = nil
	file_rakuten_security_scraper_proto_depIdxs = nil
}
