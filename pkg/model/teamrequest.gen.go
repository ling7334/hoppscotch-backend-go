// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
	"gorm.io/gorm"
)

const TableNameTeamRequest = "TeamRequest"

// TeamRequest mapped from table <TeamRequest>
type TeamRequest struct {
	ID           string     	`gorm:"column:id;type:text;primaryKey" json:"id"`
	CollectionID string     	`gorm:"column:collectionID;type:text;not null" json:"collectionID"`
	TeamID       string     	`gorm:"column:teamID;type:text;not null" json:"teamID"`
	Title        string     	`gorm:"column:title;type:text;not null" json:"title"`
	Request      string      	`gorm:"column:request;type:jsonb;not null" json:"request"`
	OrderIndex   int32      	`gorm:"column:orderIndex;type:integer;not null" json:"orderIndex"`
	CreatedOn    time.Time  	`gorm:"column:createdOn;type:timestamp(3) without time zone;not null;default:CURRENT_TIMESTAMP" json:"createdOn"`
	UpdatedOn    time.Time  	`gorm:"column:updatedOn;type:timestamp(3) without time zone;not null;autoUpdateTime" json:"updatedOn"`
	Team         Team      		`gorm:"foreignKey:TeamID" json:"team"`
	Collection   TeamCollection `gorm:"foreignKey:CollectionID" json:"collection"`
}

// TableName TeamRequest's table name
func (*TeamRequest) TableName() string {
	return TableNameTeamRequest
}

func (r *TeamRequest) GetTeamID() string {
	return r.TeamID
}

func (r *TeamRequest) Can(db *gorm.DB, uid string, role TeamMemberRole) bool {
	member := &TeamMember{}
	if db.First(member, `"userUid"=? AND "teamID"=?`, uid, r.TeamID).Error != nil {
		switch role {
		case VIEWER:
			return true
		case EDITOR:
			return member.Role == OWNER || member.Role == EDITOR
		case OWNER:
			return member.Role == OWNER
		default:
			return false
		}
	}
	return false
}
func (*TeamRequest) ParentColName() string {
	return "CollectionID"
}

func (r *TeamRequest) Move(db *gorm.DB, next *string) error {
	panic("implement me")
}