package ddbmarshal

import (
	"reflect"
	"testing"
)

func TestNewMarshaller(t *testing.T) {
	tests := []struct {
		name string
		want *DdbMarshaller
	}{
		{"construct", &DdbMarshaller{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMarshaller(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMarshaller() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSpecs(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name    string
		args    args
		want    specs
		wantErr bool
	}{
		{
			"name only",
			args{
				"myColumn",
			},
			specs{
				"myColumn",
				false,
			},
			false,
		},
		{
			"name and required ",
			args{
				"myColumn,required",
			},
			specs{
				"myColumn",
				true,
			},
			false,
		},
		{
			"name, something, and  required ",
			args{
				"myColumn,something, required",
			},
			specs{
				"myColumn",
				true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSpecs(tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSpecs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSpecs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_specs_IsRequired(t *testing.T) {
	type fields struct {
		name     string
		required bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"required is required",
			fields{"field", true},
			true,
		},
		{
			"not required is not required",
			fields{"field", false},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := specs{
				name:     tt.fields.name,
				required: tt.fields.required,
			}
			if got := s.IsRequired(); got != tt.want {
				t.Errorf("IsRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}
