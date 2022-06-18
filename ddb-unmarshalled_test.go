package ddbmarshal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"testing"
)

type testSmall struct {
	Present string `ddb:"present"`
}

func TestDdbMarshaller_GetUnmarshaledFields(t *testing.T) {
	type args struct {
		target   interface{}
		response map[string]*dynamodb.AttributeValue
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*dynamodb.AttributeValue
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "gets unmarshalled",
			args: args{
				target: &testSmall{},
				response: map[string]*dynamodb.AttributeValue{
					"present": {S: aws.String("is present")},
					"absent":  {S: aws.String("is absent")},
				},
			},
			want: map[string]*dynamodb.AttributeValue{
				"absent": {S: aws.String("is absent")},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := &DdbMarshaller{}
			got, err := me.GetUnmarshaledFields(tt.args.target, tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUnmarshaledFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUnmarshaledFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}
