package pg

import (
	"errors"
	"simrs/src/libs/validator"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service[T any] struct {
	DB    *gorm.DB
	Model T
}

func NewService[T any](service *Service[T]) *Service[T] {
	if service.DB == nil {
		logger.Panic("service.DB must exist")
	}
	return service
}

func (s *Service[T]) Count(countOption *CountOption) (*int64, error) {
	docStruct := s.Model

	countQuery := s.DB.Model(&docStruct)

	if countOption.Where != nil {
		for _, where := range *countOption.Where {
			countQuery = countQuery.Where(where.Query, where.Args...)
		}
	}

	count := new(int64)
	if err := countQuery.Count(count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log(err)
		} else {
			logger.Error(err)
		}
		return nil, err
	}

	return count, nil
}

func (s *Service[T]) FindOne(findOption *FindOneOption) (*T, error) {
	docStruct := s.Model

	selectQuery := s.DB.Model(&docStruct)

	if findOption.IncludeTables != nil {
		for _, table := range *findOption.IncludeTables {
			selectQuery = selectQuery.Preload(table.Query, table.Args...)
		}
	}
	if findOption.Where != nil {
		for _, where := range *findOption.Where {
			selectQuery = selectQuery.Where(where.Query, where.Args...)
		}
	}
	if findOption.Order != nil {
		for _, order := range *findOption.Order {
			selectQuery = selectQuery.Order(order)
		}
	}

	if err := selectQuery.Take(&docStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log(err)
		} else {
			logger.Error(err)
		}
		return nil, err
	}

	return &docStruct, nil
}

func (s *Service[T]) FindAll(findOption *FindAllOption) (*[]*T, *Pagination, error) {
	docStruct := []*T{}

	selectQuery := s.DB.Model(&docStruct)

	if findOption.IncludeTables != nil {
		for _, table := range *findOption.IncludeTables {
			selectQuery = selectQuery.Preload(table.Query, table.Args...)
		}
	}

	if findOption.Where != nil {
		for _, where := range *findOption.Where {
			if where.IncludeInCount {
				selectQuery = selectQuery.Where(where.Query, where.Args...)
			}
		}
	}

	var total int64
	selectQuery.Count(&total)

	if findOption.Where != nil {
		for _, where := range *findOption.Where {
			if !where.IncludeInCount {
				selectQuery = selectQuery.Where(where.Query, where.Args...)
			}
		}
	}

	if findOption.Order != nil {
		for _, order := range *findOption.Order {
			selectQuery = selectQuery.Order(order)
		}
	}

	if findOption.Limit != nil && *findOption.Limit > 0 {
		if *findOption.Limit > FindAllMaximumLimit {
			*findOption.Limit = FindAllMaximumLimit
		}
	} else {
		*findOption.Limit = FindAllDefaultLimit
	}
	selectQuery = selectQuery.Limit(*findOption.Limit)

	if findOption.Offset != nil && *findOption.Offset > 0 {
		selectQuery = selectQuery.Offset(*findOption.Offset)
	}

	if err := selectQuery.Find(&docStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log(err)
		} else {
			logger.Error(err)
		}
		return nil, nil, err
	}

	return &docStruct, &Pagination{
		Count: len(docStruct),
		Limit: *findOption.Limit,
		Total: int(total),
	}, nil
}

func (s *Service[T]) Create(data *T, createOption ...*CreateOption) (*T, error) {
	if err := validator.Struct(data); err != nil {
		logger.Log(err)
		return nil, err
	}

	insertQuery := s.DB.Model(data)

	if len(createOption) > 0 {
		if createOption[0].IsUpsert {
			insertQuery = insertQuery.Clauses(clause.OnConflict{UpdateAll: true})
		}
	}

	if err := insertQuery.Create(data).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *Service[T]) BulkCreate(data *[]*T, createOption ...*CreateOption) (*[]*T, error) {
	for _, doc := range *data {
		if err := validator.Struct(doc); err != nil {
			logger.Log(err)
			return nil, err
		}
	}

	insertQuery := s.DB.Model(data)

	if len(createOption) > 0 {
		if createOption[0].IsUpsert {
			insertQuery = insertQuery.Clauses(clause.OnConflict{UpdateAll: true})
		}
	}

	if err := insertQuery.Create(data).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *Service[T]) Update(data *T, updateOption ...*UpdateOption) (*T, error) {
	if err := validator.Struct(data); err != nil {
		logger.Log(err)
		return nil, err
	}

	updateQuery := s.DB.Model(data)

	if len(updateOption) > 0 && updateOption[0].Where != nil {
		for _, where := range *updateOption[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
	}

	if err := updateQuery.Updates(data).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *Service[T]) BulkUpdate(data *[]*T, updateOption ...*UpdateOption) (*[]*T, error) {
	for _, doc := range *data {
		if err := validator.Struct(doc); err != nil {
			logger.Log(err)
			return nil, err
		}
	}

	updateQuery := s.DB.Model(data)

	if len(updateOption) > 0 && updateOption[0].Where != nil {
		for _, where := range *updateOption[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
	}

	if err := updateQuery.Updates(data).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *Service[T]) Replace(data *T, replaceOption ...*ReplaceOption) error {
	if err := validator.Struct(data); err != nil {
		logger.Log(err)
		return err
	}

	updateQuery := s.DB.Model(data)

	if len(replaceOption) > 0 && replaceOption[0].Where != nil {
		for _, where := range *replaceOption[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
	}

	if err := updateQuery.Updates(data).Error; err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Service[T]) Destroy(data *T, destroyOption ...*DestroyOption) error {
	deleteQuery := s.DB.Model(data)

	if len(destroyOption) > 0 && destroyOption[0].Where != nil {
		for _, where := range *destroyOption[0].Where {
			deleteQuery = deleteQuery.Where(where.Query, where.Args...)
		}
	}

	if err := deleteQuery.Delete(data).Error; err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Service[T]) BulkDestroy(data *[]*T, destroyOption ...*DestroyOption) error {
	deleteQuery := s.DB.Model(data)

	if len(destroyOption) > 0 && destroyOption[0].Where != nil {
		for _, where := range *destroyOption[0].Where {
			deleteQuery = deleteQuery.Where(where.Query, where.Args...)
		}
	}

	if err := deleteQuery.Delete(data).Error; err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
