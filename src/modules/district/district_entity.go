package district

import "sim-puskesmas/src/libs/db/pg"

type DistrictModel struct {
	*pg.Model
	Name *string `json:"name"`
}

func (DistrictModel) TableName() string {
	return "districts"
}

func (*Module) autoMigrate() {
	if err := pg.GetDB().AutoMigrate((&DistrictModel{})); err != nil {
		logger.Panic(err)
	}
}

type districtDB struct {
	*pg.Service[DistrictModel]
}

var districtRepo *districtDB

func DistrictRepository() *districtDB {
	if districtRepo == nil {
		districtRepo = &districtDB{
			Service: pg.NewService(
				&pg.Service[DistrictModel]{
					DB: pg.GetDB(),
				},
			),
		}
	}
	return districtRepo
}
