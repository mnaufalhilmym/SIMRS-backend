package setting

import (
	"simrs/src/libs/db/pg"

	"gorm.io/datatypes"
)

type SettingModel struct {
	*pg.Model
	CoverImg *string         `json:"coverImg"`
	Workers  *datatypes.JSON `json:"workers"`
	Vision   *string         `json:"vision"`
	Mission  *string         `json:"mission"`
}

func (SettingModel) TableName() string {
	return "setting"
}

func (*Module) autoMigrate() {
	if err := pg.GetDB().AutoMigrate(&SettingModel{}); err != nil {
		logger.Panic(err)
	}
}

func (m *Module) createInitialSetting() {
	if _, err := m.getSettingService(); err == nil {
		return
	}
	if _, err := m.addSettingService(&SettingModel{}); err != nil {
		logger.Panic(err)
	}
}

type settingDB struct {
	*pg.Service[SettingModel]
}

var settingRepo *settingDB

func SettingRepository() *settingDB {
	if settingRepo == nil {
		settingRepo = &settingDB{
			Service: pg.NewService(
				&pg.Service[SettingModel]{
					DB: pg.GetDB(),
				},
			),
		}
	}
	return settingRepo
}
