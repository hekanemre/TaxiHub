package domain

// type Location struct {
// 	Lat float64 `bson:"lat" json:"lat"`
// 	Lon float64 `bson:"lon" json:"lon"`
// }

type Location struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}
