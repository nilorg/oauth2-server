module github.com/nilorg/oauth2-server

go 1.12

require (
	github.com/gin-contrib/sessions v0.0.1
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/jinzhu/gorm v1.9.10
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/nilorg/oauth2 v0.0.0-20190825141224-fdc11fb53ecf
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/ugorji/go v1.1.7 // indirect
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7 // indirect
)

replace github.com/nilorg/oauth2 => ../oauth2
