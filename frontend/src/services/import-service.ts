import {type ApiGateway, doResponse} from '../lib/api-gateway';

type SqlImportError = {
	sql: string | undefined;
	error: string;
};

export type SqlImportResponse = {
	errors: SqlImportError[];
	executed: number;
};

export class ImportService {
	constructor(private readonly api: ApiGateway) {}

	public async import({token, fd}: {token: string; fd: FormData}) {
		return this.api.post<SqlImportResponse>('/import', fd, {token});
	}
}
