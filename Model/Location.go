package Model

type CountryNameMapping struct {
	Country string			`db:"country"`
	CountryName string		`db:"countryName"`
}