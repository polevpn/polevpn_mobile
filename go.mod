module polevpnmobile

go 1.14

require (
	github.com/polevpn/anyvalue v1.0.6
	github.com/polevpn/elog v1.1.1
	github.com/polevpn/polevpn_core v1.2.6
	golang.org/x/mobile v0.0.0-20220112015953-858099ff7816 // indirect
	gorm.io/driver/sqlite v1.4.3 // indirect
	gorm.io/gorm v1.24.2 // indirect
)

replace github.com/polevpn/polevpn_core => ../polevpn_core
