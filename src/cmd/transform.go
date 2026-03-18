package cmd

import (
	// "strconv"
	// "strings"

	"github.com/MaminirinaEdwino/etl/src/model"
)

func TransformRow(row []string, e *model.Extractor, fieldName map[int]string) (map[string]string, error) {
	var transformData = make(map[string]string)
	// latStr := e.GetValue(row, "Latitude")
	// lngStr := e.GetValue(row, "Longitude")
	// sevStr := e.GetValue(row, "Accident_severity")
	// dateStr := e.GetValue(row, "Accident Date")
	// timeStr := e.GetValue(row, "Time")
	// vTypeStr := e.GetValue(row, "vehicle_type")
	// speedStr := e.GetValue(row, "speed_limit")
	// casualStr := e.GetValue(row, "Number_of_casualities")
	// day_of_week := e.GetValue(row, "Day_of_Week")
	// nbr_vehicle := e.GetValue(row, "Number_of_Vehicles")
	// //conversion
	// lat, _ := strconv.ParseFloat(latStr, 64)
	// lng, _ := strconv.ParseFloat(lngStr, 64)
	// speed, _ := strconv.Atoi(speedStr)
	// casual, _ := strconv.Atoi(casualStr)
	// vehicle, _ := strconv.Atoi(nbr_vehicle)
	// severity := strings.ToUpper(strings.TrimSpace(sevStr))

	// 4. Retourner l'objet structuré
	// return model.RawAccident{
	// 	Index:       e.GetValue(row, "Accident_Index"),
	// 	Date:        dateStr,
	// 	Time:        timeStr,
	// 	Severity:    severity,
	// 	Latitude:    lat,
	// 	Longitude:   lng,
	// 	SpeedLimit:  speed,
	// 	Casualties:  casual,
	// 	Weather:     e.GetValue(row, "weather_conditions"),
	// 	VehicleType: vTypeStr,
	// 	DayOfWeek:   day_of_week,
	// 	Vehicles:    vehicle,
	// }, nil

	for _, el := range fieldName{
		transformData[el] = e.GetValue(row, el)
	}
	return transformData, nil

}