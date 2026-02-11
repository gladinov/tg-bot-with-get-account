package domain

type ValuesMoex struct {
	ShortName       NullString
	TradeDate       NullString
	MaturityDate    NullString
	OfferDate       NullString
	BuybackDate     NullString
	YieldToMaturity NullFloat64
	YieldToOffer    NullFloat64
	FaceValue       NullFloat64
	FaceUnit        NullString
	Duration        NullFloat64
}

type NullString struct {
	Value  string
	IsSet  bool
	IsNull bool
}

func NewNullString(value string, isSet bool, isNull bool) NullString {
	return NullString{
		Value:  value,
		IsSet:  isSet,
		IsNull: isNull,
	}
}

func (ns NullString) GetValue() string {
	return ns.Value
}

func (ns NullString) GetIsSet() bool {
	return ns.IsSet
}

func (ns NullString) GetIsNull() bool {
	return ns.IsNull
}

type NullFloat64 struct {
	Value  float64
	IsSet  bool
	IsNull bool
}

func NewNullFloat64(value float64, isSet bool, isNull bool) NullFloat64 {
	return NullFloat64{
		Value:  value,
		IsSet:  isSet,
		IsNull: isNull,
	}
}

func (nf NullFloat64) GetValue() float64 {
	return nf.Value
}

func (nf NullFloat64) GetIsSet() bool {
	return nf.IsSet
}

func (nf NullFloat64) GetIsNull() bool {
	return nf.IsNull
}

func (nf NullFloat64) IsHasValue() bool {
	if !nf.IsSet || nf.IsNull {
		return false
	}
	return true
}

func (ns NullString) IsHasValue() bool {
	if !ns.IsSet || ns.IsNull {
		return false
	}
	return true
}
