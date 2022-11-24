package polevpnmobile

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

type accessServer struct {
	gorm.Model
	Name                string
	Endpoint            string
	User                string
	Password            string
	Sni                 string
	SkipVerifySSL       bool
	UseRemoteRouteRules bool
	LocalRouteRules     string
	ProxyDomains        string
}

func InitDB(path string) error {

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&accessServer{})
	Db = db
	return nil
}

func addAccessServer(server accessServer) error {
	ret := Db.Create(&server)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func getAllAccessServer() ([]accessServer, error) {
	var records []accessServer
	ret := Db.Find(&records)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return records, nil
}

func deleteAccessServer(ID uint) error {
	var server accessServer
	ret := Db.Delete(&server, ID)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func updateAccessServer(server accessServer) error {

	ret := Db.Model(&server).Updates(&map[string]interface{}{
		"ID":                  server.ID,
		"Name":                server.Name,
		"Endpoint":            server.Endpoint,
		"User":                server.User,
		"Password":            server.Password,
		"Sni":                 server.Sni,
		"UseRemoteRouteRules": server.UseRemoteRouteRules,
		"SkipVerifySSL":       server.SkipVerifySSL,
		"LocalRouteRules":     server.LocalRouteRules,
		"ProxyDomains":        server.ProxyDomains,
	})

	if ret.Error != nil {
		return ret.Error
	}

	return nil
}
