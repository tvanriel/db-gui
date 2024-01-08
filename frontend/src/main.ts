import {createApp, watch} from 'vue';
import './style.scss';
import App from './App.vue';
import {createRouter, createWebHashHistory} from 'vue-router';
import {routes} from './routes/routes';

import {createPinia} from 'pinia';

import VueHighlightJS from 'vue3-highlightjs';
import 'highlight.js/styles/atom-one-dark.css';
import 'highlight.js/lib/languages/sql';

import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css'
import ContextMenu from '@imengyu/vue3-context-menu'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';

// @ts-ignore
self.MonacoEnvironment = {
	globalAPI: true,
	getWorker() {
		return new editorWorker();
	},
};

const app = createApp(App);

const router = createRouter({
	routes,
	history: createWebHashHistory(),
});

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

app.use(router);
app.use(pinia);
app.use(VueHighlightJS);

app
.use(ContextMenu)
.mount('#app');
