package ddbmarshal

import (
	"strings"
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
	name     string
	required bool
}

func parseSpecs(tag string) (specs, error) {
	items := strings.Split(tag, ",")
	result := specs{
		name: items[0],
	}
	for _, v := range items[1:] {
		switch strings.TrimSpace(v) {
		case "required":
			result.required = true
		}
	}
	return result, nil
}

func (s specs) IsRequired() bool {
	return s.required
}
