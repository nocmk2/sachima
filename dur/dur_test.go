package dur

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// func TestCol_Rank(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		c    Col
// 		want Col
// 	}{
// 		{"case1", Col{1, 2, 3, 4}, Col{0, 1, 2, 3}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.c.Rank(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Col.Rank() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCol_Sort(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		c    Col
// 		want Col
// 	}{
// 		{"case1", Col{2, 3, 1, 4}, Col{1, 2, 3, 4}},
// 		{"case2", Col{2.34, 3.24, 1.34, 4.99}, Col{1.34, 2.34, 3.24, 4.99}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.c.Sort(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Col.Sort() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCol_Percentile(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		c    Col
// 		want Col
// 	}{
// 		// {"case1", Col{2, 3, 1, 4}, Col{1, 2, 3, 4}},
// 		{"case2", Col{2.34, 3.24, 1.34, 4.99}, Col{0.0, 0.25, 0.5, 0.75}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.c.Percentile(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Col.Percentile() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
