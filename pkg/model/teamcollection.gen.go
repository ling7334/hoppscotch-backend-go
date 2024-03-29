// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
	"gorm.io/gorm"
)

const TableNameTeamCollection = "TeamCollection"

// TeamCollection mapped from table <TeamCollection>
type TeamCollection struct {
	ID         string     			`gorm:"column:id;type:text;primaryKey" json:"id"`
	ParentID   *string    			`gorm:"column:parentID;type:text" json:"parentID"`
	TeamID     string     			`gorm:"column:teamID;type:text;not null" json:"teamID"`
	Title      string     			`gorm:"column:title;type:text;not null" json:"title"`
	OrderIndex int32      			`gorm:"column:orderIndex;type:integer;not null" json:"orderIndex"`
	CreatedOn  time.Time  			`gorm:"column:createdOn;type:timestamp(3) without time zone;not null;default:CURRENT_TIMESTAMP" json:"createdOn"`
	UpdatedOn  time.Time  			`gorm:"column:updatedOn;type:timestamp(3) without time zone;not null;autoUpdateTime" json:"updatedOn"`
	Data       *string     			`gorm:"column:data;type:jsonb" json:"data"`
	Team       Team     		   	`gorm:"foreignKey:TeamID" json:"team"`
	Parent     *TeamCollection   	`gorm:"foreignKey:ParentID" json:"parent"`
	Children   []TeamCollection   	`gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"children"`
	Requests   []TeamRequest		`gorm:"foreignKey:CollectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"request"`
}

// TableName TeamCollection's table name
func (*TeamCollection) TableName() string {
	return TableNameTeamCollection
}

func (c *TeamCollection) GetTeamID() string {
	return c.TeamID
}

func (c *TeamCollection) Can(db *gorm.DB, uid string, role TeamMemberRole) bool {
	member := &TeamMember{}
	if db.First(member, `"userUid"=? AND "teamID"=?`, uid, c.TeamID).Error != nil {
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

func (*TeamCollection) ParentColName() string {
	return "parentID"
}

func (c *TeamCollection) Move(db *gorm.DB, next *string) error {
	panic("implement me")
}