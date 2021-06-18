package i

import "testing"

func TestCheckconfigFile(t *testing.T) {
	type args struct {
		configFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"A01", args{configFile: "a.yml"}, true},
		{"A02", args{configFile: "../config/run.yml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkconfigFile(tt.args.configFile); (err != nil) != tt.wantErr {
				t.Errorf("CheckconfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
