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
func NewDatasource(driverName, datasourceName string) (*dbsql.DB, error) {
	return dbsql.Open(driverName, datasourceName)
}

// TODO: Prepared statement cache.
