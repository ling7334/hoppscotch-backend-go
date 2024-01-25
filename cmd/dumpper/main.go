package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	var dataMap = map[string]func(gorm.ColumnType) (dataType string){
		"ReqType": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "*ReqType"
			}
			return "ReqType"
		},

		"TeamMemberRole": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "*TeamMemberRole"
			}
			return "TeamMemberRole"
		},
		"jsonb": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "*JSONB"
			}
			return "JSONB"
		},
	}

	g := gen.NewGenerator(gen.Config{
		// if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable: true,
		// if you want to assign field which has a default value in the `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		FieldCoverable: true,
		// if you want to generate field with unsigned integer type, set FieldSignable true
		FieldSignable: true,
		// if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		// if you want to generate type tags from database, set FieldWithTypeTag true
		FieldWithTypeTag: true,
		// if you need unit tests for query code, set WithUnitTest true
		WithUnitTest: true,
		OutPath:      "../models",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	g.WithDataTypeMap(dataMap)

	gormdb, _ := gorm.Open(postgres.Open("postgres://postgres:example@127.0.0.1:5432/postgres?search_path=public"))
	g.UseDB(gormdb) // reuse your gorm db

	g.ApplyBasic(
		// Generate structs from all tables of current database
		g.GenerateAllTable()...,
	)
	// Generate the code
	g.Execute()
}
