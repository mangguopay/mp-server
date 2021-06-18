package test

import (
	"a.a/cu/strext"
	"fmt"
	"github.com/kellydunn/golang-geo"
	"math"
	"testing"
)

//转化为弧度(rad)
func rad(d float64) (r float64) {
	r = d * math.Pi / 180.0
	return
}

func latitude_longitude_distance(lon1, lat1, lon2, lat2 float64) (distance float64) {
	//赤道半径(单位m)
	const EARTH_RADIUS = 6378137
	rad_lat1 := rad(lat1)
	rad_lon1 := rad(lon1)
	rad_lat2 := rad(lat2)
	rad_lon2 := rad(lon2)
	if rad_lat1 < 0 {
		rad_lat1 = math.Pi/2 + math.Abs(rad_lat1)
	}
	if rad_lat1 > 0 {
		rad_lat1 = math.Pi/2 - math.Abs(rad_lat1)
	}
	if rad_lon1 < 0 {
		rad_lon1 = math.Pi*2 - math.Abs(rad_lon1)
	}
	if rad_lat2 < 0 {
		rad_lat2 = math.Pi/2 + math.Abs(rad_lat2)
	}
	if rad_lat2 > 0 {
		rad_lat2 = math.Pi/2 - math.Abs(rad_lat2)
	}
	if rad_lon2 < 0 {
		rad_lon2 = math.Pi*2 - math.Abs(rad_lon2)
	}
	x1 := EARTH_RADIUS * math.Cos(rad_lon1) * math.Sin(rad_lat1)
	y1 := EARTH_RADIUS * math.Sin(rad_lon1) * math.Sin(rad_lat1)
	z1 := EARTH_RADIUS * math.Cos(rad_lat1)

	x2 := EARTH_RADIUS * math.Cos(rad_lon2) * math.Sin(rad_lat2)
	y2 := EARTH_RADIUS * math.Sin(rad_lon2) * math.Sin(rad_lat2)
	z2 := EARTH_RADIUS * math.Cos(rad_lat2)
	d := math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2) + (z1-z2)*(z1-z2))
	theta := math.Acos((EARTH_RADIUS*EARTH_RADIUS + EARTH_RADIUS*EARTH_RADIUS - d*d) / (2 * EARTH_RADIUS * EARTH_RADIUS))
	distance = theta * EARTH_RADIUS
	return
}
func TestLat(t *testing.T) {
	//fmt.Println(strext.ToStringNoPoint(latitude_longitude_distance(30, 30, 40, 40)))
	count()

}

func count() {
	lat2 := "5e-324" // 维度
	lng2 := "5e-324" // 经度
	//lat2 := "23.0795382400"  // 维度
	//lng2 := "113.3406764300" // 经度
	// ======岭南创意园=======
	// ======创投小镇=======
	lat3 := "23.0714778800"  // 维度
	lng3 := "113.3125147400" // 经度
	fmt.Println("--------------->", strext.ToFloat64(lng3))
	// Make a few points
	p := geo.NewPoint(strext.ToFloat64(lat2), strext.ToFloat64(lng2))
	p2 := geo.NewPoint(strext.ToFloat64(lat3), strext.ToFloat64(lng3))

	// find the great circle distance between them
	dist := p.GreatCircleDistance(p2)
	fmt.Printf("great circle distance: %v\n", dist)
}

func TestTime(t *testing.T) {

}
