package database

import (
	"context"
	"database/sql"
	"time"

	"ariga.io/entcache"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/nhatquangsin/game-service/infra/repo/entc"
	viperutil "github.com/nhatquangsin/game-service/infra/utils/viper"
)

// FXModule represents a fx module for DB Conn.
var FXModule = fx.Provide(
	func(v *viper.Viper, lc fx.Lifecycle) (Client, error) {
		d, err := New(Options.ConfigLoader(v))
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				if err := d.MasterDB(ctx).Close(); err != nil {
					return err
				}

				if err := d.Master(ctx).Close(); err != nil {
					return err
				}

				return d.Slave(ctx).Close()
			},
		})

		return d, nil
	},
)

// Option is a function that sets some option on the DB Conn.
type Option func(*client) error

// Options is a factory for all available DB Conn's.
var Options options

type options struct{}

// MasterURI is an Option to set the master uri of client.
func (options) MasterURI(uri string) Option {
	return func(c *client) error {
		c.master = uri
		return nil
	}
}

// SlaveURI is an Option to set the slave uri of client.
func (options) SlaveURI(uri string) Option {
	return func(c *client) error {
		c.slave = uri
		return nil
	}
}

// MaxIdleConns is an Option to set the max idle connections for client.
func (options) MaxIdleConns(maxConns int) Option {
	return func(c *client) error {
		c.maxIdleConns = maxConns
		return nil
	}
}

// MaxActiveConns is an Option to set the max active connections for client.
func (options) MaxActiveConns(maxConns int) Option {
	return func(c *client) error {
		c.maxActiveConns = maxConns
		return nil
	}
}

// MaxConnTimeout is an Option to set the max lifetime of connections for client.
func (options) MaxConnTimeout(maxTime time.Duration) Option {
	return func(c *client) error {
		c.maxConnTimeout = maxTime
		return nil
	}
}

// RedisCacheURI is an Option to set the uri of redis client for caching.
func (options) RedisCacheURI(uri string) Option {
	return func(c *client) error {
		c.redisCacheURI = uri
		return nil
	}
}

// RedisCacheEnable is an Option to enable redis cache by query.
func (options) RedisCacheEnable(enable bool) Option {
	return func(c *client) error {
		c.redisCacheEnable = enable
		return nil
	}
}

// RedisCacheTTL is an Option to set the cache ttl in case using query redis cache.
func (options) RedisCacheTTL(ttl time.Duration) Option {
	return func(c *client) error {
		c.redisCacheTTL = ttl
		return nil
	}
}

// ConfigLoader is an Options to load all options that configured through viper.
// Client's options that are loaded from viper by following keys:
//
//   - postgres.master.uri
//     db uri to connect with master connection type.
//   - postgres.slave.uri
//     db uri to connect with slave connection type.
func (options) ConfigLoader(v *viper.Viper) Option {
	return func(c *client) error {
		cfg := struct {
			Postgres struct {
				Master struct {
					URI string `mapstructure:"uri"`
				} `mapstructure:"master"`
				Slave struct {
					URI string `mapstructure:"uri"`
				} `mapstructure:"slave"`
				MaxIdleConns   int           `mapstructure:"max_idle_conns"`
				MaxActiveConns int           `mapstructure:"max_active_conns"`
				MaxConnTimeout time.Duration `mapstructure:"max_conn_timeout"`
				RedisCache     struct {
					URI    string        `mapstructure:"uri"`
					Enable bool          `mapstructure:"enable"`
					TTL    time.Duration `mapstructure:"ttl"`
				} `mapstructure:"redis_cache"`
			} `mapstructure:"postgres"`
		}{}

		if err := viperutil.Unmarshal(v, &cfg); err != nil {
			return err
		}

		var opts []Option
		opts = append(
			opts,
			Options.MasterURI(cfg.Postgres.Master.URI),
			Options.SlaveURI(cfg.Postgres.Slave.URI),
			Options.MaxActiveConns(cfg.Postgres.MaxActiveConns),
			Options.MaxIdleConns(cfg.Postgres.MaxIdleConns),
			Options.MaxConnTimeout(cfg.Postgres.MaxConnTimeout),
			Options.RedisCacheURI(cfg.Postgres.RedisCache.URI),
			Options.RedisCacheEnable(cfg.Postgres.RedisCache.Enable),
			Options.RedisCacheTTL(cfg.Postgres.RedisCache.TTL),
		)

		for _, opt := range opts {
			if err := opt(c); err != nil {
				return err
			}
		}

		return nil
	}
}

// Client define all method provide by client.
type Client interface {
	MasterDB(context.Context) *sql.DB
	Master(context.Context) *entc.Client
	SlaveDB(context.Context) *sql.DB
	Slave(context.Context) *entc.Client
	DBName() string
}

type client struct {
	dbname           string
	master           string
	slave            string
	masterDB         *sql.DB
	slaveDB          *sql.DB
	maxIdleConns     int
	maxActiveConns   int
	maxConnTimeout   time.Duration
	redisCacheEnable bool
	redisCacheTTL    time.Duration
	redisCacheURI    string
	MasterPool       *entc.Client
	SlavePool        *entc.Client
}

// New creates and return new instance of Client.
func New(opts ...Option) (Client, error) {
	c := client{}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}

	masterDB, master, err := NewClient(c.master, c.maxIdleConns, c.maxActiveConns, c.maxConnTimeout, c.redisCacheEnable, c.redisCacheURI, c.redisCacheTTL)
	if err != nil {
		return nil, err
	}

	// get db name and assign it to client.dbname
	rows, err := masterDB.Query("SELECT current_database()")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var currentDB string
		if err := rows.Scan(&currentDB); err != nil {
			return nil, err
		}

		c.dbname = currentDB
	}

	c.MasterPool = master
	c.masterDB = masterDB

	slaveDB, slave, err := NewClient(c.slave, c.maxIdleConns, c.maxActiveConns, c.maxConnTimeout, c.redisCacheEnable, c.redisCacheURI, c.redisCacheTTL)
	if err != nil {
		return nil, err
	}

	c.SlavePool = slave
	c.slaveDB = slaveDB

	return &c, nil
}

// NewClient initialize db connection.
func NewClient(
	uri string,
	maxIdleConns,
	maxActiveConns int,
	maxConnTimeout time.Duration,
	redisCacheEnable bool,
	redisCacheURI string,
	redisCacheTTL time.Duration,
) (*sql.DB, *entc.Client, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxActiveConns)
	db.SetConnMaxLifetime(maxConnTimeout)

	drv := entsql.OpenDB(dialect.Postgres, db)
	if redisCacheEnable {
		rdb := redis.NewClient(&redis.Options{
			Addr: redisCacheURI,
		})
		drv := entcache.NewDriver(
			drv,
			entcache.TTL(redisCacheTTL),
			entcache.Levels(
				entcache.NewLRU(256),
				entcache.NewRedis(rdb),
			),
		)

		client := entc.NewClient(entc.Driver(drv))
		return db, client, nil
	}

	client := entc.NewClient(entc.Driver(drv))

	return db, client, nil
}

// DBName to return name of current db.
func (c *client) DBName() string {
	return c.dbname
}

// Master to return master pool.
func (c *client) Master(context.Context) *entc.Client {
	return c.MasterPool
}

// Slave to return slave pool.
func (c *client) Slave(context.Context) *entc.Client {
	return c.SlavePool
}

// MasterDB to return DB connection without ent.Client
func (c *client) MasterDB(context.Context) *sql.DB {
	return c.masterDB
}

// SlaveDB to return DB connection without ent.Client
func (c *client) SlaveDB(context.Context) *sql.DB {
	return c.slaveDB
}
