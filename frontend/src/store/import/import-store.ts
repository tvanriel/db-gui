import {type ImportService} from '../../services/import-service';
import {defineStore} from 'pinia';
import {type Token} from '../auth/auth-store';
import {doResponse} from '../../lib/api-gateway';

type SqlImportError = {
	sql: string | undefined;
	error: string;
};

export type SqlImportResponse = {
	errors: SqlImportError[];
	executed: number;
};

export const makeImportStore = (importService: ImportService) => defineStore('import', {
	state: () => ({
		results: {} as Record<Token, SqlImportResponse>,
		isImporting: {} as Record<Token, boolean>,
		errors: {} as Record<Token, string[]>,
	}),

	actions: {
		setImportResult(token: Token, result: SqlImportResponse) {
			this.results[token] = result;
		},

		setIsImporting(token: Token, isImporting: boolean) {
			this.isImporting[token] = isImporting;
		},

		runImport(token: Token, fd: FormData): void {
			this.setIsImporting(token, true);
			importService.import({token, fd}).then(doResponse(
				result => {
					this.setIsImporting(token, false);
					this.setImportResult(token, result);
				},
				error => {
					this.addError(token, error.error());
					this.setIsImporting(token, false);
				},
			));
		},

		addError(token: Token, error: string): void {
			this.errors[token].push(error);
		},

		dismissError(token: Token, i: number): void {
			this.errors[token].splice(i, 1);
		},
	},
});
