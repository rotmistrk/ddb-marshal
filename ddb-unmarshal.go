package ddbmarshal

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"strconv"
	"time"
)

func (me *DdbMarshaller) Unmarshal(target interface{}, source map[string]*dynamodb.AttributeValue) error {
	targetValue, err := getValidMarshallingTargetValue(target)
	if err != nil {
		return err
	}
	for i, I := 0, targetValue.NumField(); i < I; i++ {
		fieldType := targetValue.Type().Field(i)
		if ddbSpec, ok := fieldType.Tag.Lookup(TagDdb); ok {
			if !fieldType.IsExported() {
				return errors.New("can't use ddb field for unexported fieldType " + fieldType.Name)
			}
			if specs, err := ParseDdbTag(ddbSpec); err != nil {
				return err
			} else {
				if attrVal := source[specs.name]; attrVal == nil {
					if specs.required {
						return errors.New(fmt.Sprintf("missing required field (gp: %s ddb: %s)", fieldType.Name, specs.name))
					}
				} else {
					fieldValue := targetValue.Field(i)
					switch fieldValue.Interface().(type) {
					case bool:
						fieldValue.Set(reflect.ValueOf(*attrVal.BOOL))
					case string:
						fieldValue.Set(reflect.ValueOf(*attrVal.S))
					case []string:
						fieldValue.Set(reflect.ValueOf(aws.StringValueSlice(attrVal.SS)))
					case int, uint, int64, uint64, float32, float64, time.Time:
						if err := setValueWithParsedNumber(fieldValue, *attrVal.N); err != nil {
							return err
						}
					case []int, []uint, []int64, []uint64, []float32, []float64, []time.Time:
						if err := setValueWithParsedNumbers(fieldValue, attrVal.NS); err != nil {
							return err
						}
					case []byte:
						fieldValue.Set(reflect.ValueOf(attrVal.B))
					case [][]byte:
						fieldValue.Set(reflect.ValueOf(attrVal.BS))
					case
						map[string]string,
						map[string]int,
						map[string]uint,
						map[string]int64,
						map[string]uint64,
						map[string]float32,
						map[string]float64,
						map[string]time.Time:
						if err := setValueWithParsedMap(fieldValue, attrVal.M); err != nil {
							return err
						}
					default:
						return errors.New(fmt.Sprintf("Unsupported field type %v", fieldValue.Interface()))
					}
				}
			}
		}
	}
	return nil
}

func setValueWithParsedMap(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	switch value.Interface().(type) {
	case map[string]string:
		return setMapOfStrings(value, attrs)
	case map[string]int:
		return setMapOfInts(value, attrs)
	case map[string]uint:
		return setMapOfUints(value, attrs)
	case map[string]int64:
		return setMapOfInt64(value, attrs)
	case map[string]uint64:
		return setMapOfUint64(value, attrs)
	case map[string]float32:
		return setMapOfFloat32(value, attrs)
	case map[string]float64:
		return setMapOfFloat64(value, attrs)
	case map[string]time.Time:
		return setMapOfTime(value, attrs)
	default:
		return errors.New(fmt.Sprintf("Unsupported map value type %v", value.Interface()))
	}

}

func setMapOfStrings(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	result := make(map[string]string)
	for k, v := range attrs {
		// TODO: add nil protection
		result[k] = *v.S
	}
	value.Set(reflect.ValueOf(result))
	return nil
}

func setMapOfInts(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	result := make(map[string]int)
	for k, v := range attrs {
		// TODO: add nil protection
		if val, err := strconv.ParseInt(*v.N, 0, 32); err != nil {
			return err
		} else {
			result[k] = int(val)
		}
	}
	value.Set(reflect.ValueOf(result))
	return nil
}

func setMapOfUints(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	result := make(map[string]uint)
	for k, v := range attrs {
		// TODO: add nil protection
		if val, err := strconv.ParseUint(*v.N, 0, 32); err != nil {
			return err
		} else {
			result[k] = uint(val)
		}
	}
	value.Set(reflect.ValueOf(result))
	return nil
}

func setMapOfInt64(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	result := make(map[string]int64)
	for k, v := range attrs {
		// TODO: add nil protection
		if val, err := strconv.ParseInt(*v.N, 0, 64); err != nil {
			return err
		} else {
			result[k] = val
		}
	}
	value.Set(reflect.ValueOf(result))
	return nil
}

func setMapOfUint64(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	result := make(map[string]uint64)
	for k, v := range attrs {
		// TODO: add nil protection
		if val, err := strconv.ParseUint(*v.N, 0, 64); err != nil {
			return err
		} else {
			result[k] = val
		}
	}
	value.Set(reflect.ValueOf(result))
	return nil
}

func setMapOfFloat32(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	result := make(map[string]float32)
	for k, v := range attrs {
		// TODO: add nil protection
		if val, err := strconv.ParseFloat(*v.N, 32); err != nil {
			return err
		} else {
			result[k] = float32(val)
		}
	}
	value.Set(reflect.ValueOf(result))
	return nil
}

func setMapOfFloat64(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	result := make(map[string]float64)
	for k, v := range attrs {
		// TODO: add nil protection
		if val, err := strconv.ParseFloat(*v.N, 64); err != nil {
			return err
		} else {
			result[k] = val
		}
	}
	value.Set(reflect.ValueOf(result))
	return nil
}

func setMapOfTime(value reflect.Value, attrs map[string]*dynamodb.AttributeValue) error {
	result := make(map[string]time.Time)
	for k, v := range attrs {
		// TODO: add nil protection
		if val, err := strconv.ParseInt(*v.N, 0, 64); err != nil {
			return err
		} else {
			result[k] = time.Unix(val, 0)
		}
	}
	value.Set(reflect.ValueOf(result))
	return nil
}

func setValueWithParsedNumbers(value reflect.Value, strings []*string) error {
	switch value.Interface().(type) {
	case []int:
		if result, err := parseIntegers(strings); err != nil {
			return err
		} else {
			value.Set(reflect.ValueOf(result))
		}
	case []uint:
		if result, err := parseUnsignedIntegers(strings); err != nil {
			return err
		} else {
			value.Set(reflect.ValueOf(result))
		}
	case []int64:
		if result, err := parseInt64Slice(strings); err != nil {
			return err
		} else {
			value.Set(reflect.ValueOf(result))
		}
	case []uint64:
		if result, err := parseUint64Slice(strings); err != nil {
			return err
		} else {
			value.Set(reflect.ValueOf(result))
		}
	case []float32:
		if result, err := parseFloat32Slice(strings); err != nil {
			return err
		} else {
			value.Set(reflect.ValueOf(result))
		}
	case []float64:
		if result, err := parseFloat64Slice(strings); err != nil {
			return err
		} else {
			value.Set(reflect.ValueOf(result))
		}
	case []time.Time:
		if result, err := parseTimeSlice(strings); err != nil {
			return err
		} else {
			value.Set(reflect.ValueOf(result))
		}
	default:
		return errors.New(fmt.Sprintf("Unsupported numeric type %v", value.Interface()))
	}
	return nil
}

func parseIntegers(strings []*string) ([]int, error) {
	result := make([]int, 0)
	for _, v := range strings {
		if val, err := strconv.ParseInt(*v, 0, 32); err != nil {
			return nil, err
		} else {
			result = append(result, int(val))
		}
	}
	return result, nil
}

func parseUnsignedIntegers(strings []*string) ([]uint, error) {
	result := make([]uint, 0)
	for _, v := range strings {
		if val, err := strconv.ParseUint(*v, 0, 32); err != nil {
			return nil, err
		} else {
			result = append(result, uint(val))
		}
	}
	return result, nil
}

func parseInt64Slice(strings []*string) ([]int64, error) {
	result := make([]int64, 0)
	for _, v := range strings {
		if val, err := strconv.ParseInt(*v, 0, 64); err != nil {
			return nil, err
		} else {
			result = append(result, val)
		}
	}
	return result, nil
}

func parseUint64Slice(strings []*string) ([]uint64, error) {
	result := make([]uint64, 0)
	for _, v := range strings {
		if val, err := strconv.ParseUint(*v, 0, 64); err != nil {
			return nil, err
		} else {
			result = append(result, val)
		}
	}
	return result, nil
}

func parseFloat32Slice(strings []*string) ([]float32, error) {
	result := make([]float32, 0)
	for _, v := range strings {
		if val, err := strconv.ParseFloat(*v, 32); err != nil {
			return nil, err
		} else {
			result = append(result, float32(val))
		}
	}
	return result, nil
}

func parseFloat64Slice(strings []*string) ([]float64, error) {
	result := make([]float64, 0)
	for _, v := range strings {
		if val, err := strconv.ParseFloat(*v, 64); err != nil {
			return nil, err
		} else {
			result = append(result, val)
		}
	}
	return result, nil
}

func parseTimeSlice(strings []*string) ([]time.Time, error) {
	result := make([]time.Time, 0)
	for _, v := range strings {
		if val, err := strconv.ParseInt(*v, 0, 64); err != nil {
			return nil, err
		} else {
			result = append(result, time.Unix(val, 0))
		}
	}
	return result, nil
}

func setValueWithParsedNumber(value reflect.Value, str string) error {
	if val, err := parseStringToNumber(value.Interface(), str); err != nil {
		return err
	} else {
		value.Set(reflect.ValueOf(val))
		return nil
	}
}

func parseStringToNumber(typ interface{}, str string) (interface{}, error) {
	switch typ.(type) {
	case int:
		if val, err := strconv.ParseInt(str, 0, 32); err != nil {
			return 0, err
		} else {
			return int(val), nil
		}
	case int64:
		if val, err := strconv.ParseInt(str, 0, 64); err != nil {
			return 0, err
		} else {
			return val, nil
		}
	case uint:
		if val, err := strconv.ParseUint(str, 0, 32); err != nil {
			return 0, err
		} else {
			return uint(val), nil
		}
	case uint64:
		if val, err := strconv.ParseUint(str, 0, 64); err != nil {
			return 0, err
		} else {
			return val, nil
		}
	case float32:
		if val, err := strconv.ParseFloat(str, 32); err != nil {
			return 0, err
		} else {
			return float64(val), nil
		}
	case float64:
		if val, err := strconv.ParseFloat(str, 64); err != nil {
			return 0, err
		} else {
			return val, nil
		}
	case time.Time:
		if val, err := strconv.ParseInt(str, 0, 64); err != nil {
			return 0, err
		} else {
			return time.Unix(val, 0).UTC(), nil
		}
	default:
		return typ, errors.New(fmt.Sprintf("Unsupported numeric type %v", typ))
	}
}
