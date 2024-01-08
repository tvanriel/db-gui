<script lang="ts">
import { SqlResult } from '../services/sql-service';
import { useAuthStore, useSchemaStore, useSelectStore } from '../store/store';
import BreadcrumbComponent from './BreadcrumbComponent.vue';
import { BreadcrumbItem } from './BreadcrumbComponent.vue';
import ResultSet from './ResultSets/SimpleResultSet.vue';
import SelectForm from './SelectForm.vue';
import SelectPager from './SelectPager.vue';
export default {
    data: () => ({
        breadcrumbs: [
            {
                text: 'Home',
                href: "/",
                active: false,
            },
            {
                text: "Schema",
                href: "/schema",
                active: false,
            },
            {
                text: "List tables",
                href: "/list-tables",
                active: false,
            },
            {
                text: "Table",
                href: '/describe-table',
                active: false,
            },
            {
                text: "Select",
                href: '/select',
                active: true,
            }
        ] as BreadcrumbItem[]
    }),

    setup() {
        return {
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
            selectStore: useSelectStore(),
        }
    },

    mounted() {
        const db = this.schemaStore.mustGetCurrentDatabase(this.authStore.mustGetCurrent);
        this.selectStore.executeSelect(this.authStore.mustGetCurrent, db.name, db.selectedTable);
    },

    computed: {
        table() {
            const token = this.authStore.current
            
            if (token === undefined) {
                return null;
            }
            
            const db = this.schemaStore.getCurrentDatabase(token);
            
            if (db === undefined) {
                return null;
            }

            return db.tables[db.selectedTable] ?? null;
        },

        result() {
            if (!this.authStore.hasCurrentConnection && !this.schemaStore.hasCurrentDatabase(this.authStore.mustGetCurrent)) {
                return undefined;
            }
            const db = this.schemaStore.getCurrentDatabase(this.authStore.mustGetCurrent);
            if (db === undefined) return undefined;

            return this.selectStore.resultSetFor(this.authStore.mustGetCurrent, db.name, db.selectedTable);
        },

        selectSql(): string | null {
            if (!this.authStore.hasCurrentConnection && !this.schemaStore.hasCurrentDatabase(this.authStore.mustGetCurrent)) {
                return null;
            }
            const db = this.schemaStore.getCurrentDatabase(this.authStore.mustGetCurrent);
            if (db === undefined) return null;

            return this.selectStore.selectSqlFor(this.authStore.mustGetCurrent, db.name, db.selectedTable) ?? null;
        }

    },
    components: { ResultSet, BreadcrumbComponent, SelectForm, SelectPager }
}
</script>

<template>

    <div v-if="table !== null">
        <BreadcrumbComponent :items="breadcrumbs"></BreadcrumbComponent>
        <h1>Select: {{ table.name }}</h1>
        <router-link to="/describe-table">Table Structure</router-link>
        <SelectForm></SelectForm>        
        <pre v-highlightjs class="mb-0"><code class="mysql">{{ selectSql }}</code></pre>
        <SelectPager></SelectPager>
        <ResultSet v-if="result !== undefined" :result="result"></ResultSet>
        <SelectPager></SelectPager>
    </div>
    <div v-else="">
        <BreadcrumbComponent :items="breadcrumbs"></BreadcrumbComponent>
        <h1>No table selected</h1>
    </div>
</template>