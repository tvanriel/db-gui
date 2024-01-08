package communicator

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/formatting"
)

var _ domain.ScriptExecutionCommunicator = &WebsocketScriptExecutionCommunicator{}

type WebsocketScriptExecutionCommunicator struct {
	ws *websocket.Conn
}

type WsSqlQueryResult struct {
	Type   WsMessageType              `json:"type"`
	Sql    string                     `json:"sql"`
	Result []formatting.ResultSetJson `json:"results"`
	Time   int64                      `json:"microseconds"`
}

type WsSqlExecResult struct {
	Type         WsMessageType `json:"type"`
	RowsAffected int64         `json:"rowsAffected"`
	LastInsertId int64         `json:"lastInsertId"`
	Sql          string        `json:"sql"`
	Time         int64         `json:"microseconds"`
}

type WsParseProgress struct {
	Type     WsMessageType `json:"type"`
	Progress int           `json:"progress"`
}

type WsExecuteProgress struct {
	Type     WsMessageType `json:"type"`
	Executed int           `json:"executed"`
	Total    int           `json:"total"`
}

type WsStatementPlan struct {
	Type WsMessageType `json:"type"`
	Sql  string        `json:"sql"`
	Time int64         `json:"time"`
}

type WsShouldContinueFetchingMessage struct {
	Type    WsMessageType `json:"type"`
	Fetched int           `json:"fetched"`
}

type WsShouldContinueFetchingResponseMessage struct {
	Type     WsMessageType `json:"type"`
	Continue bool          `json:"continue"`
}

type WsError struct {
	Type    WsMessageType     `json:"type"`
	Message string            `json:"message"`
	Params  map[string]string `json:"params"`
}

func (w *WsError) With(key, val string) *WsError {
	w.Params[key] = val
	return w
}

type WsMessageType string

const (
	WS_PROGRESS_PARSING   WsMessageType = "sql/handler/progress_parsing"
	WS_PROGRESS_EXECUTING WsMessageType = "sql/handler/progress_executing"
	WS_STATEMENT_PLAN     WsMessageType = "sql/handler/plan"
	WS_QUERY              WsMessageType = "sql/handler/query"
	WS_EXEC               WsMessageType = "sql/handler/exec"
	WS_SHOULD_CONTINUE    WsMessageType = "sql/handler/should_continue"
	WS_STATEMENT_ERROR    WsMessageType = "sql/handler/statement_error"
	WS_ERROR              WsMessageType = "sql/handler/error"
)

func (c *WebsocketScriptExecutionCommunicator) ParsingProgress(i int) {
	_ = c.ws.WriteJSON(WsParseProgress{
		Type:     WS_PROGRESS_PARSING,
		Progress: i,
	})
}

func (c *WebsocketScriptExecutionCommunicator) ExecutionProgress(executed, total int) {
	_ = c.ws.WriteJSON(WsExecuteProgress{
		Type:     WS_PROGRESS_EXECUTING,
		Executed: executed,
		Total:    total,
	})
}

func (c *WebsocketScriptExecutionCommunicator) StatementPlan(sql string) {
  _ = c.ws.WriteJSON(WsStatementPlan{
		Type: WS_STATEMENT_PLAN,
		Sql:  sql,
		Time: time.Now().UnixMilli(),
	})
}

func (c *WebsocketScriptExecutionCommunicator) QueryResult(sql string, results domain.Resultset, Time int64) {
	_ = c.ws.WriteJSON(WsSqlQueryResult{
		Type:   WS_QUERY,
		Sql:    sql,
		Result: formatting.ToJSON(results),
		Time:   Time,
	})
}

func (c *WebsocketScriptExecutionCommunicator) ExecResult(rowsAffected, lastInsertId int64, sql string, Time int64) {
	_=c.ws.WriteJSON(WsSqlExecResult{
		Type:         WS_EXEC,
		RowsAffected: rowsAffected,
		LastInsertId: lastInsertId,
		Sql:          sql,
		Time:         Time,
	})
}

func (c *WebsocketScriptExecutionCommunicator) ShouldContinueFetching(fetched int) bool {
	_ = c.ws.WriteJSON(WsShouldContinueFetchingMessage{
		Type:    WS_SHOULD_CONTINUE,
		Fetched: fetched,
	})
	answer := &WsShouldContinueFetchingResponseMessage{}

	err := c.ws.ReadJSON(answer)

	if err != nil {
		c.GenericError(err)
		return false
	}

	return answer.Continue
}

func (c *WebsocketScriptExecutionCommunicator) StatementError(err error, sql string) {
	_ = c.ws.WriteJSON(WsError{
		Type:    WS_STATEMENT_ERROR,
		Message: err.Error(),
		Params: map[string]string{
			"sql": sql,
		},
	})
}

func (c *WebsocketScriptExecutionCommunicator) GenericError(err error) {
  _ = c.ws.WriteJSON(WsError{
		Type:    WS_ERROR,
		Message: err.Error(),
		Params:  make(map[string]string),
	})
}

func (c *WebsocketScriptExecutionCommunicator) Close() error {
	return c.ws.Close()
}

var _ domain.ScriptExecutionCommunicator = &ErrorCommunicator{}

type ErrorCommunicator struct {
	ReportStatementError func(err error, sql string)
	ReportGenericError   func(err error)
}

func (c *ErrorCommunicator) Close() error                                { return nil }
func (c *ErrorCommunicator) QueryResult(string, domain.Resultset, int64) {}
func (c *ErrorCommunicator) StatementPlan(string)                        {}
func (c *ErrorCommunicator) ExecResult(int64, int64, string, int64)      {}
func (c *ErrorCommunicator) ExecutionProgress(int, int)                  {}
func (c *ErrorCommunicator) ParsingProgress(int)                         {}
func (c *ErrorCommunicator) ShouldContinueFetching(int) bool             { return true }
func (c *ErrorCommunicator) GenericError(err error)                      { c.ReportGenericError(err) }
func (c *ErrorCommunicator) StatementError(err error, sql string)        { c.ReportStatementError(err, sql) }

func NewWsCommunicator(ws *websocket.Conn) domain.ScriptExecutionCommunicator {
	return &WebsocketScriptExecutionCommunicator{ws: ws}
}

func NewErrorCommunicator(ReportStatementError func(err error, sql string), ReportGenericError func(err error)) domain.ScriptExecutionCommunicator {
	return &ErrorCommunicator{
		ReportStatementError: ReportStatementError,
		ReportGenericError:   ReportGenericError,
	}
}
