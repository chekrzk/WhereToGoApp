package building

import "errors"

var ErrNotFound = errors.New("Building not found")

type Building struct {
	Osm_type      string
	Osm_id        int64
	Building_type string

	Street       string
	House_number string
}

type ResolvedBuilding struct {
	Building  Building
	DistanceM float64
	Match     string
}
