package ddbmarshal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"testing"
)

type testRequired struct {
	Uuid string `ddb:"uuid,required"`
	Name string `ddb:"name"`
}

func TestDdbMarshaller_Unmarshal(t *testing.T) {
	type args struct {
		target interface{}
		source map[string]*dynamodb.AttributeValue
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantData interface{}
	}{
		// TODO: Add test cases.
		{
			name: "happy",
			args: args{
				target: &testDdbMarshal{Ignored: "ignored to be"},
				source: prepareDdb(),
			},
			wantErr:  false,
			wantData: prepareStruct(),
		},
		{
			name: "required is required",
			args: args{
				target: &testRequired{},
				source: map[string]*dynamodb.AttributeValue{
					"name": {S: aws.String("a Name")},
				},
			},
			wantErr:  true,
			wantData: &testRequired{}, // edited (no changes in this case)
		},
		{
			name: "not required is optional",
			args: args{
				target: &testRequired{},
				source: map[string]*dynamodb.AttributeValue{
					"uuid": {S: aws.String("unique id")},
				},
			},
			wantErr:  false,
			wantData: &testRequired{Uuid: "unique id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := &DdbMarshaller{}
			if err := me.Unmarshal(tt.args.target, tt.args.source); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(tt.wantData, tt.args.target) {
				t.Errorf("Unmarshal() gotResult = %v, want %v", tt.args.target, tt.wantData)
			}
		})
	}
}
