// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// Vin is an object representing the database table.
type Vin struct {
	Vin                 string        `boil:"vin" json:"vin" toml:"vin" yaml:"vin"`
	FirstEarningWeek    int           `boil:"first_earning_week" json:"first_earning_week" toml:"first_earning_week" yaml:"first_earning_week"`
	FirstEarningTokenID types.Decimal `boil:"first_earning_token_id" json:"first_earning_token_id" toml:"first_earning_token_id" yaml:"first_earning_token_id"`

	R *vinR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L vinL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var VinColumns = struct {
	Vin                 string
	FirstEarningWeek    string
	FirstEarningTokenID string
}{
	Vin:                 "vin",
	FirstEarningWeek:    "first_earning_week",
	FirstEarningTokenID: "first_earning_token_id",
}

var VinTableColumns = struct {
	Vin                 string
	FirstEarningWeek    string
	FirstEarningTokenID string
}{
	Vin:                 "vins.vin",
	FirstEarningWeek:    "vins.first_earning_week",
	FirstEarningTokenID: "vins.first_earning_token_id",
}

// Generated where

var VinWhere = struct {
	Vin                 whereHelperstring
	FirstEarningWeek    whereHelperint
	FirstEarningTokenID whereHelpertypes_Decimal
}{
	Vin:                 whereHelperstring{field: "\"rewards_api\".\"vins\".\"vin\""},
	FirstEarningWeek:    whereHelperint{field: "\"rewards_api\".\"vins\".\"first_earning_week\""},
	FirstEarningTokenID: whereHelpertypes_Decimal{field: "\"rewards_api\".\"vins\".\"first_earning_token_id\""},
}

// VinRels is where relationship names are stored.
var VinRels = struct {
}{}

// vinR is where relationships are stored.
type vinR struct {
}

// NewStruct creates a new relationship struct
func (*vinR) NewStruct() *vinR {
	return &vinR{}
}

// vinL is where Load methods for each relationship are stored.
type vinL struct{}

var (
	vinAllColumns            = []string{"vin", "first_earning_week", "first_earning_token_id"}
	vinColumnsWithoutDefault = []string{"vin", "first_earning_week", "first_earning_token_id"}
	vinColumnsWithDefault    = []string{}
	vinPrimaryKeyColumns     = []string{"vin"}
	vinGeneratedColumns      = []string{}
)

type (
	// VinSlice is an alias for a slice of pointers to Vin.
	// This should almost always be used instead of []Vin.
	VinSlice []*Vin
	// VinHook is the signature for custom Vin hook methods
	VinHook func(context.Context, boil.ContextExecutor, *Vin) error

	vinQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	vinType                 = reflect.TypeOf(&Vin{})
	vinMapping              = queries.MakeStructMapping(vinType)
	vinPrimaryKeyMapping, _ = queries.BindMapping(vinType, vinMapping, vinPrimaryKeyColumns)
	vinInsertCacheMut       sync.RWMutex
	vinInsertCache          = make(map[string]insertCache)
	vinUpdateCacheMut       sync.RWMutex
	vinUpdateCache          = make(map[string]updateCache)
	vinUpsertCacheMut       sync.RWMutex
	vinUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var vinAfterSelectHooks []VinHook

var vinBeforeInsertHooks []VinHook
var vinAfterInsertHooks []VinHook

var vinBeforeUpdateHooks []VinHook
var vinAfterUpdateHooks []VinHook

var vinBeforeDeleteHooks []VinHook
var vinAfterDeleteHooks []VinHook

var vinBeforeUpsertHooks []VinHook
var vinAfterUpsertHooks []VinHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Vin) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Vin) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Vin) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Vin) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Vin) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Vin) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Vin) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Vin) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Vin) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vinAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddVinHook registers your hook function for all future operations.
func AddVinHook(hookPoint boil.HookPoint, vinHook VinHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		vinAfterSelectHooks = append(vinAfterSelectHooks, vinHook)
	case boil.BeforeInsertHook:
		vinBeforeInsertHooks = append(vinBeforeInsertHooks, vinHook)
	case boil.AfterInsertHook:
		vinAfterInsertHooks = append(vinAfterInsertHooks, vinHook)
	case boil.BeforeUpdateHook:
		vinBeforeUpdateHooks = append(vinBeforeUpdateHooks, vinHook)
	case boil.AfterUpdateHook:
		vinAfterUpdateHooks = append(vinAfterUpdateHooks, vinHook)
	case boil.BeforeDeleteHook:
		vinBeforeDeleteHooks = append(vinBeforeDeleteHooks, vinHook)
	case boil.AfterDeleteHook:
		vinAfterDeleteHooks = append(vinAfterDeleteHooks, vinHook)
	case boil.BeforeUpsertHook:
		vinBeforeUpsertHooks = append(vinBeforeUpsertHooks, vinHook)
	case boil.AfterUpsertHook:
		vinAfterUpsertHooks = append(vinAfterUpsertHooks, vinHook)
	}
}

// One returns a single vin record from the query.
func (q vinQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Vin, error) {
	o := &Vin{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for vins")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Vin records from the query.
func (q vinQuery) All(ctx context.Context, exec boil.ContextExecutor) (VinSlice, error) {
	var o []*Vin

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Vin slice")
	}

	if len(vinAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Vin records in the query.
func (q vinQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count vins rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q vinQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if vins exists")
	}

	return count > 0, nil
}

// Vins retrieves all the records using an executor.
func Vins(mods ...qm.QueryMod) vinQuery {
	mods = append(mods, qm.From("\"rewards_api\".\"vins\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"rewards_api\".\"vins\".*"})
	}

	return vinQuery{q}
}

// FindVin retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindVin(ctx context.Context, exec boil.ContextExecutor, vin string, selectCols ...string) (*Vin, error) {
	vinObj := &Vin{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"rewards_api\".\"vins\" where \"vin\"=$1", sel,
	)

	q := queries.Raw(query, vin)

	err := q.Bind(ctx, exec, vinObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from vins")
	}

	if err = vinObj.doAfterSelectHooks(ctx, exec); err != nil {
		return vinObj, err
	}

	return vinObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Vin) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no vins provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(vinColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	vinInsertCacheMut.RLock()
	cache, cached := vinInsertCache[key]
	vinInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			vinAllColumns,
			vinColumnsWithDefault,
			vinColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(vinType, vinMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(vinType, vinMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"rewards_api\".\"vins\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"rewards_api\".\"vins\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into vins")
	}

	if !cached {
		vinInsertCacheMut.Lock()
		vinInsertCache[key] = cache
		vinInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Vin.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Vin) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	vinUpdateCacheMut.RLock()
	cache, cached := vinUpdateCache[key]
	vinUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			vinAllColumns,
			vinPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update vins, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"rewards_api\".\"vins\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, vinPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(vinType, vinMapping, append(wl, vinPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update vins row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for vins")
	}

	if !cached {
		vinUpdateCacheMut.Lock()
		vinUpdateCache[key] = cache
		vinUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q vinQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for vins")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for vins")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o VinSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vinPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"rewards_api\".\"vins\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, vinPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in vin slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all vin")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Vin) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no vins provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(vinColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	vinUpsertCacheMut.RLock()
	cache, cached := vinUpsertCache[key]
	vinUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			vinAllColumns,
			vinColumnsWithDefault,
			vinColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			vinAllColumns,
			vinPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert vins, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(vinPrimaryKeyColumns))
			copy(conflict, vinPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"rewards_api\".\"vins\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(vinType, vinMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(vinType, vinMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert vins")
	}

	if !cached {
		vinUpsertCacheMut.Lock()
		vinUpsertCache[key] = cache
		vinUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Vin record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Vin) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Vin provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), vinPrimaryKeyMapping)
	sql := "DELETE FROM \"rewards_api\".\"vins\" WHERE \"vin\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from vins")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for vins")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q vinQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no vinQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from vins")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for vins")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o VinSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(vinBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vinPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"rewards_api\".\"vins\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vinPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from vin slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for vins")
	}

	if len(vinAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Vin) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindVin(ctx, exec, o.Vin)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *VinSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := VinSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vinPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"rewards_api\".\"vins\".* FROM \"rewards_api\".\"vins\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vinPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in VinSlice")
	}

	*o = slice

	return nil
}

// VinExists checks if the Vin row exists.
func VinExists(ctx context.Context, exec boil.ContextExecutor, vin string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"rewards_api\".\"vins\" where \"vin\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, vin)
	}
	row := exec.QueryRowContext(ctx, sql, vin)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if vins exists")
	}

	return exists, nil
}

// Exists checks if the Vin row exists.
func (o *Vin) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return VinExists(ctx, exec, o.Vin)
}
