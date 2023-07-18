package setting

type worker struct {
	Name     *string `json:"name"`
	Position *string `json:"position"`
}

type addSettingReq struct {
	CoverImg *string `json:"coverImg" validate:"required"`
	Workers  *worker `json:"workers" validate:"required"`
	Vision   *string `json:"vision" validate:"required"`
	Mission  *string `json:"mission" validate:"required"`
}

type updateSettingReq struct {
	CoverImg *string `json:"coverImg"`
	Workers  *worker `json:"workers"`
	Vision   *string `json:"vision"`
	Mission  *string `json:"mission"`
}

type response struct {
	Error *string     `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
