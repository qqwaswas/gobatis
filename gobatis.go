package gobatis

import (
	"context"
	"database/sql"
	"errors"
)

type ResultType string

const (
	resultTypeMap     ResultType = "map"     // result set is a map: map[string]interface{}
	resultTypeMaps    ResultType = "maps"    // result set is a slice, item is map: []map[string]interface{}
	resultTypeStruct  ResultType = "struct"  // result set is a struct
	resultTypeStructs ResultType = "structs" // result set is a slice, item is struct
	resultTypeSlice   ResultType = "slice"   // result set is a value slice, []interface{}
	resultTypeSlices  ResultType = "slices"  // result set is a value slice, item is value slice, []interface{}
	resultTypeArray   ResultType = "array"   //
	resultTypeArrays  ResultType = "arrays"  // result set is a value slice, item is value slice, []interface{}
	resultTypeValue   ResultType = "value"   // result set is single value
)

type GoBatis interface {
	Select(stmt string, param interface{}) func(res interface{}) error
	Insert(stmt string, param interface{}) (int64, int64, error)
	Update(stmt string, param interface{}) (int64, error)
	Delete(stmt string, param interface{}) (int64, error)
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go start
type sqlExecutor interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

var (
	ErrorDBNil       = errors.New("sql.db can not be nil")
	ErrorDbConn      = errors.New("db can not conn")
	ErrorEmptyMapper = errors.New("mapper xml is nil or empty")
	ErrStruct        = errors.New("struct query result must be ptr")
	ErrNoTransaction = errors.New("sql transaction is nil")
)

type DbType string

const (
	dbTypeMySQL    DbType = "mysql"
	dbTypePostgres DbType = "postgres"
)

var debug = true

type Config struct {
	Db          *sql.DB
	MapperPaths []string
	ColumnStyle int
	Debug       bool
}

func NewGoBatis(ctx context.Context, conf *Config) (*Gobatis, error) {
	if nil == conf.Db {
		return nil, ErrorDBNil
	}

	err := conf.Db.Ping()
	if nil != err {
		return nil, ErrorDbConn
	}

	debug = conf.Debug

	mapper, err := loadingMapper(conf.MapperPaths...)
	if nil != err {
		return nil, err
	}

	gb := &Gobatis{
		mappers: mapper,
		runner:&runner{
			executor:conf.Db,
			mappers:mapper,
		},
	}
	gb.ctxStd, gb.cancel = context.WithCancel(ctx)
	return gb, nil
}

type runner struct {
	executor sqlExecutor
	dbType   DbType
	mappers  *mapper
}

// Gobatis
type Gobatis struct {
	ctxStd  context.Context
	cancel  context.CancelFunc
	mappers *mapper
	*runner
}

// Begin Tx
//
// ps：
//  Tx, err := this.Begin()
func (g *Gobatis) Begin() (*runner, error) {
	return g.BeginTx(g.ctxStd, nil)
}

// Begin Tx with ctx & opts
//
// ps：
//  Tx, err := this.BeginTx(ctx, ops)
func (g *Gobatis) BeginTx(ctx context.Context, opts *sql.TxOptions) (*runner, error) {
	db := g.runner.executor.(*sql.DB)
	tx, err := db.BeginTx(ctx, opts)
	if nil != err {
		return nil, err
	}
	return &runner{
		executor: tx,
		mappers:  g.mappers,
	}, nil
}

// Close db
//
// ps：
//  err := this.Close()
func (g *Gobatis) Close() error {
	db := g.runner.executor.(*sql.DB)
	g.cancel()
	return db.Close()
}

// Commit Tx
//
// ps：
//  err := Tx.Commit()
func (r *runner) Commit() error {
	if nil == r.executor {
		return errors.New("tx no running")
	}
	tx, ok := r.executor.(*sql.Tx)
	if ok {
		return tx.Commit()
	}
	return ErrNoTransaction
}

// Rollback Tx
//
// ps：
//  err := Tx.Rollback()
func (r *runner) Rollback() error {
	if nil == r.executor {
		return errors.New("tx no running")
	}
	tx, ok := r.executor.(*sql.Tx)
	if ok {
		return tx.Rollback()
	}
	return ErrNoTransaction
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go end
func (r *runner) Select(stmt string, param interface{}) func(res interface{}) error {
	ms := r.mappers.getMappedStmt(stmt)
	if nil == ms {
		return func(res interface{}) error {
			return errors.New("mapped statement not found:" + stmt)
		}
	}
	ms.dbType = r.dbType

	params := paramProcess(param)

	return func(res interface{}) error {
		executor := &executor{r}
		return executor.query(ms, params, res)
	}
}

// insert(executor string, param interface{})
func (r *runner) Insert(stmt string, param interface{}) (int64, int64, error) {
	ms := r.mappers.getMappedStmt(stmt)
	if nil == ms {
		return 0, 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = r.dbType

	params := paramProcess(param)

	executor := &executor{r}

	lastInsertId, affected, err := executor.update(ms, params)
	if nil != err {
		return 0, 0, err
	}

	return lastInsertId, affected, nil
}

// update(executor string, param interface{})
func (r *runner) Update(stmt string, param interface{}) (int64, error) {
	ms := r.mappers.getMappedStmt(stmt)
	if nil == ms {
		return 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = r.dbType
	params := paramProcess(param)

	executor := &executor{r}

	_, affected, err := executor.update(ms, params)
	if nil != err {
		return 0, err
	}

	return affected, nil
}

// delete(executor string, param interface{})
func (r *runner) Delete(stmt string, param interface{}) (int64, error) {
	return r.Update(stmt, param)
}
