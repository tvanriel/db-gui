import {type ApiGateway} from '../lib/api-gateway';
import {type SqlResult} from './sql-service';

export type SelectionConfig = {
	where: Array<{
		columnName: string;
		operator: string;
		value: string;
	}>;
	columns: Array<{
		name: string;
		func: string;
	}>;
	order: Array<{
		direction: 'ASC' | 'DESC';
		by: string;
	}>;
	limit: {
		limit: number;
		offset: number;
	};
};

export type SelectResponse = {
	results: SqlResult;
	sql: string;
};

export class SelectService {
	constructor(private readonly api: ApiGateway) {}
	public async select(token: string, databaseName: string, tableName: string, config: SelectionConfig) {
		return this.api.post<SelectResponse>('/select', {databaseName, tableName, config}, {token});
	}
}
