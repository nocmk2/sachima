module github.com/nocmk2/sachima

replace github.com/nocmk2/score => /Users/zhangmk/go/src/github.com/nocmk2/score

replace github.com/nocmk2/sachima/auth => /Users/zhangmk/go/src/github.com/nocmk2/sachima/auth

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43

go 1.13

require (
	github.com/appleboy/gin-jwt/v2 v2.6.2
	github.com/gin-gonic/gin v1.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nocmk2/score v1.0.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
)
