import {defineStore} from 'pinia';
import {type SqlHandlerErrorMessage, type SqlHandlerExecuteResultMessage, type SqlHandlerExecutingProgressMessage, type SqlHandlerParsingProgressMessage, type SqlHandlerQueryResultMessage, type SqlHandlerStatementErrorMessage, type SqlHandlerStatementPlanMessage, type SqlService} from '../../services/sql-service';
import {type Token} from '../auth/auth-store';

export const makeSqlStore = (sqlService: SqlService) => defineStore('sql', {
	state: () => ({
		sqlResult: [] as Array<SqlHandlerExecuteResultMessage | SqlHandlerQueryResultMessage | SqlHandlerErrorMessage | SqlHandlerStatementPlanMessage | SqlHandlerStatementErrorMessage>,
		sqlCurrent: undefined as SqlHandlerStatementPlanMessage | undefined,
		sqlProgress: undefined as SqlHandlerParsingProgressMessage | SqlHandlerExecutingProgressMessage | undefined,
		sqlHistory: [] as Array<{token: string; sql: string}>,
	}),

	getters: {
		getSqlProgress: state => state.sqlProgress,
		getSqlResult: state => state.sqlResult,
	},

	actions: {
		executeSql(token: Token, databaseName: string, sql: string, sqlShouldContinue: (fetched: number) => Promise<boolean>) {
			this.sqlHistory.push({token, sql});

			sqlService.query(token, databaseName ?? '', sql, {
				clean: () => {
					this.sqlResult = [];
					this.sqlCurrent = undefined;
					this.sqlCurrent = undefined;
				},

				error: result => {
					this.sqlResult.push(result);
				},

				plan: plan => {
					this.sqlCurrent = plan;
				},

				result: result => {
					this.sqlResult.push(result);
				},

				progress: report => {
					this.sqlProgress = report;
				},

				close: () => {
					this.sqlProgress = undefined;
					this.sqlCurrent = undefined;
				},

				async shouldContinue(fetched) {
					return sqlShouldContinue(fetched);
				},

			});
		},
	},
});
