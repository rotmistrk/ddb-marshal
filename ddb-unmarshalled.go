package ddbmarshal

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
)

func (me *DdbMarshaller) GetUnmarshaledFields(target interface{}, response map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	targetValue, err := getValidMarshallingTargetValue(target)
	if err != nil {
		return nil, err
	}
	fieldmap := make(map[string]*reflect.StructField)
	for i, I := 0, targetValue.NumField(); i < I; i++ {
		fieldType := targetValue.Type().Field(i)
		if ddbSpec, ok := fieldType.Tag.Lookup("ddb"); ok {
			if !fieldType.IsExported() {
				return nil, errors.New("can't use ddb field for unexported fieldType " + fieldType.Name)
			}
			if specs, err := parseSpecs(ddbSpec); err != nil {
				return nil, err
			} else {
				fieldmap[specs.name] = &fieldType
			}
		}
	}

	result := make(map[string]*dynamodb.AttributeValue)
	for k, v := range response {
		if _, ok := fieldmap[k]; !ok {
			result[k] = v
		}
	}
	return result, nil
}
