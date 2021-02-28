module github.com/moedersvoormoeders/api.mvm.digital

go 1.14

replace github.com/schmorrison/Zoho v0.0.0-20200726181448-707d9fdc8ca7 => github.com/meyskens/Zoho v0.0.0-20200903081837-b98904914dd2

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/didip/tollbooth v4.0.2+incompatible
	github.com/didip/tollbooth_echo v0.0.0-20190918161726-5adbfff23d88
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/jinzhu/gorm v1.9.15
	github.com/labstack/echo/v4 v4.1.16
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/schmorrison/Zoho v0.0.0-20200726181448-707d9fdc8ca7
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/valyala/fasttemplate v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200728102440-3e129f6d46b1 // indirect
	gorm.io/driver/postgres v1.0.0
	gorm.io/gorm v1.20.12
)
