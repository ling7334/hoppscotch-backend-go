// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
package model

const TableNameTeam = "Team"

// Team mapped from table <Team>
type Team struct {
	ID           string            `gorm:"column:id;type:text;primaryKey" json:"id"`
	Name         string            `gorm:"column:name;type:text;not null" json:"name"`
	Collections  []TeamCollection  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"collection"`
	Environments []TeamEnvironment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"environment"`
	Invitations  []TeamInvitation  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"invitation"`
	Members      []TeamMember      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"member"`
	Requests     []TeamRequest     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"request"`
}

// TableName Team's table name
func (*Team) TableName() string {
	return TableNameTeam
}
