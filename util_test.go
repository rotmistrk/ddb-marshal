package ddbmarshal

import (
	"reflect"
	"testing"
)

func Test_getValidMarshallingTargetValue(t *testing.T) {
	var scalar string
	var entry struct {
		Name string
	}
	type args struct {
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			"scalar fails",
			args{scalar},
			reflect.Value{},
			true,
		},
		{
			"scalar ref fails",
			args{&scalar},
			reflect.Value{},
			true,
		},
		{
			"struct fails ",
			args{entry},
			reflect.Value{},
			true,
		},
		{
			"struct ref works ",
			args{&entry},
			reflect.Indirect(reflect.ValueOf(&entry).Elem()),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValidMarshallingTargetValue(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValidMarshallingTargetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValidMarshallingTargetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
