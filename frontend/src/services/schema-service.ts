import {type ApiGateway} from '../lib/api-gateway';

export type Schema = {
	getName(): string;
};

class DatabaseSchema implements Schema {
	constructor(private readonly name: string) {}

	public getName(): string {
		return this.name;
	}
}

class Table {
	constructor(private readonly name: string) {}

	public getName(): string {
		return this.name;
	}
}

type ListSchemasApiResponse = Array<{item: string}>;

export type ColumnDescription = {
	name: string;
	type: string;
	comment: string;
	default: string;
};

type TableDescription = {
	columns: ColumnDescription[];
	name: string;
};

export class CreateSchemaRequest {
	constructor(private readonly name: string) {}
	public getName(): string {
		return this.name;
	}
}

export class UpdateSchemaRequest {
	constructor(private readonly oldName: string, private readonly name: string) {}
	public getName(): string {
		return this.name;
	}

	public getOldName(): string {
		return this.oldName;
	}
}

export class DeleteSchemaRequest {
	constructor(private readonly name: string) {}
	public getName(): string {
		return this.name;
	}
}

export class SchemaService {
	private get BASE_URL() {
		return '/schemas';
	}

	constructor(private readonly api: ApiGateway) {}

	public async listSchemas(token: string) {
		return this.api.get<ListSchemasApiResponse>(this.BASE_URL, {token});
	}

	public async listTables(token: string, database: string) {
		return this.api.get<string[]>('/tables', {token, database});
	}

	public async describeTable(token: string, database: string, table: string) {
		return this.api.get<TableDescription>('/tables/describe', {token, database, table});
	}
}
