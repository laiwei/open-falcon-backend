package owl

import (
	"fmt"

	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

type ErrNotInSameHierarchy struct {
	message string
}

func (err ErrNotInSameHierarchy) Error() string {
	return err.message
}

type cityInfo struct {
	ProvinceId   int16  `db:"pv_id"`
	ProvinceName string `db:"pv_name"`
	CityId       int16  `db:"ct_id"`
	CityName     string `db:"ct_name"`
}

type city1view struct {
	Id       int16              `json:"id"`
	Name     string             `json:"name"`
	PostCode string             `json:"post_code"`
	Province *owlModel.Province `json:"province"`
}

// Checks if the hierarchy for province and city are in the same administrative region
func CheckHierarchyForCity(provinceId int16, cityId int16) error {
	if cityId == -1 {
		return nil
	}

	cityInfo := &cityInfo{}
	DbFacade.SqlxDbCtrl.Get(
		cityInfo,
		`
		SELECT pv_id, pv_name, ct_id, ct_name
		FROM owl_province
			INNER JOIN
			owl_city
			ON pv_id = ct_pv_id
				AND ct_id = ?
		`,
		cityId,
	)

	if cityInfo.ProvinceId != provinceId {
		return ErrNotInSameHierarchy{
			message: fmt.Sprintf(
				"City[ID: %d][%s] should be belonging to province[ID: %d]. But got province[ID: %d][%s]",
				cityInfo.CityId, cityInfo.CityName,
				provinceId,
				cityInfo.ProvinceId, cityInfo.ProvinceName,
			),
		}
	}

	return nil
}

func GetProvincesByName(name string) []*owlModel.Province {
	var q = DbFacade.GormDb.Model(&owlModel.Province{}).
		Select(`
		*
		`).
		Where(`
		pv_name LIKE ?
		`,
		name+"%",
	)

	var results []*owlModel.Province
	gormExt.ToDefaultGormDbExt(q.Find(&results))

	return results
}

func GetCitiesByName(name string) []*city1view {
	var q = DbFacade.GormDb.Model(&owlModel.City1{}).
		Select(`
		ct_id, ct_name, ct_post_code, pv_id, pv_name
		`).
		Joins(`
		INNER JOIN
		owl_province
		ON ct_pv_id = pv_id
		AND ct_name LIKE ?
		`,
		name+"%",
	)

	var results []*owlModel.City1
	gormExt.ToDefaultGormDbExt(q.Find(&results))

	var views = []*city1view{}
	for _, r := range results {
		v := &city1view{
			Id:       r.Id,
			Name:     r.Name,
			PostCode: r.PostCode,
			Province: &owlModel.Province{
				Id:   r.ProvinceId,
				Name: r.ProvinceName,
			},
		}
		views = append(views, v)
	}
	return views
}

func GetCitiesInProvinceByName(pvId int, name string) []*owlModel.City2 {
	var q = DbFacade.GormDb.Model(&owlModel.City2{}).
		Select(`
		ct_id, ct_name, ct_post_code
		`).
		Where(`
		ct_pv_id = ?
		AND
		ct_name LIKE ?
		`,
		pvId,
		name+"%",
	)

	var results []*owlModel.City2
	gormExt.ToDefaultGormDbExt(q.Find(&results))

	return results
}
