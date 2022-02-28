package dto


type HealthCheckDTO struct {
	Service			string		`json:"service"`
	Interval 		string		`json:"interval"`
	Timeout			string		`json:"timeout"`
}