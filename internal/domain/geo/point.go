package geo

import (
	"errors"
)

type Point struct {
	Lat float64
	Lon float64
}

func (d *Point) Validate() error {

	if !(d.Lat >= -90 && d.Lat <= 90) {
		return errors.New("invalid latitude")
	}

	if !(d.Lon >= -180 && d.Lon <= 180) {
		return errors.New("invalid longitude")
	}

	return nil
}
