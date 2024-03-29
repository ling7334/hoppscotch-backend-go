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

func newPrismaMigration(db *gorm.DB, opts ...gen.DOOption) prismaMigration {
	_prismaMigration := prismaMigration{}

	_prismaMigration.prismaMigrationDo.UseDB(db, opts...)
	_prismaMigration.prismaMigrationDo.UseModel(&model.PrismaMigration{})

	tableName := _prismaMigration.prismaMigrationDo.TableName()
	_prismaMigration.ALL = field.NewAsterisk(tableName)
	_prismaMigration.ID = field.NewString(tableName, "id")
	_prismaMigration.Checksum = field.NewString(tableName, "checksum")
	_prismaMigration.FinishedAt = field.NewTime(tableName, "finished_at")
	_prismaMigration.MigrationName = field.NewString(tableName, "migration_name")
	_prismaMigration.Logs = field.NewString(tableName, "logs")
	_prismaMigration.RolledBackAt = field.NewTime(tableName, "rolled_back_at")
	_prismaMigration.StartedAt = field.NewTime(tableName, "started_at")
	_prismaMigration.AppliedStepsCount = field.NewInt32(tableName, "applied_steps_count")

	_prismaMigration.fillFieldMap()

	return _prismaMigration
}

type prismaMigration struct {
	prismaMigrationDo

	ALL               field.Asterisk
	ID                field.String
	Checksum          field.String
	FinishedAt        field.Time
	MigrationName     field.String
	Logs              field.String
	RolledBackAt      field.Time
	StartedAt         field.Time
	AppliedStepsCount field.Int32

	fieldMap map[string]field.Expr
}

func (p prismaMigration) Table(newTableName string) *prismaMigration {
	p.prismaMigrationDo.UseTable(newTableName)
	return p.updateTableName(newTableName)
}

func (p prismaMigration) As(alias string) *prismaMigration {
	p.prismaMigrationDo.DO = *(p.prismaMigrationDo.As(alias).(*gen.DO))
	return p.updateTableName(alias)
}

func (p *prismaMigration) updateTableName(table string) *prismaMigration {
	p.ALL = field.NewAsterisk(table)
	p.ID = field.NewString(table, "id")
	p.Checksum = field.NewString(table, "checksum")
	p.FinishedAt = field.NewTime(table, "finished_at")
	p.MigrationName = field.NewString(table, "migration_name")
	p.Logs = field.NewString(table, "logs")
	p.RolledBackAt = field.NewTime(table, "rolled_back_at")
	p.StartedAt = field.NewTime(table, "started_at")
	p.AppliedStepsCount = field.NewInt32(table, "applied_steps_count")

	p.fillFieldMap()

	return p
}

func (p *prismaMigration) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := p.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (p *prismaMigration) fillFieldMap() {
	p.fieldMap = make(map[string]field.Expr, 8)
	p.fieldMap["id"] = p.ID
	p.fieldMap["checksum"] = p.Checksum
	p.fieldMap["finished_at"] = p.FinishedAt
	p.fieldMap["migration_name"] = p.MigrationName
	p.fieldMap["logs"] = p.Logs
	p.fieldMap["rolled_back_at"] = p.RolledBackAt
	p.fieldMap["started_at"] = p.StartedAt
	p.fieldMap["applied_steps_count"] = p.AppliedStepsCount
}

func (p prismaMigration) clone(db *gorm.DB) prismaMigration {
	p.prismaMigrationDo.ReplaceConnPool(db.Statement.ConnPool)
	return p
}

func (p prismaMigration) replaceDB(db *gorm.DB) prismaMigration {
	p.prismaMigrationDo.ReplaceDB(db)
	return p
}

type prismaMigrationDo struct{ gen.DO }

type IPrismaMigrationDo interface {
	gen.SubQuery
	Debug() IPrismaMigrationDo
	WithContext(ctx context.Context) IPrismaMigrationDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IPrismaMigrationDo
	WriteDB() IPrismaMigrationDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IPrismaMigrationDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IPrismaMigrationDo
	Not(conds ...gen.Condition) IPrismaMigrationDo
	Or(conds ...gen.Condition) IPrismaMigrationDo
	Select(conds ...field.Expr) IPrismaMigrationDo
	Where(conds ...gen.Condition) IPrismaMigrationDo
	Order(conds ...field.Expr) IPrismaMigrationDo
	Distinct(cols ...field.Expr) IPrismaMigrationDo
	Omit(cols ...field.Expr) IPrismaMigrationDo
	Join(table schema.Tabler, on ...field.Expr) IPrismaMigrationDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IPrismaMigrationDo
	RightJoin(table schema.Tabler, on ...field.Expr) IPrismaMigrationDo
	Group(cols ...field.Expr) IPrismaMigrationDo
	Having(conds ...gen.Condition) IPrismaMigrationDo
	Limit(limit int) IPrismaMigrationDo
	Offset(offset int) IPrismaMigrationDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IPrismaMigrationDo
	Unscoped() IPrismaMigrationDo
	Create(values ...*model.PrismaMigration) error
	CreateInBatches(values []*model.PrismaMigration, batchSize int) error
	Save(values ...*model.PrismaMigration) error
	First() (*model.PrismaMigration, error)
	Take() (*model.PrismaMigration, error)
	Last() (*model.PrismaMigration, error)
	Find() ([]*model.PrismaMigration, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.PrismaMigration, err error)
	FindInBatches(result *[]*model.PrismaMigration, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.PrismaMigration) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IPrismaMigrationDo
	Assign(attrs ...field.AssignExpr) IPrismaMigrationDo
	Joins(fields ...field.RelationField) IPrismaMigrationDo
	Preload(fields ...field.RelationField) IPrismaMigrationDo
	FirstOrInit() (*model.PrismaMigration, error)
	FirstOrCreate() (*model.PrismaMigration, error)
	FindByPage(offset int, limit int) (result []*model.PrismaMigration, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IPrismaMigrationDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (p prismaMigrationDo) Debug() IPrismaMigrationDo {
	return p.withDO(p.DO.Debug())
}

func (p prismaMigrationDo) WithContext(ctx context.Context) IPrismaMigrationDo {
	return p.withDO(p.DO.WithContext(ctx))
}

func (p prismaMigrationDo) ReadDB() IPrismaMigrationDo {
	return p.Clauses(dbresolver.Read)
}

func (p prismaMigrationDo) WriteDB() IPrismaMigrationDo {
	return p.Clauses(dbresolver.Write)
}

func (p prismaMigrationDo) Session(config *gorm.Session) IPrismaMigrationDo {
	return p.withDO(p.DO.Session(config))
}

func (p prismaMigrationDo) Clauses(conds ...clause.Expression) IPrismaMigrationDo {
	return p.withDO(p.DO.Clauses(conds...))
}

func (p prismaMigrationDo) Returning(value interface{}, columns ...string) IPrismaMigrationDo {
	return p.withDO(p.DO.Returning(value, columns...))
}

func (p prismaMigrationDo) Not(conds ...gen.Condition) IPrismaMigrationDo {
	return p.withDO(p.DO.Not(conds...))
}

func (p prismaMigrationDo) Or(conds ...gen.Condition) IPrismaMigrationDo {
	return p.withDO(p.DO.Or(conds...))
}

func (p prismaMigrationDo) Select(conds ...field.Expr) IPrismaMigrationDo {
	return p.withDO(p.DO.Select(conds...))
}

func (p prismaMigrationDo) Where(conds ...gen.Condition) IPrismaMigrationDo {
	return p.withDO(p.DO.Where(conds...))
}

func (p prismaMigrationDo) Order(conds ...field.Expr) IPrismaMigrationDo {
	return p.withDO(p.DO.Order(conds...))
}

func (p prismaMigrationDo) Distinct(cols ...field.Expr) IPrismaMigrationDo {
	return p.withDO(p.DO.Distinct(cols...))
}

func (p prismaMigrationDo) Omit(cols ...field.Expr) IPrismaMigrationDo {
	return p.withDO(p.DO.Omit(cols...))
}

func (p prismaMigrationDo) Join(table schema.Tabler, on ...field.Expr) IPrismaMigrationDo {
	return p.withDO(p.DO.Join(table, on...))
}

func (p prismaMigrationDo) LeftJoin(table schema.Tabler, on ...field.Expr) IPrismaMigrationDo {
	return p.withDO(p.DO.LeftJoin(table, on...))
}

func (p prismaMigrationDo) RightJoin(table schema.Tabler, on ...field.Expr) IPrismaMigrationDo {
	return p.withDO(p.DO.RightJoin(table, on...))
}

func (p prismaMigrationDo) Group(cols ...field.Expr) IPrismaMigrationDo {
	return p.withDO(p.DO.Group(cols...))
}

func (p prismaMigrationDo) Having(conds ...gen.Condition) IPrismaMigrationDo {
	return p.withDO(p.DO.Having(conds...))
}

func (p prismaMigrationDo) Limit(limit int) IPrismaMigrationDo {
	return p.withDO(p.DO.Limit(limit))
}

func (p prismaMigrationDo) Offset(offset int) IPrismaMigrationDo {
	return p.withDO(p.DO.Offset(offset))
}

func (p prismaMigrationDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IPrismaMigrationDo {
	return p.withDO(p.DO.Scopes(funcs...))
}

func (p prismaMigrationDo) Unscoped() IPrismaMigrationDo {
	return p.withDO(p.DO.Unscoped())
}

func (p prismaMigrationDo) Create(values ...*model.PrismaMigration) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Create(values)
}

func (p prismaMigrationDo) CreateInBatches(values []*model.PrismaMigration, batchSize int) error {
	return p.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (p prismaMigrationDo) Save(values ...*model.PrismaMigration) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Save(values)
}

func (p prismaMigrationDo) First() (*model.PrismaMigration, error) {
	if result, err := p.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.PrismaMigration), nil
	}
}

func (p prismaMigrationDo) Take() (*model.PrismaMigration, error) {
	if result, err := p.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.PrismaMigration), nil
	}
}

func (p prismaMigrationDo) Last() (*model.PrismaMigration, error) {
	if result, err := p.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.PrismaMigration), nil
	}
}

func (p prismaMigrationDo) Find() ([]*model.PrismaMigration, error) {
	result, err := p.DO.Find()
	return result.([]*model.PrismaMigration), err
}

func (p prismaMigrationDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.PrismaMigration, err error) {
	buf := make([]*model.PrismaMigration, 0, batchSize)
	err = p.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (p prismaMigrationDo) FindInBatches(result *[]*model.PrismaMigration, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return p.DO.FindInBatches(result, batchSize, fc)
}

func (p prismaMigrationDo) Attrs(attrs ...field.AssignExpr) IPrismaMigrationDo {
	return p.withDO(p.DO.Attrs(attrs...))
}

func (p prismaMigrationDo) Assign(attrs ...field.AssignExpr) IPrismaMigrationDo {
	return p.withDO(p.DO.Assign(attrs...))
}

func (p prismaMigrationDo) Joins(fields ...field.RelationField) IPrismaMigrationDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Joins(_f))
	}
	return &p
}

func (p prismaMigrationDo) Preload(fields ...field.RelationField) IPrismaMigrationDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Preload(_f))
	}
	return &p
}

func (p prismaMigrationDo) FirstOrInit() (*model.PrismaMigration, error) {
	if result, err := p.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.PrismaMigration), nil
	}
}

func (p prismaMigrationDo) FirstOrCreate() (*model.PrismaMigration, error) {
	if result, err := p.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.PrismaMigration), nil
	}
}

func (p prismaMigrationDo) FindByPage(offset int, limit int) (result []*model.PrismaMigration, count int64, err error) {
	result, err = p.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = p.Offset(-1).Limit(-1).Count()
	return
}

func (p prismaMigrationDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = p.Count()
	if err != nil {
		return
	}

	err = p.Offset(offset).Limit(limit).Scan(result)
	return
}

func (p prismaMigrationDo) Scan(result interface{}) (err error) {
	return p.DO.Scan(result)
}

func (p prismaMigrationDo) Delete(models ...*model.PrismaMigration) (result gen.ResultInfo, err error) {
	return p.DO.Delete(models)
}

func (p *prismaMigrationDo) withDO(do gen.Dao) *prismaMigrationDo {
	p.DO = *do.(*gen.DO)
	return p
}
