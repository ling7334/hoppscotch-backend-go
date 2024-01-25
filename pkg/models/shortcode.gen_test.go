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
	err := db.AutoMigrate(&model.Shortcode{})
	if err != nil {
		fmt.Printf("Error: AutoMigrate(&model.Shortcode{}) fail: %s", err)
	}
}

func Test_shortcodeQuery(t *testing.T) {
	shortcode := newShortcode(db)
	shortcode = *shortcode.As(shortcode.TableName())
	_do := shortcode.WithContext(context.Background()).Debug()

	primaryKey := field.NewString(shortcode.TableName(), clause.PrimaryKey)
	_, err := _do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <Shortcode> fail:", err)
		return
	}

	_, ok := shortcode.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from shortcode success")
	}

	err = _do.Create(&model.Shortcode{})
	if err != nil {
		t.Error("create item in table <Shortcode> fail:", err)
	}

	err = _do.Save(&model.Shortcode{})
	if err != nil {
		t.Error("create item in table <Shortcode> fail:", err)
	}

	err = _do.CreateInBatches([]*model.Shortcode{{}, {}}, 10)
	if err != nil {
		t.Error("create item in table <Shortcode> fail:", err)
	}

	_, err = _do.Select(shortcode.ALL).Take()
	if err != nil {
		t.Error("Take() on table <Shortcode> fail:", err)
	}

	_, err = _do.First()
	if err != nil {
		t.Error("First() on table <Shortcode> fail:", err)
	}

	_, err = _do.Last()
	if err != nil {
		t.Error("First() on table <Shortcode> fail:", err)
	}

	_, err = _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <Shortcode> fail:", err)
	}

	err = _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*model.Shortcode{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <Shortcode> fail:", err)
	}

	_, err = _do.Select(shortcode.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <Shortcode> fail:", err)
	}

	_, err = _do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <Shortcode> fail:", err)
	}

	_, err = _do.Select(shortcode.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <Shortcode> fail:", err)
	}

	_, err = _do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <Shortcode> fail:", err)
	}

	_, err = _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <Shortcode> fail:", err)
	}

	_, _, err = _do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <Shortcode> fail:", err)
	}

	_, err = _do.ScanByPage(&model.Shortcode{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <Shortcode> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <Shortcode> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <Shortcode> fail:", err)
	}

	var _a _another
	var _aPK = field.NewString(_a.TableName(), "id")

	err = _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <Shortcode> fail:", err)
	}

	err = _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <Shortcode> fail:", err)
	}

	_, err = _do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <Shortcode> fail:", err)
	}
}
