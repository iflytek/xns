package resource

import (
	"bufio"
	"fmt"
	"os"
	"text/template"
)

type City struct {
	Code         int
	Name         string
	ProvinceCode int
}

type Province struct {
	Code        int
	Name        string
	CountryCode int
	RegionCode  int
}

type Region struct {
	Code int
	Name string
	CountryCode int
}

//var provinces = []Province{
//	{Code: 0, Name: "", CountryCode: 0, RegionCode: 0},
//}

//var Cities = []City{
//	{Code: 1, Name: "", ProvinceCode: ""},
//}

var citiesTlp = `
package resource

var Cities = []City{
	{{ range $c:=. -}}
	{Code: {{$c.Code}},Name: "{{$c.Name}}",ProvinceCode:{{$c.ProvinceCode}}},
	{{ end }}
}

`

var provinceTlp = `
package resource

var Provinces = []Province{
	{{range $p:=. -}}
	{Code: {{$p.Code}}, Name: "{{$p.Name}}", CountryCode: {{$p.CountryCode}}, RegionCode: {{$p.RegionCode}}},
	{{ end}}
}


`

func generateCites() error {
	tlp, err := template.New("cit").Parse(citiesTlp)
	if err != nil {
		panic(err)
	}
	output, err := os.Create("city.go")
	if err != nil {
		panic(err)
	}

	city, err := os.Open("city.txt")

	if err != nil {
		return err
	}

	sc := bufio.NewScanner(city)

	cities := make([]*City, 0, 400)
	for sc.Scan() {
		txt := sc.Text()
		c := &City{}
		_, err = fmt.Sscanf(txt, "%d%s%d", &c.Code, &c.Name, &c.ProvinceCode)
		if err != nil {
			return err
		}
		cities = append(cities, c)
	}

	return tlp.Execute(output, cities)
}

func generateProvince() error {
	tlp, err := template.New("province").Parse(provinceTlp)
	if err != nil {
		panic(err)
	}
	output, err := os.Create("province.go")
	if err != nil {
		panic(err)
	}

	city, err := os.Open("province.txt")

	if err != nil {
		return err
	}

	sc := bufio.NewScanner(city)

	cities := make([]*Province, 0, 400)
	for sc.Scan() {
		txt := sc.Text()
		c := &Province{}
		_, err = fmt.Sscanf(txt, "%d%s%d%d", &c.Code, &c.Name, &c.CountryCode,&c.RegionCode)
		if err != nil {
			return err
		}
		cities = append(cities, c)
	}

	return tlp.Execute(output, cities)
}
