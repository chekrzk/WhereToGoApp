package geo

import (
	"errors"
)

var ErrWrongCoordinate = errors.New("invalid lat or lon input")

type Point struct {
	Lat float64
	Lon float64
}

func (d *Point) Validate() error {

	if !(d.Lat >= -90 && d.Lat <= 90) {
		return ErrWrongCoordinate
	}

	if !(d.Lon >= -180 && d.Lon <= 180) {
		return ErrWrongCoordinate
	}

	return nil
}
