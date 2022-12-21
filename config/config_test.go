package config

import (
	"reflect"
	"testing"
	"text/template"
	"text/template/parse"
)

func Test_extractVariables(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr error
	}{
		{name: "extract one variable",
			args: args{
				command: "curl something.com:{{.port}}",
			},
			want: []string{"port"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, wantErr := extractVariables(tt.args.command)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractVariables() = %v, want %v", got, tt.want)
			}
			if wantErr != tt.wantErr {
				t.Errorf("extractVariables() = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestService_Valid(t *testing.T) {
	type fields struct {
		Command     string
		Environment string
		Enable      bool
		Variables   []map[string]string
	}
	t1 := make(map[string]string)
	t1["port"] = "1001"

	t2 := make(map[string]string)
	t2["port"] = "1002"

	tests := []struct {
		name    string
		service Service
		want    bool
	}{
		{
			name:    "disabled service",
			service: Service{Command: "something {{.port}}", Enable: false, Variables: []map[string]string{t1, t2}},
			want:    false,
		},
		{
			name:    "enabled service, wrong variable count",
			service: Service{Command: "something {{.port}} and some other variable {{.testing}}", Enable: true, Variables: []map[string]string{t1}},
			want:    false,
		},
		{
			name:    "enabled service",
			service: Service{Command: "something {{.Port}}", Enable: true, Variables: []map[string]string{t1}},
			want:    true,
		},
		{
			name:    "enabled service, but missing variable",
			service: Service{Command: "something {{.trop}}", Enable: true, Variables: []map[string]string{t1}},
			want:    false,
		},
		{
			name:    "enabled service, but invalid variable syntax",
			service: Service{Command: "something {trop}}", Enable: true, Variables: []map[string]string{t1}},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				Command:     tt.service.Command,
				Environment: tt.service.Environment,
				Enable:      tt.service.Enable,
				Variables:   tt.service.Variables,
			}
			if got := s.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_Valid(t *testing.T) {
	// Configuration with two services
	c1 := make(map[string]Service)

	// Variable 1
	v1 := make(map[string]string)
	v1["port"] = "1234"

	// Variable 2
	v2 := make(map[string]string)
	v2["port"] = "5678"

	c1["service-1"] = Service{
		Variables: []map[string]string{v1},
	}
	c1["service-2"] = Service{
		Variables: []map[string]string{v2},
	}

	c2 := make(map[string]Service)
	c2["service-1"] = Service{
		Variables: []map[string]string{v1},
	}
	c2["service-2"] = Service{
		Variables: []map[string]string{v1},
	}

	tests := []struct {
		name string
		s    Configuration
		want bool
	}{
		{
			name: "two services with no overlapping port variable",
			s:    Configuration{Services: c1},
			want: true,
		},
		{
			name: "two services with no overlapping port variable",
			s:    Configuration{Services: c2},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListTemplateFields(t *testing.T) {
	type args struct {
		t *template.Template
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListTemplateFields(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTemplateFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_InterpolatedCommand(t *testing.T) {
	type fields struct {
		Command     string
		Environment string
		Enable      bool
		Variables   []map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				Command:     tt.fields.Command,
				Environment: tt.fields.Environment,
				Enable:      tt.fields.Enable,
				Variables:   tt.fields.Variables,
			}
			got, err := s.InterpolatedCommand()
			if (err != nil) != tt.wantErr {
				t.Errorf("InterpolatedCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InterpolatedCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Valid1(t *testing.T) {
	type fields struct {
		Command     string
		Environment string
		Enable      bool
		Variables   []map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				Command:     tt.fields.Command,
				Environment: tt.fields.Environment,
				Enable:      tt.fields.Enable,
				Variables:   tt.fields.Variables,
			}
			if got := s.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractVariables1(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractVariables(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractVariables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listNodeFields(t *testing.T) {
	type args struct {
		node parse.Node
		res  []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listNodeFields(tt.args.node, tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listNodeFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
