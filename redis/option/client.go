package option

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

// Client keeps the settings to set up redis connection.
type Client struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string
	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string
	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)
	// Hook that is called when new connection is established.
	OnConnect func(ctx context.Context, cn *redis.Conn) error
	// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
	// Default is 3.
	Protocol int
	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string
	// Optional password. Must match the password specified in the
	// requirement pass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string
	// CredentialsProvider allows the username and password to be updated
	// before reconnecting. It should return the current username and password.
	CredentialsProvider func() (username string, password string)
	// Database to be selected after connecting to the server.
	DB int
	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetReadDeadline calls completely.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.  Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetWriteDeadline calls completely.
	WriteTimeout time.Duration
	// ContextTimeoutEnabled controls whether the client respects context timeouts and deadlines.
	// See https://redis.uptrace.dev/guide/go-redis-debugging.html#timeouts
	ContextTimeoutEnabled bool
	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that FIFO has slightly higher overhead compared to LIFO,
	// but it helps to close idle connections faster reducing the pool size.
	PoolFIFO bool
	// Base number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	// If there is not enough connections in the pool, new connections will be allocated in excess of PoolSize,
	// you can limit it through MaxActiveConns
	PoolSize int
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	// Default is 0. the idle connections are not closed by default.
	MinIdleConns int
	// Maximum number of idle connections.
	// Default is 0. the idle connections are not closed by default.
	MaxIdleConns int
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActiveConns int
	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	// Should be less than server's timeout.
	//
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's idle time.
	//
	// Default is 30 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration
	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	// If <= 0, connections are not closed due to a connection's age.
	//
	// Default is to not close idle connections.
	ConnMaxLifetime time.Duration
	// TLS Config to use. When set, TLS will be negotiated.
	TLSConfig *tls.Config
	// Limiter interface used to implement circuit breaker or rate limiter.
	Limiter Limiter
	// Enables read only queries on slave/follower nodes.
	readOnly bool
	// Disable set-lib on connect. Default is false.
	DisableIndentity bool
}

// Limiter is the interface of a rate limiter or a circuit breaker.
type Limiter interface {
	// Allow returns nil if operation is allowed or an error otherwise.
	// If operation is allowed client must ReportResult of the operation
	// whether it is a success or a failure.
	Allow() error
	// ReportResult reports the result of the previously allowed operation.
	// nil indicates a success, non-nil error usually indicates a failure.
	ReportResult(result error)
}

func (c Client) ParseToRedisOptions() *redis.Options {
	return &redis.Options{
		Network:               c.Network,
		Addr:                  c.Addr,
		ClientName:            c.ClientName,
		Dialer:                c.Dialer,
		OnConnect:             c.OnConnect,
		Protocol:              c.Protocol,
		Username:              c.Username,
		Password:              c.Password,
		CredentialsProvider:   c.CredentialsProvider,
		DB:                    c.MaxIdleConns,
		MaxRetries:            c.MaxRetries,
		MinRetryBackoff:       c.MinRetryBackoff,
		MaxRetryBackoff:       c.MaxRetryBackoff,
		DialTimeout:           c.DialTimeout,
		ReadTimeout:           c.ReadTimeout,
		WriteTimeout:          c.WriteTimeout,
		ContextTimeoutEnabled: c.ContextTimeoutEnabled,
		PoolFIFO:              c.PoolFIFO,
		PoolSize:              c.PoolSize,
		PoolTimeout:           c.PoolTimeout,
		MinIdleConns:          c.MinIdleConns,
		MaxIdleConns:          c.MaxActiveConns,
		MaxActiveConns:        c.MaxActiveConns,
		ConnMaxIdleTime:       c.ConnMaxIdleTime,
		ConnMaxLifetime:       c.ConnMaxLifetime,
		TLSConfig:             c.TLSConfig,
		Limiter:               c.Limiter,
		DisableIndentity:      c.DisableIndentity,
	}
}
