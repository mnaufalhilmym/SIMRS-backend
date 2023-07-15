package patient

import (
	"errors"
	accountrole "simrs/src/common/account_role"
	"simrs/src/helpers"
	"simrs/src/libs/db/pg"
	"simrs/src/libs/parser"
	"simrs/src/middlewares/authguard"
	"simrs/src/modules/account"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/patients", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.getPatientList)
	m.App.Get("/api/v1/patient/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.getPatientDetail)
	m.App.Get("/api/v1/patients-count", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.getPatientCount)
	m.App.Post("/api/v1/patient", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.addPatient)
	m.App.Patch("/api/v1/patient/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.updatePatient)
	m.App.Delete("/api/v1/patient/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.deletePatient)
}

func (m *Module) getPatientList(c *fiber.Ctx) error {
	query := new(getPatientListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	patientListData, total, err := m.getPatientListService(&paginationOption{
		limit:  query.Limit,
		lastID: query.LastID,
	}, &searchOption{
		byFamilyCardNumber:     query.SearchByFamilyCardNumber,
		byRelationshipInFamily: query.SearchByRelationshipInFamily,
		byDistrictID:           query.SearchByDistrictID,
		byAny:                  query.Search,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Pagination: &pagination{
			Total: &total,
		},
		Data: patientListData,
	})
}

func (m *Module) getPatientDetail(c *fiber.Ctx) error {
	param := new(getPatientDetailReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	patientDetailData, err := m.getPatientDetailService(param.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: patientDetailData,
	})
}

func (m *Module) getPatientCount(c *fiber.Ctx) error {
	query := new(getPatientCountReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	count, err := m.getPatientCountService(&searchOption{
		byFamilyCardNumber:     query.SearchByFamilyCardNumber,
		byRelationshipInFamily: query.SearchByRelationshipInFamily,
		byDistrictID:           query.SearchByDistrictID,
		byAny:                  query.Search,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: count,
	})
}

func (m *Module) addPatient(c *fiber.Ctx) error {
	req := new(addPatientReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	dateOfBirth := time.Date(req.DateOfBirth.Year(), req.DateOfBirth.Month(), req.DateOfBirth.Day(), 0, 0, 0, 0, req.DateOfBirth.Location())
	patientDetailData, err := m.addPatientService(&PatientModel{
		MedicalRecordNumber:            req.MedicalRecordNumber,
		FamilyCardNumber:               req.FamilyCardNumber,
		RelationshipInFamily:           req.RelationshipInFamily,
		PopulationIdentificationNumber: req.PopulationIdentificationNumber,
		Name:                           req.Name,
		Gender:                         req.Gender,
		PlaceOfBirth:                   req.PlaceOfBirth,
		DateOfBirth:                    &dateOfBirth,
		Address:                        req.Address,
		DistrictID:                     req.DistrictID,
		Job:                            req.Job,
		Religion:                       req.Religion,
		BloodGroup:                     req.BloodGroup,
		Insurence:                      req.Insurence,
		InsurenceNumber:                req.InsurenceNumber,
		Phone:                          req.Phone,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: patientDetailData,
	})
}

func (m *Module) updatePatient(c *fiber.Ctx) error {
	param := new(updatePatientReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	req := new(updatePatientReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	dateOfBirth := time.Date(req.DateOfBirth.Year(), req.DateOfBirth.Month(), req.DateOfBirth.Day(), 0, 0, 0, 0, req.DateOfBirth.Location())
	patientDetailData := &PatientModel{
		Model: &pg.Model{
			ID: param.ID,
		},
		MedicalRecordNumber:            req.MedicalRecordNumber,
		FamilyCardNumber:               req.FamilyCardNumber,
		RelationshipInFamily:           req.RelationshipInFamily,
		PopulationIdentificationNumber: req.PopulationIdentificationNumber,
		Name:                           req.Name,
		Gender:                         req.Gender,
		PlaceOfBirth:                   req.PlaceOfBirth,
		DateOfBirth:                    &dateOfBirth,
		Address:                        req.Address,
		DistrictID:                     req.DistrictID,
		Job:                            req.Job,
		Religion:                       req.Religion,
		BloodGroup:                     req.BloodGroup,
		Insurence:                      req.Insurence,
		InsurenceNumber:                req.InsurenceNumber,
		Phone:                          req.Phone,
	}

	if req.DateOfBirth != nil {
		dateOfBirth := time.Date(req.DateOfBirth.Year(), req.DateOfBirth.Month(), req.DateOfBirth.Day(), 0, 0, 0, 0, req.DateOfBirth.Location())
		patientDetailData.DateOfBirth = &dateOfBirth
	}

	patientDetailData, err := m.updatePatientService(patientDetailData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: patientDetailData,
	})
}

func (m *Module) deletePatient(c *fiber.Ctx) error {
	param := new(deletePatientReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	if err := m.deletePatientService(param.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: param.ID,
	})
}
