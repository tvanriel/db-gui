import {defineStore } from 'pinia';
import {doResponse} from '../../lib/api-gateway';
import {type AuthService} from '../../services/auth-service';

export type Token = string;

export type Connection = {
	errors: [];
	name: string;
	colour: string;
	config?: {
		type: 'upw';
		username: string;
		host: string;
		password: string;
	} | {
		type: 'script';
		script: string;
	};
	nickname: string;
};

const colour = (token: string): string => {
	const colours = [
		'red',
		'orange',
		'yellow',
		'green',
		'teal',
		'cyan',
		'blue',
		'indigo',
		'purple',
		'pink',
	];

	let hash = 0;
	for (let i = 0; i < token.length; i++) {
		hash += token.charCodeAt(i) * 3;
	}

	return colours[hash % colours.length];
};

export const makeAuthStore = (authService: AuthService) => defineStore('auth', {
	state: () => ({
		tokens: {} as Record<Token, Connection>,
		current: undefined as string | undefined,
		errors: [] as string[],

		isAuthenticating: false,
	}),

	persist: {
		paths: [
			'tokens',
			'current',
			'errors',
		],
	},

	getters: {
		currentConnection(state): Connection | undefined {
			return this.getConnection(state.current ?? '');
		},

		hasCurrentConnection: state => state.current !== undefined,

		mustGetCurrent: state => state.current!,

		getConnection: state => (token: Token): Connection | undefined => state.tokens[token],

		getTokens(): Token[] {
			return [...Object.keys(this.tokens)];
		},
		getConnections(): Array<[Token, Connection]> {
			return this.getTokens.map((token): [Token, Connection] => [token, this.getConnection(token)!]).sort((a, b) => {
				let aname: string = a[1].name;
				let bname: string = b[1].name;
				if (a[1].nickname !== '') {
					aname = a[1].nickname;
				}

				if (b[1].nickname !== '') {
					bname = b[1].nickname;
				}

				if (aname < bname) {
					return -1;
				}

				if (aname > bname) {
					return 1;
				}

				return 0;
			});
		},
		isSelected: state => (token: Token) => state.current === token,
	},

	actions: {
		_finishAuthenticating() {
			this.isAuthenticating = false;
		},

		/**
         * Log in to a server and save the connection in the store.
         *
         * @param host Hostname of the server you want to connect to
         * @param username Username to authenticate with
         * @param password PAssword to authenticate with
         */
		async authenticate(host: string, username: string, password: string) {
			this.isAuthenticating = true;
			return authService.logIn(host, username, password).then(doResponse(
				response => {
					this.isAuthenticating = false;
					this.tokens[response.token] = {
						errors: [],
						name: `${username}@${host}`,
						colour: colour(`${username}@${host}`),
						config: {
							type: 'upw',
							password,
							host,
							username,
						},
						nickname: '',
					};
					this.current = response.token;
				},
				error => {
					this.addError(error.error());
				},
			));
		},
		async authenticateScript(name: string, script: string) {
				this.isAuthenticating = true;
				return authService.logInScript(script).then(doResponse(
				response => {
					this.isAuthenticating = false;
					this.tokens[response.token] = {
						errors: [],
						name: name,
						colour: colour(name),
						config: {
							type: 'script',
							script,
						},
						nickname: '',
					};
					this.current = response.token;
				},
				error => {
					this.isAuthenticating = false;
					this.addError(error.error());
				},
			));
		},

		selectConnection(t: Token) {
			this.current = t;
		},

		unsetCurrent() {
			this.current = undefined;
		},

		dismissError(i: number) {
			this.errors.splice(i, 1);
		},
		dismissAllErrors() {
			this.errors = [];
		},

		addError(err: string) {
			this.errors.push(err);
		},

		removeConnection(token: Token): void {
			delete this.tokens[token];
		},

		reconnect(token: Token) {
			const conn = this.tokens[token];
			if (conn === undefined) {
				return;
			}

			const {config} = conn;
			if (config === undefined) {
				return;
			}

			this.removeConnection(token);
			if (config.type === 'upw') {
				return this.authenticate(config.host, config.username, config.password).catch(() => this.tokens[token] = conn).then(() => {
					const currentToken = this.current;
					if (currentToken === undefined) {
						return;
					}

					this.tokens[currentToken].nickname = conn.nickname;
					this.tokens[currentToken].colour = conn.colour;
				});
			}

			return this.authenticateScript(conn.name, config.script).catch(() => this.tokens[token] = conn).then(() => {
				const currentToken = this.current;
				if (currentToken === undefined) {
					return;
				}

				this.tokens[currentToken].nickname = conn.nickname;
				this.tokens[currentToken].colour = conn.colour;
			});
		},

		setNickname(token: Token, to: string) {
			this.tokens[token].nickname = to;
		},
		setColour(token: Token, colour: string) {
			this.tokens[token].colour = colour;
		},
	},
});
