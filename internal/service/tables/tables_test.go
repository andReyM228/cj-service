package tables

import (
	"testing"
)

func TestService_GetDate(t *testing.T) {
	type args struct {
		day string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				day: "Воскресенье",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			got, err := s.GetDate(tt.args.day)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(got)
		})
	}
}
