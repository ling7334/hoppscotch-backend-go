package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

/* Type of a request
REST: restful request

GQL: graphql request
*/
// CREATE TYPE "ReqType" AS ENUM ('REST','GQL');
type ReqType string

const (
	REST ReqType = "REST"
	GQL  ReqType = "GQL"
)

func (rt *ReqType) Scan(value interface{}) error {
	*rt = ReqType(value.(string))
	return nil
}

func (rt ReqType) Value() (driver.Value, error) {
	return string(rt), nil
}

/* Role of a member in team
OWNER: team owner

VIEWER: can only view team assest

EDITOR: can read and write team assest
*/
// CREATE TYPE "TeamAccessRole" AS ENUM ('OWNER','VIEWER','EDITOR');
type TeamAccessRole string

const (
	OWNER  TeamAccessRole = "OWNER"
	VIEWER TeamAccessRole = "VIEWER"
	EDITOR TeamAccessRole = "EDITOR"
)

func (tmr *TeamAccessRole) Scan(value interface{}) error {
	*tmr = TeamAccessRole(value.(string))
	return nil
}

func (tmr TeamAccessRole) Value() (driver.Value, error) {
	return string(tmr), nil
}

type Auth struct {
	AuthType   string `json:"authType"`
	AuthActive bool   `json:"authActive"`
}

type Body struct {
	Body        *string `json:"body"`
	ContentType *string `json:"contentType"`
}

type Param struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Active bool   `json:"active"`
}

type Header struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Active bool   `json:"active"`
}

type ReqDetail struct {
	V                json.Number `json:"v,string"`
	Auth             Auth        `json:"auth"`
	Body             Body        `json:"body"`
	Name             string      `json:"name"`
	Method           string      `json:"method"`
	Params           []Param     `json:"params"`
	Headers          []Header    `json:"headers"`
	Endpoint         string      `json:"endpoint"`
	TestScript       string      `json:"testScript"`
	PreRequestScript string      `json:"preRequestScript"`
}

// Value Marshal
func (a ReqDetail) Value() ([]byte, error) {
	return json.Marshal(a)
	// return string(a)
}

// Scan Unmarshal
func (a *ReqDetail) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	// *a = ReqDetail(b)
	// return nil
	return json.Unmarshal(b, &a)
}

// Ownable table has field which link to user
type Ownable interface {
	IsOwner(uid string) bool
}

// Orderable table has fields indicate its order index and collection it's in
type Orderable interface {
	ParentColName() string
	Move(*gorm.DB, *string) error
}

// TeamResource table has field which link to user
type TeamResource interface {
	GetTeamID() string
	Can(*gorm.DB, string, TeamAccessRole) bool
}
