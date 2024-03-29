// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
package model

const TableNameTeam = "Team"

// Team mapped from table <Team>
type Team struct {
	ID           string            `gorm:"column:id;type:text;primaryKey" json:"id"`
	Name         string            `gorm:"column:name;type:text;not null" json:"name"`
	Collections  []TeamCollection  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"collections"`
	Environments []TeamEnvironment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"environments"`
	Invitations  []TeamInvitation  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitations"`
	Teammembers  []TeamMember      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"teamMembers"`
	Requests     []TeamRequest     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"requests"`
}

// TableName Team's table name
func (*Team) TableName() string {
	return TableNameTeam
}
