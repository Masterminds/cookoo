# Changelog

## v1.2.0 (Upcomming)

* Added Getter interface
* Added the Get* and Has* utility functions (getter.go)
* Added GetFromFirst(string, interface) (Contextvalue, Getter) function
* Added DefaultGetter struct
* Added 'subcommand' param to cli.ParseArgs
* Added 'cli.New' and 'cli.Runner'
* Added 'fmt' package

## v1.1.0 (2014-06-06)

* Added SyncContext function so Contexts are kept in sync via read/write mutex.
* Improved debugging for panics.
* Documentation fixes.

## v1.0.0 (2014-04-23)

* Initial release.
