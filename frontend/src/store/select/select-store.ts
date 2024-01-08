import {defineStore} from 'pinia';
import {doResponse} from '../../lib/api-gateway';
import {type SelectionConfig, type SelectService} from '../../services/select-service';
import {type SqlResult} from '../../services/sql-service';
import {type Token} from '../auth/auth-store';

const orDefault = <T>(x: T, subject: Record<Token, Record<string, Record<string, T>>>) => (token: Token, databaseName: string, tableName: string) => {
	const conn = subject[token];
	if (conn === undefined) {
		return x;
	}

	const db = conn[databaseName];
	if (db === undefined) {
		return x;
	}

	const table = db[tableName];
	if (table === undefined) {
		return x;
	}

	return table;
};

export const makeSelectStore = (selectService: SelectService) => defineStore('select', {
	state: () => ({
		selectionConfig: {} as Record<Token, Record<string, Record<string, SelectionConfig>>>,
		resultSet: {} as Record<Token, Record<string, Record<string, SqlResult | undefined>>>,
		selectSql: {} as Record<Token, Record<string, Record<string, string | undefined>>>,
		errors: {} as Record<Token, Record<string, Record<string, string[]>>>,
	}),

	getters: {
		selectConfigFor: state => (token: Token, databaseName: string, tableName: string): SelectionConfig => {
			const x: SelectionConfig = ({
				columns: [],
				limit: {
					limit: 50,
					offset: 0,
				},
				order: [],
				where: [],
			});
			const returnvalue = orDefault(x, state.selectionConfig)(token, databaseName, tableName);
			return returnvalue;
		},

		resultSetFor: state => (token: Token, databaseName: string, tableName: string) => {
			const x: SqlResult | undefined = undefined;
			return orDefault<SqlResult | undefined>(x, state.resultSet)(token, databaseName, tableName);
		},

		selectSqlFor: state => (token: Token, databaseName: string, tableName: string) => {
			const x: string | undefined = undefined;
			return orDefault(x, state.selectSql)(token, databaseName, tableName);
		},

		errorsFor: state => (token: Token, databaseName: string, tableName: string) => {
			const x: string[] = [];
			return orDefault<string[]>(x, state.errors)(token, databaseName, tableName);
		},
	},

	actions: {

		setConfigFor(token: Token, databaseName: string, tableName: string, config: SelectionConfig) {
			if (this.selectionConfig[token] === undefined) {
				this.selectionConfig[token] = {};
			}

			if (this.selectionConfig[token][databaseName] === undefined) {
				this.selectionConfig[token][databaseName] = {};
			}

			this.selectionConfig[token][databaseName][tableName] = config;
		},

		updateWhere(
			token: Token,
			databaseName: string,
			tableName: string,
			where: Array<{columnName: string; operator: string; value: string}>,
		) {
			const config = this.selectConfigFor(token, databaseName, tableName);
			if (config === undefined) {
				return;
			}

			config.where = where;
			this.setConfigFor(token, databaseName, tableName, config);
		},

		updateOrder(
			token: Token,
			databaseName: string,
			tableName: string,
			order: Array<{direction: 'ASC' | 'DESC'; by: string}>,
		) {
			const config = this.selectConfigFor(token, databaseName, tableName);
			if (config === undefined) {
				return;
			}

			config.order = order;
			this.setConfigFor(token, databaseName, tableName, config);
		},

		updateCols(
			token: Token,
			databaseName: string,
			tableName: string,
			columns: Array<{func: string; name: string}>,
		) {
			const config = this.selectConfigFor(token, databaseName, tableName);
			config.columns = columns;
			this.setConfigFor(token, databaseName, tableName, config);
		},
		updateLimit(
			token: Token,
			databaseName: string,
			tableName: string,
			limit: {offset: number; limit: number},
		) {
			const config = this.selectConfigFor(token, databaseName, tableName);
			config.limit = limit;
			this.setConfigFor(token, databaseName, tableName, config);
		},

		executeSelect(token: Token, databaseName: string, tableName: string) {
			const config = this.selectConfigFor(token, databaseName, tableName);

			selectService.select(token, databaseName, tableName, config).then(doResponse(
				result => {
					if (this.resultSet[token] === undefined) {
						this.resultSet[token] = {};
					}

					if (this.resultSet[token][databaseName] === undefined) {
						this.resultSet[token][databaseName] = {};
					}

					this.resultSet[token][databaseName][tableName] = result.results;
					if (this.selectSql[token] === undefined) {
						this.selectSql[token] = {};
					}

					if (this.selectSql[token][databaseName] === undefined) {
						this.selectSql[token][databaseName] = {};
					}

					this.selectSql[token][databaseName][tableName] = result.sql;
				},
				err => {
					this.addError(token, databaseName, tableName, err.error());
				},
			));
		},
		addError(token: Token, databaseName: string, tableName: string, error: string) {
			let errors = this.errors[token][databaseName][tableName];
			if (errors === undefined) {
				errors = [] as string[];
			}

			errors.push(error);
		},
		dismissError(token: Token, databaseName: string, tableName: string, i: number) {
			if (this.errors[token] === undefined) {
				this.errors[token] = {};
			}

			if (this.errors[token][databaseName] === undefined) {
				this.errors[token][databaseName] = {};
			}

			if (this.errors[token][databaseName][tableName] === undefined) {
				this.errors[token][databaseName][tableName] = [];
			}

			const errors = this.errors[token][databaseName][tableName];
			errors.splice(i, 1);
		},
	},
});
