package processor

import (
	"testing"
)

func Test_process(t *testing.T) {
	type data struct {
		Field1 string
		field2 string
	}
	type args struct {
		name    string
		pattern string
		data    interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				name:    "hp",
				pattern: "{{.Field1}}",
				data: data{
					Field1: "value1",
					field2: "value2",
				},
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "unexported field",
			args: args{
				name:    "uf",
				pattern: "{{.field2}}",
				data: data{
					Field1: "value1",
					field2: "value2",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty pattern",
			args: args{
				name:    "ep",
				pattern: "",
				data: data{
					Field1: "value1",
					field2: "value2",
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "empty data",
			args: args{
				name:    "ed",
				pattern: "{{.SomeField}}",
				data:    nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := process(tt.args.name, tt.args.pattern, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("process() got = %v, want %v", got, tt.want)
			}
		})
	}
}
