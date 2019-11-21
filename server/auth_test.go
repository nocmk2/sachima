package server

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_helloHandler(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helloHandler(tt.args.c)
		})
	}
}
