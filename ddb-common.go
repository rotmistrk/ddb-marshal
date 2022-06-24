package ddbmarshal

import (
	"strings"
)

const (
	TagDdb          = "ddb"
	TagItemHashJey  = "hash-key"
	TagItemRangeKey = "range-key"
	TagItemRequired = "required"
	TagItemTtlField = "ttl-ts"
)

type DdbMarshaller struct {
	marshalAllPublicFields     bool
	decapitalizeUntaggedFields bool
	addPrefixToTheFieldNames   string
	// TODO: options:
	//  - should we marshal fields without tags?
	//    - add ighore flag then
	//  - should we require all that are not optional?
	//  - should we use versioned records?
	//  - encryption
	//    - with KMS
	//    - with other instrumentation
	//  - signing
}

func NewMarshaller() *DdbMarshaller {
	return &DdbMarshaller{}
}

func (marshaller *DdbMarshaller) SetMarshalAllPublicFields(value bool) {
	marshaller.marshalAllPublicFields = value
}

func (marshaller *DdbMarshaller) SetDecapitalizeUntaggedFieldNames(value bool) {
	marshaller.decapitalizeUntaggedFields = value
}

func (marshaller *DdbMarshaller) SetFieldNamePrefix(prefix string) {
	marshaller.addPrefixToTheFieldNames = prefix
}

type specs struct {
	name       string
	required   bool
	isHashKey  bool
	isRangeKey bool
	isTtlField bool
}

func ParseDdbTag(tag string) (specs, error) {
	items := strings.Split(tag, ",")
	result := specs{
		name: strings.TrimSpace(items[0]),
	}
	for _, v := range items[1:] {
		switch strings.TrimSpace(v) {
		case TagItemRequired:
			result.required = true
		case TagItemHashJey:
			result.isHashKey = true
		case TagItemRangeKey:
			result.isRangeKey = true
		case TagItemTtlField:
			result.isTtlField = true
		}
	}
	return result, nil
}

func (s specs) IsRequired() bool {
	return s.required
}

func (s specs) IsHashKey() bool {
	return s.isHashKey
}

func (s specs) IsRangeKey() bool {
	return s.isRangeKey
}

func (s specs) IsKey() bool {
	return s.isHashKey || s.IsRequired()
}

func (s specs) IsTtlField() bool {
	return s.isTtlField
}

func (s specs) FieldName() string {
	return s.name
}
