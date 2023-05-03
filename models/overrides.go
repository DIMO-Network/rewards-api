// Code generated by SQLBoiler 4.13.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/strmangle"
)

// Override is an object representing the database table.
type Override struct {
	IssuanceWeekID   int    `boil:"issuance_week_id" json:"issuance_week_id" toml:"issuance_week_id" yaml:"issuance_week_id"`
	UserDeviceID     string `boil:"user_device_id" json:"user_device_id" toml:"user_device_id" yaml:"user_device_id"`
	ConnectionStreak int    `boil:"connection_streak" json:"connection_streak" toml:"connection_streak" yaml:"connection_streak"`
	Note             string `boil:"note" json:"note" toml:"note" yaml:"note"`

	R *overrideR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L overrideL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var OverrideColumns = struct {
	IssuanceWeekID   string
	UserDeviceID     string
	ConnectionStreak string
	Note             string
}{
	IssuanceWeekID:   "issuance_week_id",
	UserDeviceID:     "user_device_id",
	ConnectionStreak: "connection_streak",
	Note:             "note",
}

var OverrideTableColumns = struct {
	IssuanceWeekID   string
	UserDeviceID     string
	ConnectionStreak string
	Note             string
}{
	IssuanceWeekID:   "overrides.issuance_week_id",
	UserDeviceID:     "overrides.user_device_id",
	ConnectionStreak: "overrides.connection_streak",
	Note:             "overrides.note",
}

// Generated where

var OverrideWhere = struct {
	IssuanceWeekID   whereHelperint
	UserDeviceID     whereHelperstring
	ConnectionStreak whereHelperint
	Note             whereHelperstring
}{
	IssuanceWeekID:   whereHelperint{field: "\"rewards_api\".\"overrides\".\"issuance_week_id\""},
	UserDeviceID:     whereHelperstring{field: "\"rewards_api\".\"overrides\".\"user_device_id\""},
	ConnectionStreak: whereHelperint{field: "\"rewards_api\".\"overrides\".\"connection_streak\""},
	Note:             whereHelperstring{field: "\"rewards_api\".\"overrides\".\"note\""},
}

// OverrideRels is where relationship names are stored.
var OverrideRels = struct {
}{}

// overrideR is where relationships are stored.
type overrideR struct {
}

// NewStruct creates a new relationship struct
func (*overrideR) NewStruct() *overrideR {
	return &overrideR{}
}

// overrideL is where Load methods for each relationship are stored.
type overrideL struct{}

var (
	overrideAllColumns            = []string{"issuance_week_id", "user_device_id", "connection_streak", "note"}
	overrideColumnsWithoutDefault = []string{"issuance_week_id", "user_device_id", "connection_streak", "note"}
	overrideColumnsWithDefault    = []string{}
	overridePrimaryKeyColumns     = []string{"issuance_week_id", "user_device_id"}
	overrideGeneratedColumns      = []string{}
)

type (
	// OverrideSlice is an alias for a slice of pointers to Override.
	// This should almost always be used instead of []Override.
	OverrideSlice []*Override
	// OverrideHook is the signature for custom Override hook methods
	OverrideHook func(context.Context, boil.ContextExecutor, *Override) error

	overrideQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	overrideType                 = reflect.TypeOf(&Override{})
	overrideMapping              = queries.MakeStructMapping(overrideType)
	overridePrimaryKeyMapping, _ = queries.BindMapping(overrideType, overrideMapping, overridePrimaryKeyColumns)
	overrideInsertCacheMut       sync.RWMutex
	overrideInsertCache          = make(map[string]insertCache)
	overrideUpdateCacheMut       sync.RWMutex
	overrideUpdateCache          = make(map[string]updateCache)
	overrideUpsertCacheMut       sync.RWMutex
	overrideUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var overrideAfterSelectHooks []OverrideHook

var overrideBeforeInsertHooks []OverrideHook
var overrideAfterInsertHooks []OverrideHook

var overrideBeforeUpdateHooks []OverrideHook
var overrideAfterUpdateHooks []OverrideHook

var overrideBeforeDeleteHooks []OverrideHook
var overrideAfterDeleteHooks []OverrideHook

var overrideBeforeUpsertHooks []OverrideHook
var overrideAfterUpsertHooks []OverrideHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Override) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Override) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Override) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Override) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Override) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Override) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Override) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Override) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Override) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range overrideAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddOverrideHook registers your hook function for all future operations.
func AddOverrideHook(hookPoint boil.HookPoint, overrideHook OverrideHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		overrideAfterSelectHooks = append(overrideAfterSelectHooks, overrideHook)
	case boil.BeforeInsertHook:
		overrideBeforeInsertHooks = append(overrideBeforeInsertHooks, overrideHook)
	case boil.AfterInsertHook:
		overrideAfterInsertHooks = append(overrideAfterInsertHooks, overrideHook)
	case boil.BeforeUpdateHook:
		overrideBeforeUpdateHooks = append(overrideBeforeUpdateHooks, overrideHook)
	case boil.AfterUpdateHook:
		overrideAfterUpdateHooks = append(overrideAfterUpdateHooks, overrideHook)
	case boil.BeforeDeleteHook:
		overrideBeforeDeleteHooks = append(overrideBeforeDeleteHooks, overrideHook)
	case boil.AfterDeleteHook:
		overrideAfterDeleteHooks = append(overrideAfterDeleteHooks, overrideHook)
	case boil.BeforeUpsertHook:
		overrideBeforeUpsertHooks = append(overrideBeforeUpsertHooks, overrideHook)
	case boil.AfterUpsertHook:
		overrideAfterUpsertHooks = append(overrideAfterUpsertHooks, overrideHook)
	}
}

// One returns a single override record from the query.
func (q overrideQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Override, error) {
	o := &Override{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for overrides")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Override records from the query.
func (q overrideQuery) All(ctx context.Context, exec boil.ContextExecutor) (OverrideSlice, error) {
	var o []*Override

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Override slice")
	}

	if len(overrideAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Override records in the query.
func (q overrideQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count overrides rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q overrideQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if overrides exists")
	}

	return count > 0, nil
}

// Overrides retrieves all the records using an executor.
func Overrides(mods ...qm.QueryMod) overrideQuery {
	mods = append(mods, qm.From("\"rewards_api\".\"overrides\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"rewards_api\".\"overrides\".*"})
	}

	return overrideQuery{q}
}

// FindOverride retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindOverride(ctx context.Context, exec boil.ContextExecutor, issuanceWeekID int, userDeviceID string, selectCols ...string) (*Override, error) {
	overrideObj := &Override{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"rewards_api\".\"overrides\" where \"issuance_week_id\"=$1 AND \"user_device_id\"=$2", sel,
	)

	q := queries.Raw(query, issuanceWeekID, userDeviceID)

	err := q.Bind(ctx, exec, overrideObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from overrides")
	}

	if err = overrideObj.doAfterSelectHooks(ctx, exec); err != nil {
		return overrideObj, err
	}

	return overrideObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Override) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no overrides provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(overrideColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	overrideInsertCacheMut.RLock()
	cache, cached := overrideInsertCache[key]
	overrideInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			overrideAllColumns,
			overrideColumnsWithDefault,
			overrideColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(overrideType, overrideMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(overrideType, overrideMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"rewards_api\".\"overrides\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"rewards_api\".\"overrides\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into overrides")
	}

	if !cached {
		overrideInsertCacheMut.Lock()
		overrideInsertCache[key] = cache
		overrideInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Override.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Override) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	overrideUpdateCacheMut.RLock()
	cache, cached := overrideUpdateCache[key]
	overrideUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			overrideAllColumns,
			overridePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update overrides, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"rewards_api\".\"overrides\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, overridePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(overrideType, overrideMapping, append(wl, overridePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update overrides row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for overrides")
	}

	if !cached {
		overrideUpdateCacheMut.Lock()
		overrideUpdateCache[key] = cache
		overrideUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q overrideQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for overrides")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for overrides")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o OverrideSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), overridePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"rewards_api\".\"overrides\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, overridePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in override slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all override")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Override) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no overrides provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(overrideColumnsWithDefault, o)

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

	overrideUpsertCacheMut.RLock()
	cache, cached := overrideUpsertCache[key]
	overrideUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			overrideAllColumns,
			overrideColumnsWithDefault,
			overrideColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			overrideAllColumns,
			overridePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert overrides, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(overridePrimaryKeyColumns))
			copy(conflict, overridePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"rewards_api\".\"overrides\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(overrideType, overrideMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(overrideType, overrideMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert overrides")
	}

	if !cached {
		overrideUpsertCacheMut.Lock()
		overrideUpsertCache[key] = cache
		overrideUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Override record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Override) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Override provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), overridePrimaryKeyMapping)
	sql := "DELETE FROM \"rewards_api\".\"overrides\" WHERE \"issuance_week_id\"=$1 AND \"user_device_id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from overrides")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for overrides")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q overrideQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no overrideQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from overrides")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for overrides")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o OverrideSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(overrideBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), overridePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"rewards_api\".\"overrides\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, overridePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from override slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for overrides")
	}

	if len(overrideAfterDeleteHooks) != 0 {
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
func (o *Override) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindOverride(ctx, exec, o.IssuanceWeekID, o.UserDeviceID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OverrideSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := OverrideSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), overridePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"rewards_api\".\"overrides\".* FROM \"rewards_api\".\"overrides\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, overridePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in OverrideSlice")
	}

	*o = slice

	return nil
}

// OverrideExists checks if the Override row exists.
func OverrideExists(ctx context.Context, exec boil.ContextExecutor, issuanceWeekID int, userDeviceID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"rewards_api\".\"overrides\" where \"issuance_week_id\"=$1 AND \"user_device_id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, issuanceWeekID, userDeviceID)
	}
	row := exec.QueryRowContext(ctx, sql, issuanceWeekID, userDeviceID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if overrides exists")
	}

	return exists, nil
}
