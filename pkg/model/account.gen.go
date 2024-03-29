// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAccount = "Account"

// Account mapped from table <Account>
type Account struct {
	ID                   string     `gorm:"column:id;type:text;primaryKey" json:"id"`
	UserID               string     `gorm:"column:userId;type:text;not null" json:"userId"`
	Provider             string     `gorm:"column:provider;type:text;not null;uniqueIndex:Account_provider_providerAccountId_key,priority:1" json:"provider"`
	ProviderAccountID    string     `gorm:"column:providerAccountId;type:text;not null;uniqueIndex:Account_provider_providerAccountId_key,priority:2" json:"providerAccountId"`
	ProviderRefreshToken *string    `gorm:"column:providerRefreshToken;type:text" json:"providerRefreshToken"`
	ProviderAccessToken  *string    `gorm:"column:providerAccessToken;type:text" json:"providerAccessToken"`
	ProviderScope        *string    `gorm:"column:providerScope;type:text" json:"providerScope"`
	LoggedIn             time.Time  `gorm:"column:loggedIn;type:timestamp(3) without time zone;not null;default:CURRENT_TIMESTAMP" json:"loggedIn"`
	User                 User       `gorm:"foreignKey:UserID" json:"user"`
}

// TableName Account's table name
func (*Account) TableName() string {
	return TableNameAccount
}
