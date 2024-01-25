// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"model"
)

func newTeamRequest(db *gorm.DB, opts ...gen.DOOption) teamRequest {
	_teamRequest := teamRequest{}

	_teamRequest.teamRequestDo.UseDB(db, opts...)
	_teamRequest.teamRequestDo.UseModel(&model.TeamRequest{})

	tableName := _teamRequest.teamRequestDo.TableName()
	_teamRequest.ALL = field.NewAsterisk(tableName)
	_teamRequest.ID = field.NewString(tableName, "id")
	_teamRequest.CollectionID = field.NewString(tableName, "collectionID")
	_teamRequest.TeamID = field.NewString(tableName, "teamID")
	_teamRequest.Title = field.NewString(tableName, "title")
	_teamRequest.Request = field.NewField(tableName, "request")
	_teamRequest.OrderIndex = field.NewInt32(tableName, "orderIndex")
	_teamRequest.CreatedOn = field.NewTime(tableName, "createdOn")
	_teamRequest.UpdatedOn = field.NewTime(tableName, "updatedOn")

	_teamRequest.fillFieldMap()

	return _teamRequest
}

type teamRequest struct {
	teamRequestDo

	ALL          field.Asterisk
	ID           field.String
	CollectionID field.String
	TeamID       field.String
	Title        field.String
	Request      field.Field
	OrderIndex   field.Int32
	CreatedOn    field.Time
	UpdatedOn    field.Time

	fieldMap map[string]field.Expr
}

func (t teamRequest) Table(newTableName string) *teamRequest {
	t.teamRequestDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t teamRequest) As(alias string) *teamRequest {
	t.teamRequestDo.DO = *(t.teamRequestDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *teamRequest) updateTableName(table string) *teamRequest {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewString(table, "id")
	t.CollectionID = field.NewString(table, "collectionID")
	t.TeamID = field.NewString(table, "teamID")
	t.Title = field.NewString(table, "title")
	t.Request = field.NewField(table, "request")
	t.OrderIndex = field.NewInt32(table, "orderIndex")
	t.CreatedOn = field.NewTime(table, "createdOn")
	t.UpdatedOn = field.NewTime(table, "updatedOn")

	t.fillFieldMap()

	return t
}

func (t *teamRequest) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *teamRequest) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 8)
	t.fieldMap["id"] = t.ID
	t.fieldMap["collectionID"] = t.CollectionID
	t.fieldMap["teamID"] = t.TeamID
	t.fieldMap["title"] = t.Title
	t.fieldMap["request"] = t.Request
	t.fieldMap["orderIndex"] = t.OrderIndex
	t.fieldMap["createdOn"] = t.CreatedOn
	t.fieldMap["updatedOn"] = t.UpdatedOn
}

func (t teamRequest) clone(db *gorm.DB) teamRequest {
	t.teamRequestDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t teamRequest) replaceDB(db *gorm.DB) teamRequest {
	t.teamRequestDo.ReplaceDB(db)
	return t
}

type teamRequestDo struct{ gen.DO }

type ITeamRequestDo interface {
	gen.SubQuery
	Debug() ITeamRequestDo
	WithContext(ctx context.Context) ITeamRequestDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITeamRequestDo
	WriteDB() ITeamRequestDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITeamRequestDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITeamRequestDo
	Not(conds ...gen.Condition) ITeamRequestDo
	Or(conds ...gen.Condition) ITeamRequestDo
	Select(conds ...field.Expr) ITeamRequestDo
	Where(conds ...gen.Condition) ITeamRequestDo
	Order(conds ...field.Expr) ITeamRequestDo
	Distinct(cols ...field.Expr) ITeamRequestDo
	Omit(cols ...field.Expr) ITeamRequestDo
	Join(table schema.Tabler, on ...field.Expr) ITeamRequestDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITeamRequestDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITeamRequestDo
	Group(cols ...field.Expr) ITeamRequestDo
	Having(conds ...gen.Condition) ITeamRequestDo
	Limit(limit int) ITeamRequestDo
	Offset(offset int) ITeamRequestDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITeamRequestDo
	Unscoped() ITeamRequestDo
	Create(values ...*model.TeamRequest) error
	CreateInBatches(values []*model.TeamRequest, batchSize int) error
	Save(values ...*model.TeamRequest) error
	First() (*model.TeamRequest, error)
	Take() (*model.TeamRequest, error)
	Last() (*model.TeamRequest, error)
	Find() ([]*model.TeamRequest, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TeamRequest, err error)
	FindInBatches(result *[]*model.TeamRequest, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TeamRequest) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITeamRequestDo
	Assign(attrs ...field.AssignExpr) ITeamRequestDo
	Joins(fields ...field.RelationField) ITeamRequestDo
	Preload(fields ...field.RelationField) ITeamRequestDo
	FirstOrInit() (*model.TeamRequest, error)
	FirstOrCreate() (*model.TeamRequest, error)
	FindByPage(offset int, limit int) (result []*model.TeamRequest, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITeamRequestDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t teamRequestDo) Debug() ITeamRequestDo {
	return t.withDO(t.DO.Debug())
}

func (t teamRequestDo) WithContext(ctx context.Context) ITeamRequestDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t teamRequestDo) ReadDB() ITeamRequestDo {
	return t.Clauses(dbresolver.Read)
}

func (t teamRequestDo) WriteDB() ITeamRequestDo {
	return t.Clauses(dbresolver.Write)
}

func (t teamRequestDo) Session(config *gorm.Session) ITeamRequestDo {
	return t.withDO(t.DO.Session(config))
}

func (t teamRequestDo) Clauses(conds ...clause.Expression) ITeamRequestDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t teamRequestDo) Returning(value interface{}, columns ...string) ITeamRequestDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t teamRequestDo) Not(conds ...gen.Condition) ITeamRequestDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t teamRequestDo) Or(conds ...gen.Condition) ITeamRequestDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t teamRequestDo) Select(conds ...field.Expr) ITeamRequestDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t teamRequestDo) Where(conds ...gen.Condition) ITeamRequestDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t teamRequestDo) Order(conds ...field.Expr) ITeamRequestDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t teamRequestDo) Distinct(cols ...field.Expr) ITeamRequestDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t teamRequestDo) Omit(cols ...field.Expr) ITeamRequestDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t teamRequestDo) Join(table schema.Tabler, on ...field.Expr) ITeamRequestDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t teamRequestDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITeamRequestDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t teamRequestDo) RightJoin(table schema.Tabler, on ...field.Expr) ITeamRequestDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t teamRequestDo) Group(cols ...field.Expr) ITeamRequestDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t teamRequestDo) Having(conds ...gen.Condition) ITeamRequestDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t teamRequestDo) Limit(limit int) ITeamRequestDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t teamRequestDo) Offset(offset int) ITeamRequestDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t teamRequestDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITeamRequestDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t teamRequestDo) Unscoped() ITeamRequestDo {
	return t.withDO(t.DO.Unscoped())
}

func (t teamRequestDo) Create(values ...*model.TeamRequest) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t teamRequestDo) CreateInBatches(values []*model.TeamRequest, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t teamRequestDo) Save(values ...*model.TeamRequest) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t teamRequestDo) First() (*model.TeamRequest, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamRequest), nil
	}
}

func (t teamRequestDo) Take() (*model.TeamRequest, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamRequest), nil
	}
}

func (t teamRequestDo) Last() (*model.TeamRequest, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamRequest), nil
	}
}

func (t teamRequestDo) Find() ([]*model.TeamRequest, error) {
	result, err := t.DO.Find()
	return result.([]*model.TeamRequest), err
}

func (t teamRequestDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TeamRequest, err error) {
	buf := make([]*model.TeamRequest, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t teamRequestDo) FindInBatches(result *[]*model.TeamRequest, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t teamRequestDo) Attrs(attrs ...field.AssignExpr) ITeamRequestDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t teamRequestDo) Assign(attrs ...field.AssignExpr) ITeamRequestDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t teamRequestDo) Joins(fields ...field.RelationField) ITeamRequestDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t teamRequestDo) Preload(fields ...field.RelationField) ITeamRequestDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t teamRequestDo) FirstOrInit() (*model.TeamRequest, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamRequest), nil
	}
}

func (t teamRequestDo) FirstOrCreate() (*model.TeamRequest, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamRequest), nil
	}
}

func (t teamRequestDo) FindByPage(offset int, limit int) (result []*model.TeamRequest, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t teamRequestDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t teamRequestDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t teamRequestDo) Delete(models ...*model.TeamRequest) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *teamRequestDo) withDO(do gen.Dao) *teamRequestDo {
	t.DO = *do.(*gen.DO)
	return t
}
