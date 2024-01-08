<script lang="ts">
import BreadcrumbComponent, { BreadcrumbItem } from './BreadcrumbComponent.vue';
import * as monaco from 'monaco-editor'
import ResultSet from './ResultSets/SimpleResultSet.vue';
import SqlProgress from './SqlProgress.vue';
import SqlResult from './SqlResult.vue';
import SqlPlan from './SqlPlan.vue';
import OneDarkTheme from './monaco-editor-theme/onedark/OnedarkTextmate.tmtheme'
import { useAuthStore, useSchemaStore, useSqlStore } from '../store/store';

export default {
    data: () => ({
        options: {
            colorDecorators: true,
            lineHeight: 16,
            tabSize: 8,
        },
        resizeListener: null as null|void|(() => void),
        ready: false,
        sql: '',
        changeCallback: null as null|(() => void),
        breadcrumbs: [
            {
                text: "Home",
                active: false,
                href: "/",
            },
            {
                text: "Schema",
                active: false,
                href: "/schema"
            },
            {
                text: "SQL Editor",
                active: true,
                href: "/sql-editor"
            }
        ] as BreadcrumbItem[],
    }),
    components:
    {
        BreadcrumbComponent,
        ResultSet,
        SqlProgress,
        SqlResult,
        SqlPlan,
    },
    mounted() {
        const container = this.$refs['monaco-container'] as HTMLDivElement;

        const page = this.$refs.page as HTMLDivElement|null;
            if (page !== null) {
                const width = Math.max(page.clientWidth, 500);
                container.style.width = `${width}px`;

            }

        container.style.height = "500px"

        // monaco.languages.registerCompletionItemProvider('sql', sqlCompletion);
        // @ts-ignore
        monaco.editor.defineTheme('one-dark', OneDarkTheme);
        const editor = monaco.editor.create(container, {
            tabSize: 8,
            language: 'sql',
            theme: 'one-dark',
            minimap: {enabled: false},
            suggestOnTriggerCharacters: true,
            value: this.sql,
            fontFamily: 'var(--bs-font-monospace)',
            fontSize: 14,
        });


        this.resizeListener = window.addEventListener('resize', () => {
            const page = this.$refs.page as HTMLDivElement|null;
            if (page === null) return;
            const width = Math.max(page.clientWidth, 500);
            editor.layout({width: width, height: 500});
        });

        this.changeCallback = () => {
            editor.setValue(this.sql);
        }
        editor.onDidChangeCursorSelection(() => {
            if (editor === null) return;
            const model = editor.getModel();
            const selection = editor.getSelection();
            if (model === null || selection === null) return;

            const selected = model.getValueInRange(selection);
            if (selected.length > 0) {
                this.sql = selected;
            } else {
                this.sql = editor.getValue();
            }
        });

        this.ready = true;
    },
    unmounted() {
        const x = this.resizeListener;
        if (typeof x !== "function") {
            return;
        }
        window.removeEventListener('resize', x);
        this.resizeListener = null;
    },
    methods: {

        changeSql(to: string) {
            this.sql = to;
            if (this.changeCallback !== null)
            this.changeCallback();
        },

        executeSql() {
            
            this.sqlStore.executeSql(this.authStore.mustGetCurrent, this.schemaStore.getCurrentDatabase(this.authStore.mustGetCurrent)?.name ?? '', this.sql, (_: number) => {return new Promise<boolean>((resolve) => resolve(true))});
        }
    },
    setup() {
        return {
            sqlStore: useSqlStore(),
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
        }
    },
    computed: {
        sqlResult() {
            return this.sqlStore.getSqlResult
        },
        sqlHistory() {
            return this.sqlStore.sqlHistory.filter(item => {
                return this.authStore.mustGetCurrent === item.token;
            });
        }

    },
}
</script>
<template>

    <div ref="page">
        <BreadcrumbComponent :items="breadcrumbs"></BreadcrumbComponent>
        <h1>SQL Editor</h1>
        <div class="card bg-dark mb-4" v-for="(_, i) in sqlResult">
            <h3 class="card-header">Query #{{ i+1 }}</h3>
            <div class="card-body">
                <SqlResult :i="i"></SqlResult>
            </div>
        </div>
        <SqlPlan></SqlPlan>

        <div class="my-2" ref="monaco-container">
        </div>
        <button class="btn btn-success" @click="executeSql">
            <svg id="i-lightning" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                <path d="M18 13 L26 2 8 13 14 19 6 30 24 19 Z" />
            </svg>
            Execute
        </button>
        <ul class="list-group">
            <li  v-for="item in sqlHistory" class="list-group-item list-group-item-action" style="cursor:pointer" @click="changeSql(item.sql)">
                <pre v-highlightjs class="mb-0"><code class="mysql">{{ item.sql }}</code></pre>
            </li>
        </ul>
    </div>
</template>