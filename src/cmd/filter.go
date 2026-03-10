package cmd

import "github.com/MaminirinaEdwino/etl/src/model"

func ShouldKeep(acc model.RawAccident, cfg model.FilterConfig) bool {
	if acc.Vehicles < cfg.MinVehicles {
		return false
	}
	if cfg.Severity != "" && acc.Severity != cfg.Severity {
		return false
	}
	return true
}