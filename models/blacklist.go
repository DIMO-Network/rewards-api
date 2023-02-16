// Code generated by SQLBoiler 4.14.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// Blacklist is an object representing the database table.
type Blacklist struct {
	UserEthereumAddress string    `boil:"user_ethereum_address" json:"user_ethereum_address" toml:"user_ethereum_address" yaml:"user_ethereum_address"`
	CreatedAt           time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Note                string    `boil:"note" json:"note" toml:"note" yaml:"note"`

	R *blacklistR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L blacklistL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var BlacklistColumns = struct {
	UserEthereumAddress string
	CreatedAt           string
	Note                string
}{
	UserEthereumAddress: "user_ethereum_address",
	CreatedAt:           "created_at",
	Note:                "note",
}

var BlacklistTableColumns = struct {
	UserEthereumAddress string
	CreatedAt           string
	Note                string
}{
	UserEthereumAddress: "blacklist.user_ethereum_address",
	CreatedAt:           "blacklist.created_at",
	Note:                "blacklist.note",
}

// Generated where

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var BlacklistWhere = struct {
	UserEthereumAddress whereHelperstring
	CreatedAt           whereHelpertime_Time
	Note                whereHelperstring
}{
	UserEthereumAddress: whereHelperstring{field: "\"rewards_api\".\"blacklist\".\"user_ethereum_address\""},
	CreatedAt:           whereHelpertime_Time{field: "\"rewards_api\".\"blacklist\".\"created_at\""},
	Note:                whereHelperstring{field: "\"rewards_api\".\"blacklist\".\"note\""},
}

// BlacklistRels is where relationship names are stored.
var BlacklistRels = struct {
}{}

// blacklistR is where relationships are stored.
type blacklistR struct {
}

// NewStruct creates a new relationship struct
func (*blacklistR) NewStruct() *blacklistR {
	return &blacklistR{}
}

// blacklistL is where Load methods for each relationship are stored.
type blacklistL struct{}

var (
	blacklistAllColumns            = []string{"user_ethereum_address", "created_at", "note"}
	blacklistColumnsWithoutDefault = []string{"user_ethereum_address", "note"}
	blacklistColumnsWithDefault    = []string{"created_at"}
	blacklistPrimaryKeyColumns     = []string{"user_ethereum_address"}
	blacklistGeneratedColumns      = []string{}
)

type (
	// BlacklistSlice is an alias for a slice of pointers to Blacklist.
	// This should almost always be used instead of []Blacklist.
	BlacklistSlice []*Blacklist
	// BlacklistHook is the signature for custom Blacklist hook methods
	BlacklistHook func(context.Context, boil.ContextExecutor, *Blacklist) error

	blacklistQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	blacklistType                 = reflect.TypeOf(&Blacklist{})
	blacklistMapping              = queries.MakeStructMapping(blacklistType)
	blacklistPrimaryKeyMapping, _ = queries.BindMapping(blacklistType, blacklistMapping, blacklistPrimaryKeyColumns)
	blacklistInsertCacheMut       sync.RWMutex
	blacklistInsertCache          = make(map[string]insertCache)
	blacklistUpdateCacheMut       sync.RWMutex
	blacklistUpdateCache          = make(map[string]updateCache)
	blacklistUpsertCacheMut       sync.RWMutex
	blacklistUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var blacklistAfterSelectHooks []BlacklistHook

var blacklistBeforeInsertHooks []BlacklistHook
var blacklistAfterInsertHooks []BlacklistHook

var blacklistBeforeUpdateHooks []BlacklistHook
var blacklistAfterUpdateHooks []BlacklistHook

var blacklistBeforeDeleteHooks []BlacklistHook
var blacklistAfterDeleteHooks []BlacklistHook

var blacklistBeforeUpsertHooks []BlacklistHook
var blacklistAfterUpsertHooks []BlacklistHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Blacklist) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Blacklist) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Blacklist) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Blacklist) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Blacklist) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Blacklist) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Blacklist) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Blacklist) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Blacklist) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blacklistAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddBlacklistHook registers your hook function for all future operations.
func AddBlacklistHook(hookPoint boil.HookPoint, blacklistHook BlacklistHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		blacklistAfterSelectHooks = append(blacklistAfterSelectHooks, blacklistHook)
	case boil.BeforeInsertHook:
		blacklistBeforeInsertHooks = append(blacklistBeforeInsertHooks, blacklistHook)
	case boil.AfterInsertHook:
		blacklistAfterInsertHooks = append(blacklistAfterInsertHooks, blacklistHook)
	case boil.BeforeUpdateHook:
		blacklistBeforeUpdateHooks = append(blacklistBeforeUpdateHooks, blacklistHook)
	case boil.AfterUpdateHook:
		blacklistAfterUpdateHooks = append(blacklistAfterUpdateHooks, blacklistHook)
	case boil.BeforeDeleteHook:
		blacklistBeforeDeleteHooks = append(blacklistBeforeDeleteHooks, blacklistHook)
	case boil.AfterDeleteHook:
		blacklistAfterDeleteHooks = append(blacklistAfterDeleteHooks, blacklistHook)
	case boil.BeforeUpsertHook:
		blacklistBeforeUpsertHooks = append(blacklistBeforeUpsertHooks, blacklistHook)
	case boil.AfterUpsertHook:
		blacklistAfterUpsertHooks = append(blacklistAfterUpsertHooks, blacklistHook)
	}
}

// One returns a single blacklist record from the query.
func (q blacklistQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Blacklist, error) {
	o := &Blacklist{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for blacklist")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Blacklist records from the query.
func (q blacklistQuery) All(ctx context.Context, exec boil.ContextExecutor) (BlacklistSlice, error) {
	var o []*Blacklist

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Blacklist slice")
	}

	if len(blacklistAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Blacklist records in the query.
func (q blacklistQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count blacklist rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q blacklistQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if blacklist exists")
	}

	return count > 0, nil
}

// Blacklists retrieves all the records using an executor.
func Blacklists(mods ...qm.QueryMod) blacklistQuery {
	mods = append(mods, qm.From("\"rewards_api\".\"blacklist\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"rewards_api\".\"blacklist\".*"})
	}

	return blacklistQuery{q}
}

// FindBlacklist retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindBlacklist(ctx context.Context, exec boil.ContextExecutor, userEthereumAddress string, selectCols ...string) (*Blacklist, error) {
	blacklistObj := &Blacklist{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"rewards_api\".\"blacklist\" where \"user_ethereum_address\"=$1", sel,
	)

	q := queries.Raw(query, userEthereumAddress)

	err := q.Bind(ctx, exec, blacklistObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from blacklist")
	}

	if err = blacklistObj.doAfterSelectHooks(ctx, exec); err != nil {
		return blacklistObj, err
	}

	return blacklistObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Blacklist) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no blacklist provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(blacklistColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	blacklistInsertCacheMut.RLock()
	cache, cached := blacklistInsertCache[key]
	blacklistInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			blacklistAllColumns,
			blacklistColumnsWithDefault,
			blacklistColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(blacklistType, blacklistMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(blacklistType, blacklistMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"rewards_api\".\"blacklist\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"rewards_api\".\"blacklist\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into blacklist")
	}

	if !cached {
		blacklistInsertCacheMut.Lock()
		blacklistInsertCache[key] = cache
		blacklistInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Blacklist.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Blacklist) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	blacklistUpdateCacheMut.RLock()
	cache, cached := blacklistUpdateCache[key]
	blacklistUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			blacklistAllColumns,
			blacklistPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update blacklist, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"rewards_api\".\"blacklist\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, blacklistPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(blacklistType, blacklistMapping, append(wl, blacklistPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update blacklist row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for blacklist")
	}

	if !cached {
		blacklistUpdateCacheMut.Lock()
		blacklistUpdateCache[key] = cache
		blacklistUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q blacklistQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for blacklist")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for blacklist")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o BlacklistSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blacklistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"rewards_api\".\"blacklist\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, blacklistPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in blacklist slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all blacklist")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Blacklist) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no blacklist provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(blacklistColumnsWithDefault, o)

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

	blacklistUpsertCacheMut.RLock()
	cache, cached := blacklistUpsertCache[key]
	blacklistUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			blacklistAllColumns,
			blacklistColumnsWithDefault,
			blacklistColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			blacklistAllColumns,
			blacklistPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert blacklist, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(blacklistPrimaryKeyColumns))
			copy(conflict, blacklistPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"rewards_api\".\"blacklist\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(blacklistType, blacklistMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(blacklistType, blacklistMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert blacklist")
	}

	if !cached {
		blacklistUpsertCacheMut.Lock()
		blacklistUpsertCache[key] = cache
		blacklistUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Blacklist record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Blacklist) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Blacklist provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), blacklistPrimaryKeyMapping)
	sql := "DELETE FROM \"rewards_api\".\"blacklist\" WHERE \"user_ethereum_address\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from blacklist")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for blacklist")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q blacklistQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no blacklistQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from blacklist")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for blacklist")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o BlacklistSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(blacklistBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blacklistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"rewards_api\".\"blacklist\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, blacklistPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from blacklist slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for blacklist")
	}

	if len(blacklistAfterDeleteHooks) != 0 {
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
func (o *Blacklist) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindBlacklist(ctx, exec, o.UserEthereumAddress)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BlacklistSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := BlacklistSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blacklistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"rewards_api\".\"blacklist\".* FROM \"rewards_api\".\"blacklist\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, blacklistPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in BlacklistSlice")
	}

	*o = slice

	return nil
}

// BlacklistExists checks if the Blacklist row exists.
func BlacklistExists(ctx context.Context, exec boil.ContextExecutor, userEthereumAddress string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"rewards_api\".\"blacklist\" where \"user_ethereum_address\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, userEthereumAddress)
	}
	row := exec.QueryRowContext(ctx, sql, userEthereumAddress)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if blacklist exists")
	}

	return exists, nil
}

// Exists checks if the Blacklist row exists.
func (o *Blacklist) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return BlacklistExists(ctx, exec, o.UserEthereumAddress)
}