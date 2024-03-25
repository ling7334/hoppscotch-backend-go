package dto

import "model"

type TeamRequestExportJSON struct {
	V                string         `json:"v"`
	Auth             model.Auth     `json:"auth"`
	Body             model.Body     `json:"body"`
	Name             string         `json:"name"`
	Method           string         `json:"method"`
	Params           []model.Param  `json:"params"`
	Headers          []model.Header `json:"headers"`
	Endpoint         string         `json:"endpoint"`
	TestScript       string         `json:"testScript"`
	PreRequestScript string         `json:"preRequestScript"`
}

type TeamCollectionExportJSON struct {
	Name     string                     `json:"name"`
	Folders  []TeamCollectionExportJSON `json:"folders"`
	Requests []TeamRequestExportJSON    `json:"requests"`
	Data     *string                    `json:"data"`
}

type TeamCollectionImportDataJSON struct {
	Headers []model.Header `json:"headers"`
	Auth    model.Auth     `json:"auth"`
}
type TeamCollectionImportJSON struct {
	Name     string                       `json:"name"`
	Folders  []TeamCollectionImportJSON   `json:"folders"`
	Requests []TeamRequestExportJSON      `json:"requests"`
	Headers  []model.Header               `json:"headers"`
	Auth     model.Auth                   `json:"auth"`
	Data     TeamCollectionImportDataJSON `json:"data"`
}

type UserRequestExportJSON struct {
	ID               string         `json:"id"`
	V                string         `json:"v"`
	Auth             model.Auth     `json:"auth"`
	Body             model.Body     `json:"body"`
	Name             string         `json:"name"`
	Method           string         `json:"method"`
	Params           []model.Param  `json:"params"`
	Headers          []model.Header `json:"headers"`
	Endpoint         string         `json:"endpoint"`
	TestScript       string         `json:"testScript"`
	PreRequestScript string         `json:"preRequestScript"`
}

type UserCollectionExportJSON struct {
	ID       string                     `json:"id"`
	Name     string                     `json:"name"`
	Folders  []UserCollectionExportJSON `json:"folders"`
	Requests []UserRequestExportJSON    `json:"requests"`
	Data     *string                    `json:"data"`
}
