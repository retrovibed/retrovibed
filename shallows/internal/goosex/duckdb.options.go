package goosex

// Option configures InitializeDatabase.
type Option func(*config)

type config struct {
	table string
	pre   []string
	post  []string
}

// WithTableName sets the table tracking the applied migrations. qualify it with a
// catalog, e.g.) cache.goose_db_version, to keep it in an attached database
// regardless of the default catalog of the connection that checks for pending
// migrations, which goose does before the migration lock is taken and so before
// any WithPreStatements have run.
func WithTableName(name string) Option {
	return func(c *config) {
		c.table = name
	}
}

// WithPreStatements runs statements on the connection the migrations run on
// before the migration lock is taken, e.g.) to ATTACH a database and USE it so
// the lock table, the version table, and the unqualified statements of the
// migrations all apply to it. statements must be safe to repeat.
func WithPreStatements(statements ...string) Option {
	return func(c *config) {
		c.pre = append(c.pre, statements...)
	}
}

// WithPostStatements runs statements on the connection the migrations ran on
// after the migration lock is released, e.g.) to restore the catalog selected
// by WithPreStatements.
func WithPostStatements(statements ...string) Option {
	return func(c *config) {
		c.post = append(c.post, statements...)
	}
}
