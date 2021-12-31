package tools

import "testing"

func TestIsValidIpAddress(t *testing.T) {
	type args struct {
		ipp string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{args: args{ipp: "10.1.23.21"},wantErr: false},
		{args: args{ipp: "10.1.23.256"},wantErr: true},
		{args: args{ipp: "10.1.23.254:80"},wantErr: false},
		{args: args{ipp: "10.1.23.254:324345"},wantErr: true},
		{args: args{ipp: "[2400:da00::dbf:0:100]"},wantErr: false},
		{args: args{ipp: "[2400:da00::dbf:0:100]:80"},wantErr: false},
		{args: args{ipp: "[2400:da00::dbf:0:100]:"},wantErr: true},
		{args: args{ipp: "[2400:da00::dbf:0:100]:234121"},wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IsValidIpAddress(tt.args.ipp); (err != nil) != tt.wantErr {
				t.Errorf("IsValidIpAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
