// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"context"
	"fmt"
	"testing"

	"model"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

func init() {
	InitializeDB()
	err := db.AutoMigrate(&model.TeamCollection{})
	if err != nil {
		fmt.Printf("Error: AutoMigrate(&model.TeamCollection{}) fail: %s", err)
	}
}

func Test_teamCollectionQuery(t *testing.T) {
	teamCollection := newTeamCollection(db)
	teamCollection = *teamCollection.As(teamCollection.TableName())
	_do := teamCollection.WithContext(context.Background()).Debug()

	primaryKey := field.NewString(teamCollection.TableName(), clause.PrimaryKey)
	_, err := _do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <TeamCollection> fail:", err)
		return
	}

	_, ok := teamCollection.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from teamCollection success")
	}

	err = _do.Create(&model.TeamCollection{})
	if err != nil {
		t.Error("create item in table <TeamCollection> fail:", err)
	}

	err = _do.Save(&model.TeamCollection{})
	if err != nil {
		t.Error("create item in table <TeamCollection> fail:", err)
	}

	err = _do.CreateInBatches([]*model.TeamCollection{{}, {}}, 10)
	if err != nil {
		t.Error("create item in table <TeamCollection> fail:", err)
	}

	_, err = _do.Select(teamCollection.ALL).Take()
	if err != nil {
		t.Error("Take() on table <TeamCollection> fail:", err)
	}

	_, err = _do.First()
	if err != nil {
		t.Error("First() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Last()
	if err != nil {
		t.Error("First() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <TeamCollection> fail:", err)
	}

	err = _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*model.TeamCollection{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Select(teamCollection.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Select(teamCollection.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <TeamCollection> fail:", err)
	}

	_, _, err = _do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <TeamCollection> fail:", err)
	}

	_, err = _do.ScanByPage(&model.TeamCollection{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <TeamCollection> fail:", err)
	}

	var _a _another
	var _aPK = field.NewString(_a.TableName(), "id")

	err = _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <TeamCollection> fail:", err)
	}

	err = _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <TeamCollection> fail:", err)
	}

	_, err = _do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <TeamCollection> fail:", err)
	}
}
