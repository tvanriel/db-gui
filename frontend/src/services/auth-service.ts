import {type ApiGateway} from '../lib/api-gateway';

export type SessionApiResponse = {
	token: string;
};

export class AuthService {
	constructor(private readonly api: ApiGateway) {}
	public async logIn(host: string, username: string, password: string) {
		return this.api.post<SessionApiResponse>('/sessions', {host, username, password});
	}

	public async logInScript(script: string) {
		return this.api.post<SessionApiResponse>('/sessions/script', {script});
	}
}
