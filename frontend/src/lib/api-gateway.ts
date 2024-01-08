import axios, {type AxiosRequestConfig, type AxiosResponse} from 'axios';

export class ApiError {
	constructor(private readonly value: string) {}
	public error(): string {
		return this.value;
	}
}

export type Result<T> = SuccessResult<T> | ErrorResult;
export type ErrorResult = [undefined, ApiError];
export type SuccessResult<T> = [T, undefined];

export class ApiGateway {
	constructor(private readonly baseUrl: string) {

	}

	public getBaseUrl(): string {
		return this.baseUrl;
	}

	public async rawGet<T>(url: string, options: AxiosRequestConfig): Promise<AxiosResponse<T>> {
		return axios.get(url, options);
	}

	public async rawPost<T>(url: string, data: any, options: AxiosRequestConfig): Promise<AxiosResponse<T>> {
		return axios.post<T>(url, data, options);
	}

	public async rawPatch<T>(url: string, data: any, options: AxiosRequestConfig): Promise<AxiosResponse<T>> {
		return axios.patch<T>(url, data, options);
	}

	public async rawDelete<T>(url: string, options: AxiosRequestConfig): Promise<AxiosResponse<T>> {
		return axios.delete<T>(url, options);
	}

	public async get<T>(to: string, queryParams?: any): Promise<Result<T>> {
		return this.decoratePromise<T>(this.rawGet<T>(this.baseUrl + to, {
			params: queryParams ?? {},
		}));
	}

	public async post<T>(to: string, body: Record<string, unknown> | FormData, queryParams?: any): Promise<Result<T>> {
		return this.decoratePromise<T>(this.rawPost<T>(this.baseUrl + to, body, {
			params: queryParams ?? {},
		}));
	}

	public async patch<T>(to: string, body: Record<string, unknown>, queryParams?: any): Promise<Result<T>> {
		return this.decoratePromise<T>(this.rawPatch(this.baseUrl + to, body, {
			params: queryParams ?? {},
		}));
	}

	public async delete<T>(to: string, queryParams?: any): Promise<Result<T>> {
		return this.decoratePromise(this.rawDelete(this.baseUrl + to, {
			params: queryParams ?? {},
		}));
	}

	public async ws(to: string, queryParams?: Map<string, string>): Promise<WebSocket> {
		const url = new URL(`ws://${location.host}${this.baseUrl}${to}`);

		if (queryParams !== undefined) {
			queryParams.forEach((v, k) => {
				url.searchParams.append(k, v);
			});
		}

		return new Promise<WebSocket>((resolve, reject) => {
			try {
				const ws = new WebSocket(url, []);

				ws.addEventListener('open', () => {
					resolve(ws);
				});
			} catch (e) {
				reject(e);
			}
		});
	}

	private async decoratePromise<T>(p: Promise<AxiosResponse<T>>): Promise<Result<T>> {
		try {
			const response = await p;
			return [response.data, undefined];
		} catch (e: any) {
			if (typeof e.response !== 'undefined') {
				if (typeof e.response.data === 'object') {
					return [undefined, new ApiError(e.response.data.message)];
				}
			}

			return [undefined, new ApiError(e.message)];
		}
	}
}

type Mapper<T, U> = (response: T, index: number) => U;
type Perform<U> = (item: U) => void;
type ErrHandler = (item: ApiError) => void;

/**
 * Map a successful response from one to the other.
 *
 * This function behaves as the identity if result is not successful.
 *
 * @param f Map function to pass every element to.
 * @returns
 */
export const mapResponse = <T, U>(f: Mapper<T, U>, e?: ErrHandler) =>
	(result: Result<T[]>): Result<U[]> => {
		const [response, err] = result;

		// If the result is of the Error Type, we just return the complete result
		if (err !== undefined) {
			if (e !== undefined) {
				e(err);
			}

			return result;
		}

		// Else we map over the items in the response and return them.
		return [response.map((item, index) => f(item, index)), undefined];
	};

/**
 * Perform an action on a successful result.
 *
 * This function behaves as identity if the result is not successful
 *
 * @param p Action to perform on the successful result
 * @returns Result<null>
 */
export const doResponse = <U>(p: Perform<U>, e?: ErrHandler) =>
	(result: Result<U>): Result<null> => {
		const [item, err] = result;

		if (err !== undefined) {
			if (e !== undefined) {
				e(err);
			}

			return [undefined, err];
		}

		p(item);
		return [null, undefined] as SuccessResult<null>;
	};

export const mapAndDo = <T, U>(f: Mapper<T, U>, p: Perform<U[]>, e?: ErrHandler) =>
	(result: Result<T[]>): Result<null> => doResponse(p, e)(mapResponse(f, e)(result));
