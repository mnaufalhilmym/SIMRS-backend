package patientexamination

import (
	"errors"
	accountrole "simrs/src/common/account_role"
	"simrs/src/helpers"
	"simrs/src/libs/db/pg"
	"simrs/src/libs/parser"
	"simrs/src/middlewares/authguard"
	"simrs/src/modules/account"
	"simrs/src/modules/patient"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/patient-examinations", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.getPatientExaminationList)
	m.App.Get("/api/v1/patient-examination/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.getPatientExaminationDetail)
	m.App.Post("/api/v1/patient-examination", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.addPatientExamination)
	m.App.Patch("/api/v1/patient-examination/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.updatePatientExamination)
	m.App.Delete("/api/v1/patient-examination/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.deletePatientExamination)
}

func (m *Module) getPatientExaminationList(c *fiber.Ctx) error {
	query := new(getPatientExaminationListQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	patientExaminationListData, total, err := m.getPatientExaminationListService(&paginationOption{
		limit:  query.Limit,
		lastID: query.LastID,
	}, query.PatientID, query.Search)
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
		Data: patientExaminationListData,
	})
}

func (m *Module) getPatientExaminationDetail(c *fiber.Ctx) error {
	param := new(getPatientExaminationDetailReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	patientExaminationDetailData, err := m.getPatientExaminationDetailService(param.ID)
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
		Data: patientExaminationDetailData,
	})
}

func (m *Module) addPatientExamination(c *fiber.Ctx) error {
	req := new(addPatientExaminationReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	patientDetailData, err := m.getPatientDetailService(req.PatientID)
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

	if patientDetailData.LastHealthCheckTime == nil || req.ExaminationTime.After(*patientDetailData.LastHealthCheckTime) {
		if _, err := m.updatePatientService(&patient.PatientModel{
			Model: &pg.Model{
				ID: req.PatientID,
			},
			LastHealthCheckTime: req.ExaminationTime,
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
			})
		}
	}

	patientExaminationDetailData, err := m.addPatientExaminationService(&PatientExaminationModel{
		PatientID:       req.PatientID,
		ExaminationTime: req.ExaminationTime,
		Examination:     req.Examination,
		Treatment:       req.Treatment,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: patientExaminationDetailData,
	})
}

func (m *Module) updatePatientExamination(c *fiber.Ctx) error {
	param := new(updatePatientExaminationReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	req := new(updatePatientExaminationReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	if req.ExaminationTime != nil {
		patientExaminationDetailData, err := m.getPatientExaminationDetailService(param.ID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
			})
		}

		patientDetailData, err := m.getPatientDetailService(patientExaminationDetailData.PatientID)
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

		if patientDetailData.LastHealthCheckTime == nil || req.ExaminationTime.After(*patientDetailData.LastHealthCheckTime) {
			if _, err := m.updatePatientService(&patient.PatientModel{
				Model: &pg.Model{
					ID: patientExaminationDetailData.PatientID,
				},
				LastHealthCheckTime: req.ExaminationTime,
			}); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&response{
					Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
				})
			}
		}
	}

	patientExaminationDetailData, err := m.updatePatientExaminationService(&PatientExaminationModel{
		Model: &pg.Model{
			ID: param.ID,
		},
		ExaminationTime: req.ExaminationTime,
		Examination:     req.Examination,
		Treatment:       req.Treatment,
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
		Data: patientExaminationDetailData,
	})
}

func (m *Module) deletePatientExamination(c *fiber.Ctx) error {
	param := new(deletePatientExaminationReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	if err := m.deletePatientExaminationService(param.ID); err != nil {
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
