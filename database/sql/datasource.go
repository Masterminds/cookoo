// SQL datasource and commands for Cookoo.
package sql

import (
	dbsql "database/sql"
)

// Create a new SQL datasource.
//
// Currently, this is an empty wrapper around the built-in DB object.
//
// Example:
//	ds, err := sql.NewDatasource("mysql", "root@/mpbtest")
//	if err != nil {
//		panic("Could not create a database connection.")
//		return
//	}
//
//	cxt.AddDatasource("db", ds)
//
// In the example above, we create a new datasource and then add it to
// the context. This should be done at server init, before web.Serve
// or router.HandleRequest().
func NewDbDatasource(driverName, datasourceName string) (*dbsql.DB, error) {
	return dbsql.Open(driverName, datasourceName)
}

// TODO: Prepared statement cache.
// Create a new cache for prepared statements.
//
// Initial capacity determines how big the cache will be.
//
// Warning: The implementation of the caching layer will likely
// change from relatively static to an LRU.
func NewStmtCache(dbHandle *dbsql.DB, initialCapacity int) StmtCache {
	c := new(StmtCacheMap)
	c.cache = make(map[string]*dbsql.Stmt, initialCapacity)
	c.capacity = initialCapacity
	c.dbh = dbHandle

	return c
}

// A StmtCache caches SQL prepared statements.
//
// It's intended use is as a datsource for a long-running SQL-backed
// application. Prepared statements can exist across requests and be
// shared by separate goroutines. For frequently executed statements,
// this is both more performant and more secure (at least for some
// drivers).
//
// IMPORTANT: Statments are cached by string key, so it is important that to
// get the most out of the cache, you re-use the same strings. Otherwise,
// 'SELECT surname, name FROM names' will generate a different cache entry
// than 'SELECT name, surname FROM names'.
//
// The cache is driver-agnostic.
type StmtCache interface {
	Get(statment string) (*dbsql.Stmt, error)
	Clear() error
}

type StmtCacheMap struct {
	cache    map[string]*dbsql.Stmt
	capacity int
	dbh      *dbsql.DB
}

// Get a prepared statement from a SQL string.
//
// This will return a cached statement if one exists, otherwise
// this will generate one, insert it into the cache, and return
// the new statement.
//
// It is assumed that the underlying database layer can handle
// parallelism with prepared statements, and we make no effort
// to deal with locking or synchronization.
func (c *StmtCacheMap) Get(statement string) (*dbsql.Stmt, error) {
	if stmt, ok := c.cache[statement]; ok {
		return stmt, nil
	}
	// Else we prepare the statement and then cache it.
	stmt, err := c.dbh.Prepare(statement)
	if err != nil {
		return nil, err
	}

	// Cache by string key.
	c.cache[statement] = stmt

	return stmt, nil
}

// For compatibility with database/sql.DB.Prepare
func (c *StmtCacheMap) Prepare(statement string) (*dbsql.Stmt, error) {
	return c.Get(statement)
}

// Clear the cache.
func (c *StmtCacheMap) Clear() error {
	// While I don't think this is a good idea, it might be necessary. On the
	// flip side, it might cause race conditions if one goroutine is running
	// a query while another is clearing the cache. For now, leaving this
	// to the memory manager.
	//for _, stmt := range c.cache {
	//	stmt.Close()
	//}

	c.cache = make(map[string]*dbsql.Stmt, c.capacity)
	return nil
}
