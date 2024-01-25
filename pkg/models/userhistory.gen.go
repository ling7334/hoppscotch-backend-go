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

func newUserHistory(db *gorm.DB, opts ...gen.DOOption) userHistory {
	_userHistory := userHistory{}

	_userHistory.userHistoryDo.UseDB(db, opts...)
	_userHistory.userHistoryDo.UseModel(&model.UserHistory{})

	tableName := _userHistory.userHistoryDo.TableName()
	_userHistory.ALL = field.NewAsterisk(tableName)
	_userHistory.ID = field.NewString(tableName, "id")
	_userHistory.UserUID = field.NewString(tableName, "userUid")
	_userHistory.ReqType = field.NewField(tableName, "reqType")
	_userHistory.Request = field.NewField(tableName, "request")
	_userHistory.ResponseMetadata = field.NewField(tableName, "responseMetadata")
	_userHistory.IsStarred = field.NewBool(tableName, "isStarred")
	_userHistory.ExecutedOn = field.NewTime(tableName, "executedOn")

	_userHistory.fillFieldMap()

	return _userHistory
}

type userHistory struct {
	userHistoryDo

	ALL              field.Asterisk
	ID               field.String
	UserUID          field.String
	ReqType          field.Field
	Request          field.Field
	ResponseMetadata field.Field
	IsStarred        field.Bool
	ExecutedOn       field.Time

	fieldMap map[string]field.Expr
}

func (u userHistory) Table(newTableName string) *userHistory {
	u.userHistoryDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userHistory) As(alias string) *userHistory {
	u.userHistoryDo.DO = *(u.userHistoryDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userHistory) updateTableName(table string) *userHistory {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewString(table, "id")
	u.UserUID = field.NewString(table, "userUid")
	u.ReqType = field.NewField(table, "reqType")
	u.Request = field.NewField(table, "request")
	u.ResponseMetadata = field.NewField(table, "responseMetadata")
	u.IsStarred = field.NewBool(table, "isStarred")
	u.ExecutedOn = field.NewTime(table, "executedOn")

	u.fillFieldMap()

	return u
}

func (u *userHistory) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userHistory) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 7)
	u.fieldMap["id"] = u.ID
	u.fieldMap["userUid"] = u.UserUID
	u.fieldMap["reqType"] = u.ReqType
	u.fieldMap["request"] = u.Request
	u.fieldMap["responseMetadata"] = u.ResponseMetadata
	u.fieldMap["isStarred"] = u.IsStarred
	u.fieldMap["executedOn"] = u.ExecutedOn
}

func (u userHistory) clone(db *gorm.DB) userHistory {
	u.userHistoryDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userHistory) replaceDB(db *gorm.DB) userHistory {
	u.userHistoryDo.ReplaceDB(db)
	return u
}

type userHistoryDo struct{ gen.DO }

type IUserHistoryDo interface {
	gen.SubQuery
	Debug() IUserHistoryDo
	WithContext(ctx context.Context) IUserHistoryDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserHistoryDo
	WriteDB() IUserHistoryDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserHistoryDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserHistoryDo
	Not(conds ...gen.Condition) IUserHistoryDo
	Or(conds ...gen.Condition) IUserHistoryDo
	Select(conds ...field.Expr) IUserHistoryDo
	Where(conds ...gen.Condition) IUserHistoryDo
	Order(conds ...field.Expr) IUserHistoryDo
	Distinct(cols ...field.Expr) IUserHistoryDo
	Omit(cols ...field.Expr) IUserHistoryDo
	Join(table schema.Tabler, on ...field.Expr) IUserHistoryDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserHistoryDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserHistoryDo
	Group(cols ...field.Expr) IUserHistoryDo
	Having(conds ...gen.Condition) IUserHistoryDo
	Limit(limit int) IUserHistoryDo
	Offset(offset int) IUserHistoryDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserHistoryDo
	Unscoped() IUserHistoryDo
	Create(values ...*model.UserHistory) error
	CreateInBatches(values []*model.UserHistory, batchSize int) error
	Save(values ...*model.UserHistory) error
	First() (*model.UserHistory, error)
	Take() (*model.UserHistory, error)
	Last() (*model.UserHistory, error)
	Find() ([]*model.UserHistory, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserHistory, err error)
	FindInBatches(result *[]*model.UserHistory, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.UserHistory) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserHistoryDo
	Assign(attrs ...field.AssignExpr) IUserHistoryDo
	Joins(fields ...field.RelationField) IUserHistoryDo
	Preload(fields ...field.RelationField) IUserHistoryDo
	FirstOrInit() (*model.UserHistory, error)
	FirstOrCreate() (*model.UserHistory, error)
	FindByPage(offset int, limit int) (result []*model.UserHistory, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserHistoryDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (u userHistoryDo) Debug() IUserHistoryDo {
	return u.withDO(u.DO.Debug())
}

func (u userHistoryDo) WithContext(ctx context.Context) IUserHistoryDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userHistoryDo) ReadDB() IUserHistoryDo {
	return u.Clauses(dbresolver.Read)
}

func (u userHistoryDo) WriteDB() IUserHistoryDo {
	return u.Clauses(dbresolver.Write)
}

func (u userHistoryDo) Session(config *gorm.Session) IUserHistoryDo {
	return u.withDO(u.DO.Session(config))
}

func (u userHistoryDo) Clauses(conds ...clause.Expression) IUserHistoryDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userHistoryDo) Returning(value interface{}, columns ...string) IUserHistoryDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userHistoryDo) Not(conds ...gen.Condition) IUserHistoryDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userHistoryDo) Or(conds ...gen.Condition) IUserHistoryDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userHistoryDo) Select(conds ...field.Expr) IUserHistoryDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userHistoryDo) Where(conds ...gen.Condition) IUserHistoryDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userHistoryDo) Order(conds ...field.Expr) IUserHistoryDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userHistoryDo) Distinct(cols ...field.Expr) IUserHistoryDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userHistoryDo) Omit(cols ...field.Expr) IUserHistoryDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userHistoryDo) Join(table schema.Tabler, on ...field.Expr) IUserHistoryDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userHistoryDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserHistoryDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userHistoryDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserHistoryDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userHistoryDo) Group(cols ...field.Expr) IUserHistoryDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userHistoryDo) Having(conds ...gen.Condition) IUserHistoryDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userHistoryDo) Limit(limit int) IUserHistoryDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userHistoryDo) Offset(offset int) IUserHistoryDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userHistoryDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserHistoryDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userHistoryDo) Unscoped() IUserHistoryDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userHistoryDo) Create(values ...*model.UserHistory) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userHistoryDo) CreateInBatches(values []*model.UserHistory, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userHistoryDo) Save(values ...*model.UserHistory) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userHistoryDo) First() (*model.UserHistory, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserHistory), nil
	}
}

func (u userHistoryDo) Take() (*model.UserHistory, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserHistory), nil
	}
}

func (u userHistoryDo) Last() (*model.UserHistory, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserHistory), nil
	}
}

func (u userHistoryDo) Find() ([]*model.UserHistory, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserHistory), err
}

func (u userHistoryDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserHistory, err error) {
	buf := make([]*model.UserHistory, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userHistoryDo) FindInBatches(result *[]*model.UserHistory, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userHistoryDo) Attrs(attrs ...field.AssignExpr) IUserHistoryDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userHistoryDo) Assign(attrs ...field.AssignExpr) IUserHistoryDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userHistoryDo) Joins(fields ...field.RelationField) IUserHistoryDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userHistoryDo) Preload(fields ...field.RelationField) IUserHistoryDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userHistoryDo) FirstOrInit() (*model.UserHistory, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserHistory), nil
	}
}

func (u userHistoryDo) FirstOrCreate() (*model.UserHistory, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserHistory), nil
	}
}

func (u userHistoryDo) FindByPage(offset int, limit int) (result []*model.UserHistory, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userHistoryDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userHistoryDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userHistoryDo) Delete(models ...*model.UserHistory) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userHistoryDo) withDO(do gen.Dao) *userHistoryDo {
	u.DO = *do.(*gen.DO)
	return u
}
