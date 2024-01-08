package connectionpool

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/scriptauthentication"
)

var _ domain.ConnectionPool = &MemoryConnectionPool{}

type MemoryConnectionPool struct {
	addConnection func(string, string, string) domain.Connection
	Connections   map[string]PoolEntry
}
type PoolEntry struct {
	Connection domain.Connection
	Expires    time.Time
	Host       string
	Username   string
	Password   string
	Script     string
}

func (m *MemoryConnectionPool) Add(config *domain.ConnectionConfig, token string) error {

	conn, resp, err := m.refreshBasedOnScript(config.Script)
	if err != nil {
		return err
	}

	if err = conn.Connect(); err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}

	m.Connections[token] = PoolEntry{
		Connection: conn,
		Username:   resp.Username,
		Password:   resp.Password,
		Host:       resp.Host,
		Script:     config.Script,
		Expires:    time.Now().Add(time.Duration(int(time.Second) * resp.Expires)),
	}

	conn.Close()
	return nil
}

func (m *MemoryConnectionPool) Delete(token string) {
	delete(m.Connections, token)
}
func (m *MemoryConnectionPool) Get(token string) domain.Connection {
	entry, ok := m.Connections[token]

	if !ok {
		return nil
	}

	now := time.Now()

	if now.Unix() > entry.Expires.Unix() {
		conn, resp, err := m.refreshBasedOnScript(entry.Script)
		if err != nil {
			m.Delete(token)
			return nil
		}
		m.Connections[token] = PoolEntry{
			Connection: conn,
			Username:   resp.Username,
			Password:   resp.Password,
			Host:       resp.Host,
			Script:     entry.Script,
			Expires:    time.Now().Add(time.Duration(int(time.Second) * resp.Expires)),
		}

	}

	return entry.Connection
}

func (m *MemoryConnectionPool) NewToken() string {

	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(b)
}

func (m *MemoryConnectionPool) refreshBasedOnScript(script string) (domain.Connection, *scriptauthentication.AuthScriptResponse, error) {

	resp, err := scriptauthentication.ExecuteAuthScript(script)

	if err != nil {
		return nil, nil, err
	}

	conn := m.addConnection(resp.Username, resp.Password, resp.Host)

	return conn, resp, err
}

func NewMemoryConnectionPool(addConnection func(string, string, string) domain.Connection) domain.ConnectionPool {
	return &MemoryConnectionPool{
		addConnection: addConnection,
		Connections:   make(map[string]PoolEntry),
	}
}
