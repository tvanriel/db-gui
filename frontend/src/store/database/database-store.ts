import {type SchemaService} from '../../services/schema-service';
import {defineStore} from 'pinia';
import {type Token} from '../auth/auth-store';
import {doResponse, mapAndDo} from '../../lib/api-gateway';
import {Fzf} from 'fzf';

export type Column = {
	name: string;
	type: string;
	default: string;
	comment: string;
};
export type Index = {
	type: string;
	content: string;
};

export type Table = {
	name: string;
	columns: Column[];
	indexes?: Index[];
};

export type Database = {
	tableNames: string[];
	tables: Record<string, Table | undefined>;
	name: string;

	tableSearch: string;
	selectedTable: string;
};


export const makeDatabaseStore = (schemaService: SchemaService) => defineStore('database', {
	state: () => ({
		loading: false,
		databasesList: {} as Record<Token, Record<string, Database>>,
		errors: {} as Record<Token, string[]>,
		currentDatabase: {} as Record<Token, string | undefined>,
	}),

	persist: {
		paths: [
			'databasesList',
			'errors',
			'currentDatabase',
		],
	},

	getters: {
		getCurrentDatabase: state => (token: Token): Database | undefined => {
			const currentDatabaseName = state.currentDatabase[token];
			if (currentDatabaseName === undefined || currentDatabaseName === null) {
				return undefined;
			}

			const databases = state.databasesList[token] ?? null;
			if (databases === null) {
				return undefined;
			}

			return databases[currentDatabaseName] ?? null;
		},

		getCurrentTable: state => (token: Token, database: string): Table | undefined => {
			const currentDatabaseName = state.currentDatabase[token];
			if (currentDatabaseName === undefined || currentDatabaseName === null) {
				return undefined;
			}

			const databases = state.databasesList[token] ?? null;
			if (databases === null) {
				return undefined;
			}

			const db = databases[currentDatabaseName];
			if (db === undefined || db === null) {
				return undefined;
			}

			const table = db.tables[db.selectedTable];
			if (table === undefined || table === null) {
				return undefined;
			}

			return table;
		},

		hasCurrentDatabase: state => (token: Token) => typeof state.currentDatabase[token] === 'string',

		mustGetConnection: state => (token: Token) =>
			state.databasesList[token],

		mustGetCurrentDatabase() {
			return (token: Token): Database => this.mustGetConnection(token)[this.currentDatabase[token] ?? ''];
		},
	},

	actions: {

		selectTable(token: Token, databaseName: string, table: string) {
			this.databasesList[token][databaseName].selectedTable = table;
			this.loadTables(token, databaseName);
			this.loadTableDescription(token, databaseName, table);
		},

		removeDatabases(token: Token): void {
			delete this.databasesList[token];
		},

		selectDatabase(token: Token, databaseName: string) {
			this.currentDatabase[token] = databaseName;
		},

		deselectDatabase(token: Token) {
			delete this.currentDatabase[token];
		},

		loadDatabases(token: Token): void {
			this.loading = true;
			schemaService.listSchemas(token).then(mapAndDo(
				({item}): Database => ({
					tableNames: [],
					tables: {},
					name: item,
					tableSearch: '',
					selectedTable: '',
				}),
				result => {
					this.loading = false;
					const m = {} as Record<string, Database>;
					result.forEach(database => m[database.name] = database);
					this.databasesList[token] = m;
				},
				e => {
					this.loading = false;
					this.addError(token, e.error());
				},
			));
		},

		loadTableDescription(token: Token, databaseName: string, tableName: string): void {
			this.loading = true;
			schemaService.describeTable(token, databaseName, tableName).then(doResponse(
				item => {
					this.loading = false;
					this.describeTable(token, databaseName, tableName, item);
				},
				error => {
					this.loading = false;
					this.addError(token, error.error());
				},
			));
		},

		describeTable(token: Token, databaseName: string, tableName: string, tableDescription: Table): void {
			const conn = this.databasesList[token];
			if (conn === undefined) {
				return;
			}

			const db = conn[databaseName];
			if (db === undefined) {
				return;
			}

			db.tables[tableName] = tableDescription;
		},

		loadTables(token: Token, databaseName: string) {
			this.loading = true;
			return schemaService.listTables(token, databaseName).then(doResponse(
				result => {
					this.loading = false;
					this.addTables(token, databaseName, result);
				},
				error => {
					this.loading = false;
					this.addError(token, error.error());
				},
			));
		},

		addTables(token: Token, databaseName: string, result: string[]) {
			const conn = this.databasesList[token];
			if (conn === undefined) {
				return;
			}

			const db = conn[databaseName];
			if (db === undefined) {
				return;
			}

			db.tableNames = result;
		},

		addError(token: Token, error: string): void {
			const errors = this.errors[token];

			if (errors !== undefined) {
				errors.push(error);
			}
		},

		dismissError(token: Token, i: number): void {
			const errors = this.errors[token];

			if (errors !== undefined) {
				errors.splice(i, 1);
			}
		},
	},
});
