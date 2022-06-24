package ddbmarshal

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (me *DdbMarshaller) Marshal(source interface{}) (result map[string]*dynamodb.AttributeValue, err error) {
	return me.MarshalTagFilter(source, func(specs) bool {
		return true
	})
}

func IsRequired(spec specs) bool {
	return spec.required
}

func IsHashKey(spec specs) bool {
	return spec.isHashKey
}

func IsRangeKey(spec specs) bool {
	return spec.isRangeKey
}

func IsKeyField(spec specs) bool {
	return spec.isHashKey || spec.isRangeKey
}

func (me *DdbMarshaller) MarshalTagFilter(source interface{}, filter func(spec specs) bool) (result map[string]*dynamodb.AttributeValue, err error) {
	sourceValue, err := getValidMarshallingTargetValue(source)
	if err != nil {
		return nil, err
	}
	result = make(map[string]*dynamodb.AttributeValue)
	for i, I := 0, sourceValue.NumField(); i < I; i++ {
		fieldType := sourceValue.Type().Field(i)
		checkUntagged := me.marshalAllPublicFields && fieldType.IsExported()
		if ddbSpecStr, ok := fieldType.Tag.Lookup(TagDdb); ok || checkUntagged {
			// TODO: refactor me
			// TODO: add unit tests on marshaller flags
			if !fieldType.IsExported() {
				return nil, errors.New("can't use ddb field for unexported fieldType " + fieldType.Name)
			}
			var ddbSpecs specs
			if !ok {
				if checkUntagged {
					ddbSpecs = specs{name: fieldType.Name}
					if me.decapitalizeUntaggedFields {
						ddbSpecs.name = strings.ToLower(ddbSpecs.name[0:1]) + ddbSpecs.name[1:]
					}
				} else {
					continue
				}
			} else if ddbSpecs, err = ParseDdbTag(ddbSpecStr); err != nil {
				return nil, err
			}
			ddbSpecs.name = me.addPrefixToTheFieldNames + ddbSpecs.name
			if filter(ddbSpecs) {
				fieldValue := sourceValue.Field(i)
				if result[ddbSpecs.name], err = ddbBasicMarshal(fieldValue); err != nil {
					return nil, err
				}
			}
		}
	}
	return result, nil
}

func ddbBasicMarshal(value reflect.Value) (*dynamodb.AttributeValue, error) {
	switch value := value.Interface().(type) {
	case bool:
		return &dynamodb.AttributeValue{BOOL: aws.Bool(value)}, nil
	case string:
		return &dynamodb.AttributeValue{S: aws.String(value)}, nil
	case []string:
		return &dynamodb.AttributeValue{SS: aws.StringSlice(value)}, nil
	case []byte:
		return &dynamodb.AttributeValue{B: value}, nil
	case [][]byte:
		return &dynamodb.AttributeValue{BS: value}, nil
	case
		map[string]string,
		map[string]int,
		map[string]uint,
		map[string]int64,
		map[string]uint64,
		map[string]float32,
		map[string]float64,
		map[string]time.Time:
		if theMap, err := ddbMarshalMap(value); err != nil {
			return nil, err
		} else {
			return &dynamodb.AttributeValue{M: theMap}, nil
		}
	case int, int64, uint, uint64, float32, float64, time.Time:
		if str, err := ddbFormatNum(value); err != nil {
			return nil, err
		} else {
			return &dynamodb.AttributeValue{N: aws.String(str)}, nil
		}
	case []int, []int64, []uint, []uint64, []float32, []float64, []time.Time:
		if strs, err := ddbFormatNums(value); err != nil {
			return nil, err
		} else {
			return &dynamodb.AttributeValue{NS: aws.StringSlice(strs)}, nil
		}
	default:
		return nil, errors.New(fmt.Sprintf("Can't format type %t", value))
	}
}

func ddbFormatNums(value interface{}) ([]string, error) {
	result := make([]string, 0)
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(value)

		for i := 0; i < slice.Len(); i++ {
			if str, err := ddbFormatNum(slice.Index(i).Interface()); err != nil {
				return nil, err
			} else {
				result = append(result, str)
			}
		}
	default:
		return nil, errors.New(fmt.Sprintf("slice of numbers is expected, got %t", value))
	}
	return result, nil
}

func ddbFormatNum(val interface{}) (string, error) {
	switch val := val.(type) {
	case int:
		return strconv.Itoa(val), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'G', -1, 64), nil
	case float64:
		return strconv.FormatFloat(val, 'G', -1, 64), nil
	case time.Time:
		return strconv.FormatInt(val.Unix(), 10), nil
	default:
		return "", errors.New(fmt.Sprintf("Unsupported number type: %t", val))
	}
}

func ddbMarshalMap(value interface{}) (result map[string]*dynamodb.AttributeValue, err error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		result = make(map[string]*dynamodb.AttributeValue)
		valueMap := reflect.ValueOf(value)
		iter := valueMap.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			if result[k.String()], err = ddbBasicMarshal(v); err != nil {
				return nil, err
			}
		}
	default:
		return nil, errors.New(fmt.Sprintf("map[string] is expected, got %t", value))
	}
	return
}
