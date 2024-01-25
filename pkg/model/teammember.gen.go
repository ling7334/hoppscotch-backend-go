// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"gorm.io/gorm"
)

const TableNameTeamMember = "TeamMember"

// TeamMember mapped from table <TeamMember>
type TeamMember struct {
	ID      string         `gorm:"column:id;type:text;primaryKey" json:"id"`
	Role    TeamMemberRole `gorm:"column:role;type:team_member_role;not null" json:"role"`
	UserUID string         `gorm:"column:userUid;type:text;not null;uniqueIndex:TeamMember_teamID_userUid_key,priority:1" json:"userUid"`
	TeamID  string         `gorm:"column:teamID;type:text;not null;uniqueIndex:TeamMember_teamID_userUid_key,priority:2" json:"teamID"`
	Team    Team     		`gorm:"foreignKey:TeamID" json:"team"`
	User    User     		`gorm:"foreignKey:UserUID" json:"user"`
}

// TableName TeamMember's table name
func (*TeamMember) TableName() string {
	return TableNameTeamMember
}

func (m *TeamMember) GetTeamID() string {
	return m.TeamID
}

func (m *TeamMember) Can(db *gorm.DB, uid string, role TeamMemberRole) bool {
	member := &TeamMember{}
	if db.First(member, `"userUid"=? AND "teamID"=?`, uid, m.TeamID).Error != nil {
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