package utils

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"net/netip"

	"github.com/oschwald/geoip2-golang/v2"
)

func SearchIP(ip string) (*geoip2.City, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, err
	}
	db, err := geoip2.Open(config.Env.GeoCityDB)
	if err != nil {
		log.Logger.Warningf("Failed to open GeoCityDB: %v", err)
		return nil, err
	}
	defer func(db *geoip2.Reader) {
		if err = db.Close(); err != nil {
			log.Logger.Warningf("Failed to close GeoCityDB: %v", err)
		}
	}(db)
	record, err := db.City(addr)
	if err != nil {
		log.Logger.Warningf("Failed to read GeoCityDB: %v", err)
		return nil, err
	}
	return record, nil
}
