package gormrepo

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/banbridge/common/repository"
)

func wrapDBErr(op string, err error) error {
	if isConflictError(err) {
		return fmt.Errorf("%w: %w: %s: %w", repository.ErrConflict, repository.ErrDBFailed, op, err)
	}
	return fmt.Errorf("%w: %s: %w", repository.ErrDBFailed, op, err)
}

func wrapNotFoundErr(err error) error {
	return fmt.Errorf("%w: get record from database failed: %w", repository.ErrNotFound, err)
}

func invalidArgf(format string, args ...any) error {
	return fmt.Errorf("gormrepo: "+format, args...)
}

// CreateModels inserts rows in batches.
func CreateModels[M any](ctx context.Context, db *gorm.DB, rows ...M) (int, error) {
	if db == nil {
		return 0, invalidArgf("db is nil")
	}
	if len(rows) == 0 {
		return 0, nil
	}
	db = db.Session(&gorm.Session{CreateBatchSize: 200})
	var model M
	res := db.WithContext(ctx).
		Model(&model).
		Create(rows)
	if res.Error != nil {
		return 0, wrapDBErr("insert record into database failed", res.Error)
	}
	return int(res.RowsAffected), nil
}

// CreateModelsWithConds inserts pointer rows in batches with extra clauses.
func CreateModelsWithConds[M any](ctx context.Context, db *gorm.DB, rows []*M, conds ...clause.Expression) (int, error) {
	if db == nil {
		return 0, invalidArgf("db is nil")
	}
	if len(rows) == 0 {
		return 0, nil
	}
	res := db.WithContext(ctx).
		Clauses(conds...).
		CreateInBatches(rows, 100)
	if res.Error != nil {
		return 0, wrapDBErr("insert record into database failed", res.Error)
	}
	return int(res.RowsAffected), nil
}

// DeleteModels deletes rows matching the given scopes.
func DeleteModels[M any](ctx context.Context, db *gorm.DB, _ *M,
	conds ...func(*gorm.DB) *gorm.DB) (int, error) {
	if db == nil {
		return 0, invalidArgf("db is nil")
	}
	var m M
	res := db.WithContext(ctx).
		Scopes(conds...).
		Delete(&m)
	err := ignoreNotFound(res.Error)
	if err != nil {
		return 0, wrapDBErr("delete record from database failed", err)
	}
	return int(res.RowsAffected), err
}

// UpdateModelWithConds updates fields on rows matching the given scopes.
func UpdateModelWithConds[M any](ctx context.Context, db *gorm.DB, _ *M, fields map[string]any,
	conds ...func(*gorm.DB) *gorm.DB) (int, error) {
	if db == nil {
		return 0, invalidArgf("db is nil")
	}
	if len(fields) == 0 {
		return 0, invalidArgf("update fields is empty")
	}
	var m M
	db = db.WithContext(ctx).
		Model(&m).
		Scopes(conds...).
		Updates(fields)
	count, err := db.RowsAffected, db.Error
	if err != nil {
		return 0, wrapDBErr("update record into database failed", err)
	}
	return int(count), nil
}

// UpdatesModelWithValue updates rows matching the given scopes with a model.
func UpdatesModelWithValue[M any](ctx context.Context, db *gorm.DB, _ *M, v *M,
	conds ...func(*gorm.DB) *gorm.DB) (int, error) {
	if db == nil {
		return 0, invalidArgf("db is nil")
	}
	if v == nil {
		return 0, invalidArgf("update model is nil")
	}
	var m M
	db = db.WithContext(ctx).
		Model(&m).
		Scopes(conds...).
		Updates(v)

	count, err := db.RowsAffected, db.Error
	if err != nil {
		return 0, wrapDBErr("update record into database failed", err)
	}
	return int(count), nil
}

// GetModelByConds returns the first row matching the given scopes.
func GetModelByConds[M any](ctx context.Context, db *gorm.DB, model *M,
	conds ...func(*gorm.DB) *gorm.DB) (*M, error) {
	if db == nil {
		return nil, invalidArgf("db is nil")
	}
	if model == nil {
		return nil, invalidArgf("model is nil")
	}
	var res M
	err := db.WithContext(ctx).
		Model(model).
		Scopes(conds...).
		First(&res).Error
	if err != nil {
		if isNotFoundError(err) {
			return nil, wrapNotFoundErr(err)
		}
		return nil, wrapDBErr("get record from database failed", err)
	}
	return &res, nil
}

// ListModelsByConds returns all rows matching the given scopes.
func ListModelsByConds[M any](ctx context.Context, db *gorm.DB, model *M,
	conds ...func(*gorm.DB) *gorm.DB) ([]M, error) {
	if db == nil {
		return nil, invalidArgf("db is nil")
	}
	if model == nil {
		return nil, invalidArgf("model is nil")
	}
	var res []M
	err := db.WithContext(ctx).
		Model(model).
		Scopes(conds...).
		Find(&res).Error
	if err != nil {
		return nil, wrapDBErr("list record from database failed", err)
	}
	return res, nil
}

// ListPointerModelsByConds returns pointer rows matching the given scopes.
func ListPointerModelsByConds[M any](ctx context.Context, db *gorm.DB, model *M,
	conds ...func(*gorm.DB) *gorm.DB) ([]*M, error) {
	if db == nil {
		return nil, invalidArgf("db is nil")
	}
	if model == nil {
		return nil, invalidArgf("model is nil")
	}
	var res []*M
	err := db.WithContext(ctx).
		Model(model).
		Scopes(conds...).
		Find(&res).Error
	if err != nil {
		return nil, wrapDBErr("list record from database failed", err)
	}
	return res, nil
}

// CountByConds counts rows matching the given scopes, ignoring paging scopes.
func CountByConds[M any](ctx context.Context, db *gorm.DB, model *M,
	conds ...func(*gorm.DB) *gorm.DB) (int64, error) {
	if db == nil {
		return 0, invalidArgf("db is nil")
	}
	if model == nil {
		return 0, invalidArgf("model is nil")
	}
	var total int64
	err := db.WithContext(ctx).
		Model(model).
		Scopes(conds...).
		Scopes(func(tx *gorm.DB) *gorm.DB {
			return tx.Limit(-1).Offset(-1)
		}).
		Count(&total).Error
	if err != nil {
		return 0, wrapDBErr("count record from database failed", err)
	}
	return total, nil
}

// GroupModelsByConds returns grouped rows matching the given scopes.
func GroupModelsByConds[M, M2 any](ctx context.Context, db *gorm.DB, model *M, groupModel M2,
	groupSelect string, groupColumns string, conds ...func(*gorm.DB) *gorm.DB) ([]M2, error) {
	if db == nil {
		return nil, invalidArgf("db is nil")
	}
	if model == nil {
		return nil, invalidArgf("model is nil")
	}
	var res []M2
	err := db.WithContext(ctx).
		Model(model).
		Select(groupSelect).
		Scopes(conds...).
		Group(groupColumns).
		Find(&res).Error
	if err != nil {
		return nil, wrapDBErr("group record from database failed", err)
	}
	return res, nil
}

// FindModelInBatches iterates rows in batches. res holds the current batch
// inside fc; do not rely on its contents after FindModelInBatches returns.
func FindModelInBatches[M any](ctx context.Context, db *gorm.DB, model *M, res *[]*M, batchSize int, fc func(tx *gorm.DB, batch int) error, scopes ...func(tx *gorm.DB) *gorm.DB) error {
	if db == nil {
		return invalidArgf("db is nil")
	}
	if model == nil {
		return invalidArgf("model is nil")
	}
	if res == nil {
		return invalidArgf("result slice is nil")
	}
	if batchSize <= 0 {
		return invalidArgf("batch size must be positive")
	}
	if fc == nil {
		return invalidArgf("batch callback is nil")
	}
	err := db.WithContext(ctx).
		Model(model).
		Scopes(scopes...).
		FindInBatches(res, batchSize, fc).Error
	if err != nil {
		return wrapDBErr("find record from database failed", err)
	}
	return nil
}
