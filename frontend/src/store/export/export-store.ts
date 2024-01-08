import {type ExportService} from '../../services/export-service';
import {defineStore} from 'pinia';
import {type Token} from '../auth/auth-store';

export type ExportProgress = {
	tables: string[];
	exportedTables: number;

	views: string[];
	exportedViews: number;
};

export const makeExportStore = (exportService: ExportService) => defineStore('export', {
	state: () => ({
		runningExport: new Map<Token, ExportProgress>(),
		selectedDatabase: new Map<Token, string>(),
	}),

	getters: {
		getCurrentExport() {
			return (token: Token): ExportProgress | undefined => this.runningExport.get(token);
		},
	},

	actions: {
		selectDatabase(token: Token, databaseName: string) {
			this.selectedDatabase.set(token, databaseName);
		},

		startExport(token: Token): void {
			exportService.export({
				databaseName: this.selectedDatabase.get(token) ?? '',
				token,
			});
		},
	},
});
