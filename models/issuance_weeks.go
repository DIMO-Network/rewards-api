// Code generated by SQLBoiler 4.8.6 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// IssuanceWeek is an object representing the database table.
type IssuanceWeek struct {
	ID        int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	JobStatus string    `boil:"job_status" json:"job_status" toml:"job_status" yaml:"job_status"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	StartsAt  time.Time `boil:"starts_at" json:"starts_at" toml:"starts_at" yaml:"starts_at"`
	EndsAt    time.Time `boil:"ends_at" json:"ends_at" toml:"ends_at" yaml:"ends_at"`

	R *issuanceWeekR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L issuanceWeekL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var IssuanceWeekColumns = struct {
	ID        string
	JobStatus string
	CreatedAt string
	UpdatedAt string
	StartsAt  string
	EndsAt    string
}{
	ID:        "id",
	JobStatus: "job_status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	StartsAt:  "starts_at",
	EndsAt:    "ends_at",
}

var IssuanceWeekTableColumns = struct {
	ID        string
	JobStatus string
	CreatedAt string
	UpdatedAt string
	StartsAt  string
	EndsAt    string
}{
	ID:        "issuance_weeks.id",
	JobStatus: "issuance_weeks.job_status",
	CreatedAt: "issuance_weeks.created_at",
	UpdatedAt: "issuance_weeks.updated_at",
	StartsAt:  "issuance_weeks.starts_at",
	EndsAt:    "issuance_weeks.ends_at",
}

// Generated where

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

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

var IssuanceWeekWhere = struct {
	ID        whereHelperint
	JobStatus whereHelperstring
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
	StartsAt  whereHelpertime_Time
	EndsAt    whereHelpertime_Time
}{
	ID:        whereHelperint{field: "\"rewards_api\".\"issuance_weeks\".\"id\""},
	JobStatus: whereHelperstring{field: "\"rewards_api\".\"issuance_weeks\".\"job_status\""},
	CreatedAt: whereHelpertime_Time{field: "\"rewards_api\".\"issuance_weeks\".\"created_at\""},
	UpdatedAt: whereHelpertime_Time{field: "\"rewards_api\".\"issuance_weeks\".\"updated_at\""},
	StartsAt:  whereHelpertime_Time{field: "\"rewards_api\".\"issuance_weeks\".\"starts_at\""},
	EndsAt:    whereHelpertime_Time{field: "\"rewards_api\".\"issuance_weeks\".\"ends_at\""},
}

// IssuanceWeekRels is where relationship names are stored.
var IssuanceWeekRels = struct {
	Rewards string
}{
	Rewards: "Rewards",
}

// issuanceWeekR is where relationships are stored.
type issuanceWeekR struct {
	Rewards RewardSlice `boil:"Rewards" json:"Rewards" toml:"Rewards" yaml:"Rewards"`
}

// NewStruct creates a new relationship struct
func (*issuanceWeekR) NewStruct() *issuanceWeekR {
	return &issuanceWeekR{}
}

// issuanceWeekL is where Load methods for each relationship are stored.
type issuanceWeekL struct{}

var (
	issuanceWeekAllColumns            = []string{"id", "job_status", "created_at", "updated_at", "starts_at", "ends_at"}
	issuanceWeekColumnsWithoutDefault = []string{"id", "job_status", "starts_at", "ends_at"}
	issuanceWeekColumnsWithDefault    = []string{"created_at", "updated_at"}
	issuanceWeekPrimaryKeyColumns     = []string{"id"}
	issuanceWeekGeneratedColumns      = []string{}
)

type (
	// IssuanceWeekSlice is an alias for a slice of pointers to IssuanceWeek.
	// This should almost always be used instead of []IssuanceWeek.
	IssuanceWeekSlice []*IssuanceWeek
	// IssuanceWeekHook is the signature for custom IssuanceWeek hook methods
	IssuanceWeekHook func(context.Context, boil.ContextExecutor, *IssuanceWeek) error

	issuanceWeekQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	issuanceWeekType                 = reflect.TypeOf(&IssuanceWeek{})
	issuanceWeekMapping              = queries.MakeStructMapping(issuanceWeekType)
	issuanceWeekPrimaryKeyMapping, _ = queries.BindMapping(issuanceWeekType, issuanceWeekMapping, issuanceWeekPrimaryKeyColumns)
	issuanceWeekInsertCacheMut       sync.RWMutex
	issuanceWeekInsertCache          = make(map[string]insertCache)
	issuanceWeekUpdateCacheMut       sync.RWMutex
	issuanceWeekUpdateCache          = make(map[string]updateCache)
	issuanceWeekUpsertCacheMut       sync.RWMutex
	issuanceWeekUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var issuanceWeekAfterSelectHooks []IssuanceWeekHook

var issuanceWeekBeforeInsertHooks []IssuanceWeekHook
var issuanceWeekAfterInsertHooks []IssuanceWeekHook

var issuanceWeekBeforeUpdateHooks []IssuanceWeekHook
var issuanceWeekAfterUpdateHooks []IssuanceWeekHook

var issuanceWeekBeforeDeleteHooks []IssuanceWeekHook
var issuanceWeekAfterDeleteHooks []IssuanceWeekHook

var issuanceWeekBeforeUpsertHooks []IssuanceWeekHook
var issuanceWeekAfterUpsertHooks []IssuanceWeekHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *IssuanceWeek) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *IssuanceWeek) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *IssuanceWeek) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *IssuanceWeek) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *IssuanceWeek) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *IssuanceWeek) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *IssuanceWeek) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *IssuanceWeek) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *IssuanceWeek) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range issuanceWeekAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddIssuanceWeekHook registers your hook function for all future operations.
func AddIssuanceWeekHook(hookPoint boil.HookPoint, issuanceWeekHook IssuanceWeekHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		issuanceWeekAfterSelectHooks = append(issuanceWeekAfterSelectHooks, issuanceWeekHook)
	case boil.BeforeInsertHook:
		issuanceWeekBeforeInsertHooks = append(issuanceWeekBeforeInsertHooks, issuanceWeekHook)
	case boil.AfterInsertHook:
		issuanceWeekAfterInsertHooks = append(issuanceWeekAfterInsertHooks, issuanceWeekHook)
	case boil.BeforeUpdateHook:
		issuanceWeekBeforeUpdateHooks = append(issuanceWeekBeforeUpdateHooks, issuanceWeekHook)
	case boil.AfterUpdateHook:
		issuanceWeekAfterUpdateHooks = append(issuanceWeekAfterUpdateHooks, issuanceWeekHook)
	case boil.BeforeDeleteHook:
		issuanceWeekBeforeDeleteHooks = append(issuanceWeekBeforeDeleteHooks, issuanceWeekHook)
	case boil.AfterDeleteHook:
		issuanceWeekAfterDeleteHooks = append(issuanceWeekAfterDeleteHooks, issuanceWeekHook)
	case boil.BeforeUpsertHook:
		issuanceWeekBeforeUpsertHooks = append(issuanceWeekBeforeUpsertHooks, issuanceWeekHook)
	case boil.AfterUpsertHook:
		issuanceWeekAfterUpsertHooks = append(issuanceWeekAfterUpsertHooks, issuanceWeekHook)
	}
}

// One returns a single issuanceWeek record from the query.
func (q issuanceWeekQuery) One(ctx context.Context, exec boil.ContextExecutor) (*IssuanceWeek, error) {
	o := &IssuanceWeek{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for issuance_weeks")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all IssuanceWeek records from the query.
func (q issuanceWeekQuery) All(ctx context.Context, exec boil.ContextExecutor) (IssuanceWeekSlice, error) {
	var o []*IssuanceWeek

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to IssuanceWeek slice")
	}

	if len(issuanceWeekAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all IssuanceWeek records in the query.
func (q issuanceWeekQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count issuance_weeks rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q issuanceWeekQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if issuance_weeks exists")
	}

	return count > 0, nil
}

// Rewards retrieves all the reward's Rewards with an executor.
func (o *IssuanceWeek) Rewards(mods ...qm.QueryMod) rewardQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"rewards_api\".\"rewards\".\"issuance_week_id\"=?", o.ID),
	)

	query := Rewards(queryMods...)
	queries.SetFrom(query.Query, "\"rewards_api\".\"rewards\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"rewards_api\".\"rewards\".*"})
	}

	return query
}

// LoadRewards allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (issuanceWeekL) LoadRewards(ctx context.Context, e boil.ContextExecutor, singular bool, maybeIssuanceWeek interface{}, mods queries.Applicator) error {
	var slice []*IssuanceWeek
	var object *IssuanceWeek

	if singular {
		object = maybeIssuanceWeek.(*IssuanceWeek)
	} else {
		slice = *maybeIssuanceWeek.(*[]*IssuanceWeek)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &issuanceWeekR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &issuanceWeekR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`rewards_api.rewards`),
		qm.WhereIn(`rewards_api.rewards.issuance_week_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load rewards")
	}

	var resultSlice []*Reward
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice rewards")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on rewards")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for rewards")
	}

	if len(rewardAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Rewards = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &rewardR{}
			}
			foreign.R.IssuanceWeek = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.IssuanceWeekID {
				local.R.Rewards = append(local.R.Rewards, foreign)
				if foreign.R == nil {
					foreign.R = &rewardR{}
				}
				foreign.R.IssuanceWeek = local
				break
			}
		}
	}

	return nil
}

// AddRewards adds the given related objects to the existing relationships
// of the issuance_week, optionally inserting them as new records.
// Appends related to o.R.Rewards.
// Sets related.R.IssuanceWeek appropriately.
func (o *IssuanceWeek) AddRewards(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Reward) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.IssuanceWeekID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"rewards_api\".\"rewards\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"issuance_week_id"}),
				strmangle.WhereClause("\"", "\"", 2, rewardPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.IssuanceWeekID, rel.UserDeviceID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.IssuanceWeekID = o.ID
		}
	}

	if o.R == nil {
		o.R = &issuanceWeekR{
			Rewards: related,
		}
	} else {
		o.R.Rewards = append(o.R.Rewards, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &rewardR{
				IssuanceWeek: o,
			}
		} else {
			rel.R.IssuanceWeek = o
		}
	}
	return nil
}

// IssuanceWeeks retrieves all the records using an executor.
func IssuanceWeeks(mods ...qm.QueryMod) issuanceWeekQuery {
	mods = append(mods, qm.From("\"rewards_api\".\"issuance_weeks\""))
	return issuanceWeekQuery{NewQuery(mods...)}
}

// FindIssuanceWeek retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindIssuanceWeek(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*IssuanceWeek, error) {
	issuanceWeekObj := &IssuanceWeek{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"rewards_api\".\"issuance_weeks\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, issuanceWeekObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from issuance_weeks")
	}

	if err = issuanceWeekObj.doAfterSelectHooks(ctx, exec); err != nil {
		return issuanceWeekObj, err
	}

	return issuanceWeekObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *IssuanceWeek) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no issuance_weeks provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(issuanceWeekColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	issuanceWeekInsertCacheMut.RLock()
	cache, cached := issuanceWeekInsertCache[key]
	issuanceWeekInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			issuanceWeekAllColumns,
			issuanceWeekColumnsWithDefault,
			issuanceWeekColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(issuanceWeekType, issuanceWeekMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(issuanceWeekType, issuanceWeekMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"rewards_api\".\"issuance_weeks\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"rewards_api\".\"issuance_weeks\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into issuance_weeks")
	}

	if !cached {
		issuanceWeekInsertCacheMut.Lock()
		issuanceWeekInsertCache[key] = cache
		issuanceWeekInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the IssuanceWeek.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *IssuanceWeek) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	issuanceWeekUpdateCacheMut.RLock()
	cache, cached := issuanceWeekUpdateCache[key]
	issuanceWeekUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			issuanceWeekAllColumns,
			issuanceWeekPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update issuance_weeks, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"rewards_api\".\"issuance_weeks\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, issuanceWeekPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(issuanceWeekType, issuanceWeekMapping, append(wl, issuanceWeekPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update issuance_weeks row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for issuance_weeks")
	}

	if !cached {
		issuanceWeekUpdateCacheMut.Lock()
		issuanceWeekUpdateCache[key] = cache
		issuanceWeekUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q issuanceWeekQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for issuance_weeks")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for issuance_weeks")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o IssuanceWeekSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), issuanceWeekPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"rewards_api\".\"issuance_weeks\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, issuanceWeekPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in issuanceWeek slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all issuanceWeek")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *IssuanceWeek) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no issuance_weeks provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(issuanceWeekColumnsWithDefault, o)

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

	issuanceWeekUpsertCacheMut.RLock()
	cache, cached := issuanceWeekUpsertCache[key]
	issuanceWeekUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			issuanceWeekAllColumns,
			issuanceWeekColumnsWithDefault,
			issuanceWeekColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			issuanceWeekAllColumns,
			issuanceWeekPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert issuance_weeks, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(issuanceWeekPrimaryKeyColumns))
			copy(conflict, issuanceWeekPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"rewards_api\".\"issuance_weeks\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(issuanceWeekType, issuanceWeekMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(issuanceWeekType, issuanceWeekMapping, ret)
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
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert issuance_weeks")
	}

	if !cached {
		issuanceWeekUpsertCacheMut.Lock()
		issuanceWeekUpsertCache[key] = cache
		issuanceWeekUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single IssuanceWeek record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *IssuanceWeek) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no IssuanceWeek provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), issuanceWeekPrimaryKeyMapping)
	sql := "DELETE FROM \"rewards_api\".\"issuance_weeks\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from issuance_weeks")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for issuance_weeks")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q issuanceWeekQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no issuanceWeekQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from issuance_weeks")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for issuance_weeks")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o IssuanceWeekSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(issuanceWeekBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), issuanceWeekPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"rewards_api\".\"issuance_weeks\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, issuanceWeekPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from issuanceWeek slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for issuance_weeks")
	}

	if len(issuanceWeekAfterDeleteHooks) != 0 {
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
func (o *IssuanceWeek) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindIssuanceWeek(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *IssuanceWeekSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := IssuanceWeekSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), issuanceWeekPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"rewards_api\".\"issuance_weeks\".* FROM \"rewards_api\".\"issuance_weeks\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, issuanceWeekPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in IssuanceWeekSlice")
	}

	*o = slice

	return nil
}

// IssuanceWeekExists checks if the IssuanceWeek row exists.
func IssuanceWeekExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"rewards_api\".\"issuance_weeks\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if issuance_weeks exists")
	}

	return exists, nil
}
