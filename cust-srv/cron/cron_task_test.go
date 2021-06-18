package cron

import (
	"a.a/mp-server/common/cache"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
	"time"
)

func TestGetDistributedLock(t *testing.T) {
	type args struct {
		key    string
		value  string
		expire time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"A01", args{key: "myLockKey", value: "myLockValue", expire: time.Second * 60}, true},
		{"A02", args{key: "myLockKey", value: "myLockValue", expire: time.Second * 60}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cache.GetDistributedLock(tt.args.key, tt.args.value, tt.args.expire); got != tt.want {
				t.Errorf("GetDistributedLock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReleaseDistributedLock(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"A01", args{key: "myLockKey", value: "myLockValueXXXXX"}, false},
		{"A02", args{key: "myLockKey", value: "myLockValue"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cache.ReleaseDistributedLock(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("ReleaseDistributedLock() = %v, want %v", got, tt.want)
			}
		})
	}
}
