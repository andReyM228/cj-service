package tables

import (
	"cj_service/internal/domain"
	"github.com/andReyM228/lib/log"
	"net/http"
	"reflect"
	"testing"
)

func TestRepository_Get(t *testing.T) {
	type fields struct {
		log    log.Logger
		client *http.Client
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Students
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				log:    log.Init(),
				client: http.DefaultClient,
			},
			args: args{
				name: "Ко",
			},
			want:    domain.Students{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				log:    tt.fields.log,
				client: tt.fields.client,
			}
			got, err := r.Get(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
