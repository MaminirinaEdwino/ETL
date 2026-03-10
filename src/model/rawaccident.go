package model

type RawAccident struct {
	Index           string
	Date            string
	DayOfWeek       string
	JunctionControl string
	JunctionDetails string
	Severity        string
	Latitude        float64
	LightConditions string
	LocalAuthority  string
	Hazards         string
	Longitude       float64
	Casualties      int
	Vehicles        int
	PoliceForce     string
	RoadSurface     string
	RoadType        string
	SpeedLimit      int
	Time            string
	AreaType        string // Urban or Rural
	Weather         string
	VehicleType     string
}