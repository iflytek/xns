package api

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
	"strings"
)

type Region struct {
	Code        int    `json:"code" desc:"大区的唯一代码" minimum:"1000" maximum:"9999"`
	Name        string `json:"name" desc:"大区名称"`
	Description string `json:"description"`
	IdcAffinity string `json:"idc_affinity" desc:"大区机房亲和性，值为机房的id，多个用, 隔开，优先级依次降低"`
}

type regionService struct {
	RegionDao   dao.Region
	ProvinceDao dao.Province
	CityDao     dao.City
	IdcDao      dao.Idc
	CountryDao  dao.Country
	RouteDao    dao.Route
}

func (r *regionService) Create(req *Region) (region *models.Region, code int, err error) {
	region, err = r.RegionDao.GetByCode(code)
	if err == nil {
		code = CodeConflict
		err = fmt.Errorf("region code  '%d' already exists", req.Code)
		return
	}

	afs := parseIdcAffinity(req.IdcAffinity)
	req.IdcAffinity, code, err = checkAffinity(r.IdcDao, afs)
	if err != nil {
		return
	}

	region = &models.Region{
		Base:        models.Base{Description: req.Description},
		Name:        req.Name,
		Code:        req.Code,
		IdcAffinity: req.IdcAffinity,
	}

	err = r.RegionDao.Create(region)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (r *regionService) Delete(regionCode int) (code int, err error) {
	provs, err := r.ProvinceDao.GeProvinceByRegionCode(regionCode)
	if err != nil {
		return CodeDbError, err
	}
	if len(provs) > 0 {
		return CodeDbError, fmt.Errorf("You should modify the region of provinces belongs to this region before delete this region. ")
	}
	err = r.RegionDao.Delete(regionCode)
	if err != nil {
		code, err = convertErrorf(err, "%w", err)
	}
	return
}

type UpdateRegion struct {
	Code        int    `json:"code" from:"path"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IdcAffinity string `json:"idc_affinity" desc:"机房亲和性，机房id，多个用逗号分隔，优先级依次降低"`
}

func checkAffinity(idc dao.Idc, afs []string) (afsss string, code int, err error) {
	afss := make([]string, 0, len(afs))
	for _, af := range afs {
		var i *models.Idc
		i, err = idc.GetByIdOrName(af)
		if err != nil {
			code, err = convertErrorf(err, "idc affinity '%s' not found", af)
			return
		}
		afss = append(afss, i.Id)
	}
	return strings.Join(afss, ","), 0, nil
}

func (r *regionService) Update(regionCode int, req *UpdateRegion) (region *models.Region, code int, err error) {
	region, err = r.RegionDao.GetByCode(regionCode)
	if err != nil {
		code, err = convertErrorf(err, "get region '%d error:%w", code, err)
		return
	}
	region.Description = or(req.Description, region.Description)
	region.Name = or(req.Name, region.Name)

	afs := parseIdcAffinity(req.IdcAffinity)
	req.IdcAffinity, code, err = checkAffinity(r.IdcDao, afs)
	if err != nil {
		return
	}
	region.IdcAffinity = req.IdcAffinity

	err = r.RegionDao.Update(region.Id, region)
	if err != nil {
		code = CodeDbError
	}
	return
}

type regionWrap struct {
	*models.Region
	Idcs []*models.Idc `json:"idcs"`
}

func (r *regionService) GetRegions() (ress []*regionWrap, code int, err error) {
	var res []*models.Region
	res, err = r.RegionDao.GetList()
	if err != nil {
		code = CodeDbError
		return
	}
	resWraps := make([]*regionWrap, 0, len(res))
	m := r.IdcMap()
	for _, re := range res {
		var idcs []*models.Idc
		idcs, code, err = r.idcs(re.IdcAffinity, m)
		if err != nil {
			return
		}
		resWraps = append(resWraps, &regionWrap{
			Region: re,
			Idcs:   idcs,
		})
	}
	return resWraps, 0, nil
}

func (r *regionService) GetRegion(regionCode int) (ress *regionWrap, code int, err error) {
	var res *models.Region
	res, err = r.RegionDao.GetByCode(regionCode)
	if err != nil {
		code, err = convertErrorf(err, "get region '%d' err :%w", code, err)
		return
	}

	idcs, code, err := r.idcs(res.IdcAffinity, r.IdcMap())
	if err != nil {
		return nil, code, err
	}
	return &regionWrap{Region: res, Idcs: idcs}, 0, nil
}

type provinceWrap struct {
	*models.Province
	Idcs   []*models.Idc  `json:"idcs"`
	Region *models.Region `json:"region,omitempty"`
}

func (r *regionService) regionMap() map[int]*models.Region {
	res := make(map[int]*models.Region)
	regions, _ := r.RegionDao.GetList()
	for _, region := range regions {
		res[region.Code] = region
	}
	return res
}

func (r *regionService) provinceMap() map[int]*models.Province {
	res := make(map[int]*models.Province)
	regions, _ := r.ProvinceDao.GetList()
	for _, region := range regions {
		res[region.Code] = region
	}
	return res
}

func (r *regionService) GetRegionProvince(regionCode int) (ress []*provinceWrap, code int, err error) {
	var res []*models.Province
	res, err = r.ProvinceDao.GeProvinceByRegionCode(regionCode)
	if err != nil {
		code = CodeDbError
		return
	}
	rgMap := r.regionMap()
	provWraps := make([]*provinceWrap, 0, len(res))
	m := r.IdcMap()
	for _, re := range res {
		var idcs []*models.Idc
		idcs, code, err = r.idcs(re.IdcAffinity, m)
		if err != nil {
			return
		}
		provWraps = append(provWraps, &provinceWrap{
			Province: re,
			Idcs:     idcs,
			Region:   rgMap[re.RegionCode],
		})
	}
	return provWraps, 0, nil
}

type cityWrap struct {
	*models.City
	Idcs     []*models.Idc    `json:"idcs"`
	Province *models.Province `json:"province,omitempty"`
}

func (r *regionService) GetProvinceCity(provCode int) (ress []*cityWrap, code int, err error) {
	var res []*models.City
	res, err = r.CityDao.GetProvinceCities(provCode)
	if err != nil {
		code = CodeDbError
		return
	}

	m := r.IdcMap()
	pm := r.provinceMap()
	provWraps := make([]*cityWrap, 0, len(res))
	for _, re := range res {
		var idcs []*models.Idc
		idcs, code, err = r.idcs(re.IdcAffinity, m)
		if err != nil {
			return
		}
		provWraps = append(provWraps, &cityWrap{
			City:     re,
			Idcs:     idcs,
			Province: pm[re.ProvinceCode],
		})
	}
	return provWraps, 0, nil
}

func (r *regionService) Idcs(idcf string) (res []*models.Idc, code int, err error) {
	for _, idc := range strings.Split(idcf, ",") {
		if idc == "" {
			continue
		}
		idc, err := r.IdcDao.GetByIdOrName(idc)
		if err != nil {
			return nil, CodeRequestError, err
		}
		res = append(res, idc)
	}
	return res, CodeRequestError, nil
}

func (r *regionService) IdcMap() (res map[string]*models.Idc) {
	res = map[string]*models.Idc{}
	idcs, err := r.IdcDao.GetList()
	if err != nil {
		return
	}

	for _, idc := range idcs {
		res[idc.Id] = idc
		res[idc.Name] = idc
	}
	return res
}

func (r *regionService) idcs(idcf string, m map[string]*models.Idc) (res []*models.Idc, code int, err error) {
	for _, idc := range strings.Split(idcf, ",") {
		if idc == "" {
			continue
		}
		idc, ok := m[idc]
		if !ok {
			return nil, CodeRequestError, fmt.Errorf("idc %s not found", idc)
		}
		res = append(res, idc)
	}
	return res, 0, nil
}

func (r *regionService) GetProvinces() (ress []*provinceWrap, code int, err error) {
	var res []*models.Province
	res, err = r.ProvinceDao.GetList()
	if err != nil {
		code = CodeDbError
		return
	}
	m := r.IdcMap()
	rm := r.regionMap()
	provWraps := make([]*provinceWrap, 0, len(res))
	for _, re := range res {
		var idcs []*models.Idc
		idcs, code, err = r.idcs(re.IdcAffinity, m)
		if err != nil {
			return
		}
		provWraps = append(provWraps, &provinceWrap{
			Province: re,
			Idcs:     idcs,
			Region:   rm[re.RegionCode],
		})
	}
	return provWraps, 0, nil
}

func (r *regionService) wrapRegion(region *models.Region) *regionWrap {
	res, _, _ := r.idcs(region.IdcAffinity, r.IdcMap())
	return &regionWrap{Region: region, Idcs: res}
}

func (r *regionService) wrapCity(city *models.City) *cityWrap {
	res, _, _ := r.idcs(city.IdcAffinity, r.IdcMap())
	return &cityWrap{City: city, Idcs: res}
}

func (r *regionService) wrapProvince(prov *models.Province) *provinceWrap {
	res, _, _ := r.idcs(prov.IdcAffinity, r.IdcMap())
	return &provinceWrap{Province: prov, Idcs: res}
}

func (r *regionService) GetProvince(provCode int) (ress *provinceWrap, code int, err error) {
	var res *models.Province
	res, err = r.ProvinceDao.GetByCode(provCode)
	if err != nil {
		code = CodeDbError
		return
	}

	idcs, code, err := r.idcs(res.IdcAffinity, r.IdcMap())
	if err != nil {
		return nil, code, err
	}

	return &provinceWrap{Province: res, Idcs: idcs}, 0, nil
}

func (r *regionService) AddProvince(req *AddProvinceReq) (res *models.Province, code int, err error) {

	if req.CountryCode == 0 {
		req.CountryCode = 101 // china
	}
	res = &models.Province{
		Base:        models.Base{Description: req.Description},
		Name:        req.Name,
		Code:        req.Code,
		RegionCode:  req.RegionCode,
		CountryCode: req.CountryCode,
		IdcAffinity: req.IdcAffinity,
	}

	idcf := parseIdcAffinity(req.IdcAffinity)
	req.IdcAffinity, code, err = checkAffinity(r.IdcDao, idcf)
	if err != nil {
		return
	}
	err = r.ProvinceDao.Create(res)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (r *regionService) checkRuleReferenceProvince(province int) (code int ,err error){
	routes ,err := r.RouteDao.QueryRoutesByRuleCond(fmt.Sprintf("rules like '%%province=%d%%' and rules like '%%region=%%'",province))
	if err != nil{
		return CodeDbError,err
	}
	if len(routes) > 0{
		return CodeRequestError,fmt.Errorf("province and its region is referenced in some routes ,cannot modify region of it. ")
	}
	return 0,nil
}

func (r *regionService) UpdateProvince(pcode int, req map[string]interface{}) (res *models.Province, code int, err error) {
	res, err = r.ProvinceDao.GetByCode(pcode)
	if err != nil {
		code, err = convertErrorf(err, "update prov get prov '%d' error :%w", pcode, err)
		return
	}
	regionCode ,_ := req["region_code"].(float64)
	if int(regionCode) != res.RegionCode{ // region code 改变了，需要改变所有规则中的region code，由于加载顺序关系，会出现部分请求失败,需要禁止该操作
		code ,err = r.checkRuleReferenceProvince(pcode)
		if err != nil{
			return nil, 0, err
		}
	}
	delete(req, "code")
	err = patch(res, req)
	if err != nil {
		code = CodeRequestError
		return
	}

	idcf := parseIdcAffinity(res.IdcAffinity)
	res.IdcAffinity, code, err = checkAffinity(r.IdcDao, idcf)
	if err != nil {
		return
	}

	err = r.ProvinceDao.Update(pcode, res)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (r *regionService) GetCities() (ress []*cityWrap, code int, err error) {
	var res []*models.City
	res, err = r.CityDao.GetList()
	if err != nil {
		code = CodeDbError
		return
	}
	m := r.IdcMap()
	cityWraps := make([]*cityWrap, 0, len(res))
	for _, re := range res {
		var idcs []*models.Idc
		idcs, code, err = r.idcs(re.IdcAffinity, m)
		if err != nil {
			return
		}
		cityWraps = append(cityWraps, &cityWrap{
			City: re,
			Idcs: idcs,
		})
	}
	return cityWraps, 0, nil
}

func (r *regionService) GetCity(cityCode int) (ress *cityWrap, code int, err error) {
	var res *models.City
	res, err = r.CityDao.GetByCode(cityCode)
	if err != nil {
		code = CodeDbError
		return
	}
	m := r.IdcMap()
	idcs, code, err := r.idcs(res.IdcAffinity, m)
	if err != nil {
		return nil, code, err
	}
	return &cityWrap{City: res, Idcs: idcs}, 0, nil
}

func (r *regionService) AddCity(req *AddCityReq) (res *models.City, code int, err error) {
	m := &models.City{
		Base: models.Base{
			Description: req.Description,
		},
		Name:         req.Name,
		Code:         req.Code,
		ProvinceCode: req.ProvinceCode,
		IdcAffinity:  req.IdcAffinity,
	}
	afs := parseIdcAffinity(req.IdcAffinity)
	req.IdcAffinity, code, err = checkAffinity(r.IdcDao, afs)
	if err != nil {
		return
	}
	res = m
	err = r.CityDao.Create(m)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (r *regionService) UpdateCity(ccode int, req map[string]interface{}) (city *models.City, code int, err error) {
	city, err = r.CityDao.GetByCode(ccode)
	if err != nil {
		code, err = convertErrorf(err, "update city ,get city '%d' error:%w", ccode, err)
		return
	}
	delete(req, "code")
	delete(req, "province_code")
	err = patch(city, req)
	if err != nil {
		code = CodeRequestError
		return
	}

	afs := parseIdcAffinity(city.IdcAffinity)
	var afss string
	afss, code, err = checkAffinity(r.IdcDao, afs)
	if err != nil {
		return
	}
	city.IdcAffinity = afss
	err = r.CityDao.Update(ccode, city)
	if err != nil {
		code = CodeDbError
	}
	return
}

func parseIdcAffinity(f string) []string {
	if f == "" {
		return nil
	}
	return strings.Split(f, ",")
}

func (r *regionService) AddCountry(ccode int, name, desc string) (res *models.Country, code int, err error) {
	res = &models.Country{
		Base: models.Base{
			Description: desc,
		},
		Code: ccode,
		Name: name,
	}
	err = r.CountryDao.Create(res)
	if err != nil {
		code = CodeDbError
		return
	}
	return
}

func (r *regionService) DeleteCountry(ccode int, ) (code int, err error) {

	err = r.CountryDao.Delete(ccode)
	if err != nil {
		code, err = convertErrorf(err, "%w", err)
		return
	}
	return
}

//
func (r *regionService) GetCountryList() (res []*models.Country, code int, err error) {
	res, err = r.CountryDao.GetList()
	if err != nil {
		code = CodeDbError
		return
	}
	return
}

//
//todo add country to region
