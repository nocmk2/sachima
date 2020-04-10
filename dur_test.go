package main

import (
	"reflect"
	"testing"

	"github.com/nocmk2/sachima/dur"
)

func TestDurReadSQL(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		db   string
		want int
	}{
		{"case1", "select user_name from users limit 1", "sachima_local", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dur.ReadSQL(tt.sql, tt.db); !reflect.DeepEqual(got.Rows(), tt.want) {
				t.Errorf("dur.ReadSQL = %v, want %v", got, tt.want)
			}
		})
	}
}
