package geo

import (
	"errors"
	"net"
	"os"
	"path"

	"github.com/oschwald/geoip2-golang"
)

var (
	db           *geoip2.Reader
	EmptyDBError = errors.New("empty db")
)

func init() {
	var (
		geoPath     = os.Getenv("GEO_PATH")
		geoFileName = os.Getenv("GEO_FILE_NAME")
		err         error
	)

	if geoPath == "" {
		return
	}

	if geoFileName == "" {
		geoFileName = "GeoLite2-City.mmdb"
	}

	db, err = geoip2.Open(path.Join(geoPath, geoFileName))
	if err != nil {
		panic(err)
	}
}

func Close() error {
	if db == nil {
		return nil
	}

	return db.Close()
}

func GetOriginCity(ip string) (city *geoip2.City, err error) {
	if db == nil {
		err = EmptyDBError
		return
	}

	return db.City(net.ParseIP(ip))
}

func GetLocation(ip string) (lat, lon float64, err error) {
	city, err := GetOriginCity(ip)
	if err != nil {
		return
	}

	lat = city.Location.Latitude
	lon = city.Location.Longitude
	return
}

// GetCityName 获取城市名称（默认返回英语名称)
func GetCityName(ip string, opts ...Option) (string, error) {
	city, err := GetOriginCity(ip)
	if err != nil {
		return "", err
	}

	o := getOptions(opts...)
	return city.City.Names[o.languageCode], nil
}

// GetCountryName 获取国家名称（默认返回英语名称)
func GetCountryName(ip string, opts ...Option) (string, error) {
	city, err := GetOriginCity(ip)
	if err != nil {
		return "", err
	}

	o := getOptions(opts...)
	return city.Country.Names[o.languageCode], nil
}

// GetProvinceName 获取省份名称（默认返回英语名称)
func GetProvinceName(ip string, opts ...Option) (string, error) {
	city, err := GetOriginCity(ip)
	if err != nil {
		return "", err
	}

	if len(city.Subdivisions) <= 0 {
		return "", nil
	}

	o := getOptions(opts...)
	return city.Subdivisions[0].Names[o.languageCode], nil
}

func GetCountryISOCode(ip string) (string, error) {
	city, err := GetOriginCity(ip)
	if err != nil {
		return "", err
	}

	return city.Country.IsoCode, nil
}

// GetContinent 获取所属洲，比如亚洲
func GetContinent(ip string, opts ...Option) (string, error) {
	city, err := GetOriginCity(ip)
	if err != nil {
		return "", err
	}

	o := getOptions(opts...)
	return city.Continent.Names[o.languageCode], nil
}

type City struct {
	Continent       string
	CountryISOCode  string
	CountryName     string
	ProvinceName    string
	ProvinceIsoCode string
	CityName        string
	Latitude        float64
	Longitude       float64
	PostalCode      string
}

func GetCity(ip string, opts ...Option) (*City, error) {
	city, err := GetOriginCity(ip)
	if err != nil {
		return nil, err
	}

	o := getOptions(opts...)
	return &City{
		Continent:      city.Continent.Names[o.languageCode],
		CountryISOCode: city.Country.IsoCode,
		CountryName:    city.Country.Names[o.languageCode],
		ProvinceName: func() string {
			if len(city.Subdivisions) <= 0 {
				return ""
			}
			return city.Subdivisions[0].Names[o.languageCode]
		}(),
		ProvinceIsoCode: func() string {
			if len(city.Subdivisions) <= 0 {
				return ""
			}
			return city.Subdivisions[0].IsoCode
		}(),
		CityName:   city.City.Names[o.languageCode],
		Latitude:   city.Location.Latitude,
		Longitude:  city.Location.Longitude,
		PostalCode: city.Postal.Code,
	}, nil
}

//func GetIp() {
//
//	// If you are using strings that may be invalid, check that ip is not nil
//	//ip := net.ParseIP("81.2.69.142")
//	ip := net.ParseIP("115.192.211.101")
//	record, err := db.City(ip)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
//	if len(record.Subdivisions) > 0 {
//		fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
//	}
//	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
//	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
//	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
//	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
//	// Output:
//	// Portuguese (BR) city name: Londres
//	// English subdivision name: England
//	// Russian country name: Великобритания
//	// ISO country code: GB
//	// Time zone: Europe/London
//	// Coordinates: 51.5142, -0.0931
//
//	fmt.Println("中文结果")
//	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["zh-CN"])
//	if len(record.Subdivisions) > 0 {
//		fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["zh-CN"])
//	}
//	fmt.Printf("Represented country name:%v\n", record.RepresentedCountry.Names["zh-CN"])
//	fmt.Printf("Continent:%v\n", record.Continent.Names["zh-CN"])
//	fmt.Printf("Registered country:%v\n", record.RepresentedCountry.Names["zh-CN"])
//	fmt.Printf("Russian country name: %v\n", record.Country.Names["zh-CN"])
//	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
//	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
//	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
//}
