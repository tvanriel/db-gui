import {type ApiGateway} from '../lib/api-gateway';

export type SqlResultColumnDescription = {
	name: string;
	type: string;
};

export type SqlResultValue = {
	null: boolean;
	value: string;
};

export type SqlResult = Array<{
	columns: SqlResultColumnDescription[];
	result: SqlResultValue[][];
	error: string;
}>;

export type SqlHandlerProgressReporter = {
	clean: () => void;
	progress: (p: SqlHandlerParsingProgressMessage | SqlHandlerExecutingProgressMessage) => void;
	error: (e: SqlHandlerErrorMessage | SqlHandlerStatementErrorMessage) => void;
	result: (r: SqlHandlerExecuteResultMessage | SqlHandlerQueryResultMessage) => void;
	plan: (r: SqlHandlerStatementPlanMessage) => void;
	close: () => void;
	shouldContinue: (fetched: number) => Promise<boolean>;
};

export type SqlHandlerParsingProgressMessage = {
	type: 'sql/handler/progress_parsing';
	progress: number;
};

export type SqlHandlerStatementPlanMessage = {
	type: 'sql/handler/plan';
	sql: string;
	time: number;
};

export type SqlHandlerExecutingProgressMessage = {
	type: 'sql/handler/progress_executing';
	executed: number;
	total: number;
};

export type SqlHandlerErrorMessage = {
	type: 'sql/handler/error';
	message: string;
};
export type SqlHandlerStatementErrorMessage = {
	type: 'sql/handler/statement_error';
	message: string;
	params: {
		sql: string;
	};
};

export type SqlHandlerQueryResultMessage = {
	type: 'sql/handler/query';
	sql: string;
	microseconds: number;
	results: SqlResult;
};
export type SqlHandlerExecuteResultMessage = {
	type: 'sql/handler/exec';
	sql: string;
	microseconds: number;
	lastInsertId: number;
	rowsAffected: number;
};

export type SqlShouldContinueMessage = {
	type: 'sql/handler/should_continue';
	fetched: number;
};

export type SqlHandlerMessage = SqlShouldContinueMessage | SqlHandlerParsingProgressMessage | SqlHandlerExecutingProgressMessage | SqlHandlerErrorMessage | SqlHandlerQueryResultMessage | SqlHandlerExecuteResultMessage | SqlHandlerStatementPlanMessage | SqlHandlerStatementErrorMessage;

export class SqlService {
	constructor(private readonly api: ApiGateway) {}
	public query(token: string, databaseName: string, sql: string, handler: SqlHandlerProgressReporter) {
		const queryString = new Map<string, string>();
		queryString.set('token', token);

		this.api.ws('/sql', queryString).then(ws => {
			ws.send(JSON.stringify({sql, databaseName}));
			handler.clean();
			ws.addEventListener('message', ev => {
				try {
					const data = JSON.parse(ev.data) as SqlHandlerMessage;

					switch (data.type) {
						case 'sql/handler/error':
						case 'sql/handler/statement_error':
							handler.error(data);
							break;
						case 'sql/handler/exec':
						case 'sql/handler/query':
							handler.result(data);
							break;
						case 'sql/handler/progress_executing':
						case 'sql/handler/progress_parsing':
							handler.progress(data);
							break;
						case 'sql/handler/plan':
							handler.plan(data);
							break;
						case 'sql/handler/should_continue':
							handler.shouldContinue(data.fetched).then(shouldContinue => {
								ws.send(JSON.stringify({
									continue: shouldContinue,
									type: 'sql/handler/should_continue',
								}));
							});
					}
				} catch (e) {

				}
			});
			ws.addEventListener('close', () => {
				handler.close();
			});
		});
	}
}
