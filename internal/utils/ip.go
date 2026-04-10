package utils

import (
	"net/netip"

	"github.com/oschwald/geoip2-golang/v2"
)

func SearchIP(ip string, file string) (*geoip2.City, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, err
	}
	db, err := geoip2.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(db *geoip2.Reader) {
		_ = db.Close()
	}(db)
	record, err := db.City(addr)
	if err != nil {
		return nil, err
	}
	return record, nil
}
