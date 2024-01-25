package model

import (
	"database/sql/driver"
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
// CREATE TYPE "TeamMemberRole" AS ENUM ('OWNER','VIEWER','EDITOR');
type TeamMemberRole string

const (
	OWNER  TeamMemberRole = "OWNER"
	VIEWER TeamMemberRole = "VIEWER"
	EDITOR TeamMemberRole = "EDITOR"
)

func (tmr *TeamMemberRole) Scan(value interface{}) error {
	*tmr = TeamMemberRole(value.(string))
	return nil
}

func (tmr TeamMemberRole) Value() (driver.Value, error) {
	return string(tmr), nil
}

// JSONB Interface for JSONB Field
type JSONB string

// Value Marshal
func (a JSONB) Value() string {
	// return json.Marshal(a)
	return string(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	*a = JSONB(b)
	// return json.Unmarshal(b, &a)
	return nil
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
	Can(*gorm.DB, string, TeamMemberRole) bool
}
