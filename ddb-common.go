package ddbmarshal

import (
	"strings"
)

const (
	TagDdb          = "ddb"
	TagItemHashJey  = "hash-key"
	TagItemRangeKey = "range-key"
	TagItemRequired = "required"
)

type DdbMarshaller struct {
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

type specs struct {
	name       string
	required   bool
	isHashKey  bool
	isRangeKey bool
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
		}
	}
	return result, nil
}

func (s specs) IsRequired() bool {
	return s.required
}
