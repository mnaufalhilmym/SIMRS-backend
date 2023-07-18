package setting

import "simrs/src/libs/db/pg"

func (m *Module) getSettingService() (*SettingModel, error) {
	return SettingRepository().FindOne(&pg.FindOneOption{})
}

func (m *Module) addSettingService(data *SettingModel) (*SettingModel, error) {
	return SettingRepository().Create(data)
}

func (m *Module) updateSettingService(data *SettingModel) (*SettingModel, error) {
	settingData, err := m.getSettingService()
	if err != nil {
		return nil, err
	}
	data.ID = settingData.ID
	return SettingRepository().Update(data)
}
