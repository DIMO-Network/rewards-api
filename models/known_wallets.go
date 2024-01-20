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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// KnownWallet is an object representing the database table.
type KnownWallet struct {
	ChainID     int64       `boil:"chain_id" json:"chain_id" toml:"chain_id" yaml:"chain_id"`
	Address     []byte      `boil:"address" json:"address" toml:"address" yaml:"address"`
	Description string      `boil:"description" json:"description" toml:"description" yaml:"description"`
	Type        null.String `boil:"type" json:"type,omitempty" toml:"type" yaml:"type,omitempty"`

	R *knownWalletR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L knownWalletL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var KnownWalletColumns = struct {
	ChainID     string
	Address     string
	Description string
	Type        string
}{
	ChainID:     "chain_id",
	Address:     "address",
	Description: "description",
	Type:        "type",
}

var KnownWalletTableColumns = struct {
	ChainID     string
	Address     string
	Description string
	Type        string
}{
	ChainID:     "known_wallets.chain_id",
	Address:     "known_wallets.address",
	Description: "known_wallets.description",
	Type:        "known_wallets.type",
}

// Generated where

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelper__byte struct{ field string }

func (w whereHelper__byte) EQ(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelper__byte) NEQ(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelper__byte) LT(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelper__byte) LTE(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelper__byte) GT(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelper__byte) GTE(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelpernull_String struct{ field string }

func (w whereHelpernull_String) EQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_String) NEQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_String) LT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_String) LTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_String) GT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_String) GTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelpernull_String) LIKE(x null.String) qm.QueryMod {
	return qm.Where(w.field+" LIKE ?", x)
}
func (w whereHelpernull_String) NLIKE(x null.String) qm.QueryMod {
	return qm.Where(w.field+" NOT LIKE ?", x)
}
func (w whereHelpernull_String) ILIKE(x null.String) qm.QueryMod {
	return qm.Where(w.field+" ILIKE ?", x)
}
func (w whereHelpernull_String) NILIKE(x null.String) qm.QueryMod {
	return qm.Where(w.field+" NOT ILIKE ?", x)
}
func (w whereHelpernull_String) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelpernull_String) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

func (w whereHelpernull_String) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_String) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var KnownWalletWhere = struct {
	ChainID     whereHelperint64
	Address     whereHelper__byte
	Description whereHelperstring
	Type        whereHelpernull_String
}{
	ChainID:     whereHelperint64{field: "\"rewards_api\".\"known_wallets\".\"chain_id\""},
	Address:     whereHelper__byte{field: "\"rewards_api\".\"known_wallets\".\"address\""},
	Description: whereHelperstring{field: "\"rewards_api\".\"known_wallets\".\"description\""},
	Type:        whereHelpernull_String{field: "\"rewards_api\".\"known_wallets\".\"type\""},
}

// KnownWalletRels is where relationship names are stored.
var KnownWalletRels = struct {
}{}

// knownWalletR is where relationships are stored.
type knownWalletR struct {
}

// NewStruct creates a new relationship struct
func (*knownWalletR) NewStruct() *knownWalletR {
	return &knownWalletR{}
}

// knownWalletL is where Load methods for each relationship are stored.
type knownWalletL struct{}

var (
	knownWalletAllColumns            = []string{"chain_id", "address", "description", "type"}
	knownWalletColumnsWithoutDefault = []string{"chain_id", "address", "description"}
	knownWalletColumnsWithDefault    = []string{"type"}
	knownWalletPrimaryKeyColumns     = []string{"chain_id", "address"}
	knownWalletGeneratedColumns      = []string{}
)

type (
	// KnownWalletSlice is an alias for a slice of pointers to KnownWallet.
	// This should almost always be used instead of []KnownWallet.
	KnownWalletSlice []*KnownWallet
	// KnownWalletHook is the signature for custom KnownWallet hook methods
	KnownWalletHook func(context.Context, boil.ContextExecutor, *KnownWallet) error

	knownWalletQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	knownWalletType                 = reflect.TypeOf(&KnownWallet{})
	knownWalletMapping              = queries.MakeStructMapping(knownWalletType)
	knownWalletPrimaryKeyMapping, _ = queries.BindMapping(knownWalletType, knownWalletMapping, knownWalletPrimaryKeyColumns)
	knownWalletInsertCacheMut       sync.RWMutex
	knownWalletInsertCache          = make(map[string]insertCache)
	knownWalletUpdateCacheMut       sync.RWMutex
	knownWalletUpdateCache          = make(map[string]updateCache)
	knownWalletUpsertCacheMut       sync.RWMutex
	knownWalletUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var knownWalletAfterSelectHooks []KnownWalletHook

var knownWalletBeforeInsertHooks []KnownWalletHook
var knownWalletAfterInsertHooks []KnownWalletHook

var knownWalletBeforeUpdateHooks []KnownWalletHook
var knownWalletAfterUpdateHooks []KnownWalletHook

var knownWalletBeforeDeleteHooks []KnownWalletHook
var knownWalletAfterDeleteHooks []KnownWalletHook

var knownWalletBeforeUpsertHooks []KnownWalletHook
var knownWalletAfterUpsertHooks []KnownWalletHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *KnownWallet) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *KnownWallet) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *KnownWallet) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *KnownWallet) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *KnownWallet) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *KnownWallet) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *KnownWallet) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *KnownWallet) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *KnownWallet) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range knownWalletAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddKnownWalletHook registers your hook function for all future operations.
func AddKnownWalletHook(hookPoint boil.HookPoint, knownWalletHook KnownWalletHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		knownWalletAfterSelectHooks = append(knownWalletAfterSelectHooks, knownWalletHook)
	case boil.BeforeInsertHook:
		knownWalletBeforeInsertHooks = append(knownWalletBeforeInsertHooks, knownWalletHook)
	case boil.AfterInsertHook:
		knownWalletAfterInsertHooks = append(knownWalletAfterInsertHooks, knownWalletHook)
	case boil.BeforeUpdateHook:
		knownWalletBeforeUpdateHooks = append(knownWalletBeforeUpdateHooks, knownWalletHook)
	case boil.AfterUpdateHook:
		knownWalletAfterUpdateHooks = append(knownWalletAfterUpdateHooks, knownWalletHook)
	case boil.BeforeDeleteHook:
		knownWalletBeforeDeleteHooks = append(knownWalletBeforeDeleteHooks, knownWalletHook)
	case boil.AfterDeleteHook:
		knownWalletAfterDeleteHooks = append(knownWalletAfterDeleteHooks, knownWalletHook)
	case boil.BeforeUpsertHook:
		knownWalletBeforeUpsertHooks = append(knownWalletBeforeUpsertHooks, knownWalletHook)
	case boil.AfterUpsertHook:
		knownWalletAfterUpsertHooks = append(knownWalletAfterUpsertHooks, knownWalletHook)
	}
}

// One returns a single knownWallet record from the query.
func (q knownWalletQuery) One(ctx context.Context, exec boil.ContextExecutor) (*KnownWallet, error) {
	o := &KnownWallet{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for known_wallets")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all KnownWallet records from the query.
func (q knownWalletQuery) All(ctx context.Context, exec boil.ContextExecutor) (KnownWalletSlice, error) {
	var o []*KnownWallet

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to KnownWallet slice")
	}

	if len(knownWalletAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all KnownWallet records in the query.
func (q knownWalletQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count known_wallets rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q knownWalletQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if known_wallets exists")
	}

	return count > 0, nil
}

// KnownWallets retrieves all the records using an executor.
func KnownWallets(mods ...qm.QueryMod) knownWalletQuery {
	mods = append(mods, qm.From("\"rewards_api\".\"known_wallets\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"rewards_api\".\"known_wallets\".*"})
	}

	return knownWalletQuery{q}
}

// FindKnownWallet retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindKnownWallet(ctx context.Context, exec boil.ContextExecutor, chainID int64, address []byte, selectCols ...string) (*KnownWallet, error) {
	knownWalletObj := &KnownWallet{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"rewards_api\".\"known_wallets\" where \"chain_id\"=$1 AND \"address\"=$2", sel,
	)

	q := queries.Raw(query, chainID, address)

	err := q.Bind(ctx, exec, knownWalletObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from known_wallets")
	}

	if err = knownWalletObj.doAfterSelectHooks(ctx, exec); err != nil {
		return knownWalletObj, err
	}

	return knownWalletObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *KnownWallet) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no known_wallets provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(knownWalletColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	knownWalletInsertCacheMut.RLock()
	cache, cached := knownWalletInsertCache[key]
	knownWalletInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			knownWalletAllColumns,
			knownWalletColumnsWithDefault,
			knownWalletColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(knownWalletType, knownWalletMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(knownWalletType, knownWalletMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"rewards_api\".\"known_wallets\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"rewards_api\".\"known_wallets\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into known_wallets")
	}

	if !cached {
		knownWalletInsertCacheMut.Lock()
		knownWalletInsertCache[key] = cache
		knownWalletInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the KnownWallet.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *KnownWallet) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	knownWalletUpdateCacheMut.RLock()
	cache, cached := knownWalletUpdateCache[key]
	knownWalletUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			knownWalletAllColumns,
			knownWalletPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update known_wallets, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"rewards_api\".\"known_wallets\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, knownWalletPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(knownWalletType, knownWalletMapping, append(wl, knownWalletPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update known_wallets row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for known_wallets")
	}

	if !cached {
		knownWalletUpdateCacheMut.Lock()
		knownWalletUpdateCache[key] = cache
		knownWalletUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q knownWalletQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for known_wallets")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for known_wallets")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o KnownWalletSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), knownWalletPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"rewards_api\".\"known_wallets\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, knownWalletPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in knownWallet slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all knownWallet")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *KnownWallet) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no known_wallets provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(knownWalletColumnsWithDefault, o)

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

	knownWalletUpsertCacheMut.RLock()
	cache, cached := knownWalletUpsertCache[key]
	knownWalletUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			knownWalletAllColumns,
			knownWalletColumnsWithDefault,
			knownWalletColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			knownWalletAllColumns,
			knownWalletPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert known_wallets, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(knownWalletPrimaryKeyColumns))
			copy(conflict, knownWalletPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"rewards_api\".\"known_wallets\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(knownWalletType, knownWalletMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(knownWalletType, knownWalletMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert known_wallets")
	}

	if !cached {
		knownWalletUpsertCacheMut.Lock()
		knownWalletUpsertCache[key] = cache
		knownWalletUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single KnownWallet record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *KnownWallet) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no KnownWallet provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), knownWalletPrimaryKeyMapping)
	sql := "DELETE FROM \"rewards_api\".\"known_wallets\" WHERE \"chain_id\"=$1 AND \"address\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from known_wallets")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for known_wallets")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q knownWalletQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no knownWalletQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from known_wallets")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for known_wallets")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o KnownWalletSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(knownWalletBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), knownWalletPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"rewards_api\".\"known_wallets\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, knownWalletPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from knownWallet slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for known_wallets")
	}

	if len(knownWalletAfterDeleteHooks) != 0 {
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
func (o *KnownWallet) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindKnownWallet(ctx, exec, o.ChainID, o.Address)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *KnownWalletSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := KnownWalletSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), knownWalletPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"rewards_api\".\"known_wallets\".* FROM \"rewards_api\".\"known_wallets\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, knownWalletPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in KnownWalletSlice")
	}

	*o = slice

	return nil
}

// KnownWalletExists checks if the KnownWallet row exists.
func KnownWalletExists(ctx context.Context, exec boil.ContextExecutor, chainID int64, address []byte) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"rewards_api\".\"known_wallets\" where \"chain_id\"=$1 AND \"address\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, chainID, address)
	}
	row := exec.QueryRowContext(ctx, sql, chainID, address)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if known_wallets exists")
	}

	return exists, nil
}

// Exists checks if the KnownWallet row exists.
func (o *KnownWallet) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return KnownWalletExists(ctx, exec, o.ChainID, o.Address)
}
