import {type ApiGateway} from '../lib/api-gateway';
import {type AxiosResponse} from 'axios';

export class ExportService {
	constructor(private readonly api: ApiGateway) {}

	public async export({token, databaseName}: {token: string; databaseName: string}) {
		return this.api.rawPost(this.api.getBaseUrl() + '/export', {databaseName, createTableStrategy: 'DROP+CREATE'}, {responseType: 'blob', params: {token}}).then((response: AxiosResponse) => {

		});
	}
}
