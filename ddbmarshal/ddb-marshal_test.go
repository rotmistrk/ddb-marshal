package ddbmarshal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type empty struct{}

type testDdbMarshal struct {
	Ignored    string
	Name       string             `ddb:"Name"`
	Active     bool               `ddb:"active"`
	Ordinal    int                `ddb:"ordinal"`
	Value      float64            `ddb:"value"`
	Groups     []string           `ddb:"groups"`
	Counts     []uint64           `ddb:"counts"`
	Properties map[string]string  `ddb:"properties"`
	Numbers    map[string]float64 `ddb:"numbers"`
	Bytes      []byte             `ddb:"bytes"`
	BytesList  [][]byte           `ddb:"bytesList"`
	Expire     time.Time          `ddb:"expire"`
}

const (
	THE_TIME = "2022-02-02T22:02:20Z"
)

func mustParseTime(t string) time.Time {
	if tm, err := time.Parse(time.RFC3339, t); err != nil {
		log.Panicln("Failed to parse ", t)
		return tm
	} else {
		return tm
	}
}

func prepareStruct() *testDdbMarshal {
	return &testDdbMarshal{
		Ignored: "ignored to be",
		Name:    "a Name",
		Active:  true,
		Ordinal: 55,
		Value:   3.14159,
		Groups:  []string{"admins", "users"},
		Counts:  []uint64{22, 33},
		Properties: map[string]string{
			"hello": "world",
		},
		Numbers: map[string]float64{
			"pi":     3.14159,
			"e":      2.71,
			"golden": 1.61,
		},
		Bytes:     []byte("Oopsy"),
		BytesList: [][]byte{[]byte("aaa"), []byte("bbb")},
		Expire:    mustParseTime(THE_TIME),
	}
}

func prepareDdb() map[string]*dynamodb.AttributeValue {
	result := make(map[string]*dynamodb.AttributeValue)
	result["Name"] = &dynamodb.AttributeValue{S: aws.String("a Name")}
	result["active"] = &dynamodb.AttributeValue{BOOL: aws.Bool(true)}
	result["ordinal"] = &dynamodb.AttributeValue{N: aws.String("55")}
	result["value"] = &dynamodb.AttributeValue{N: aws.String("3.14159")}
	result["groups"] = &dynamodb.AttributeValue{SS: aws.StringSlice([]string{"admins", "users"})}
	result["counts"] = &dynamodb.AttributeValue{NS: aws.StringSlice([]string{"22", "33"})}
	result["properties"] = &dynamodb.AttributeValue{
		M: map[string]*dynamodb.AttributeValue{
			"hello": {S: aws.String("world")},
		},
	}
	result["numbers"] = &dynamodb.AttributeValue{
		M: map[string]*dynamodb.AttributeValue{
			"pi":     {N: aws.String("3.14159")},
			"e":      {N: aws.String("2.71")},
			"golden": {N: aws.String("1.61")},
		},
	}
	result["bytes"] = &dynamodb.AttributeValue{
		B: []byte("Oopsy"),
	}
	result["bytesList"] = &dynamodb.AttributeValue{
		BS: [][]byte{[]byte("aaa"), []byte("bbb")},
	}
	result["expire"] = &dynamodb.AttributeValue{
		N: aws.String(strconv.FormatInt(mustParseTime(THE_TIME).Unix(), 10)),
	}
	return result
}

func TestDdbMarshaller_Marshal(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantResult map[string]*dynamodb.AttributeValue
		wantErr    bool
	}{
		// TODO: Add more test cases.
		{
			name:       "empty",
			args:       args{&empty{}},
			wantResult: make(map[string]*dynamodb.AttributeValue),
			wantErr:    false,
		},
		{
			name:       "happy",
			args:       args{prepareStruct()},
			wantResult: prepareDdb(),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := &DdbMarshaller{}
			gotResult, err := me.Marshal(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Marshal() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
