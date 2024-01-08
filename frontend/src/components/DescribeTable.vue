<script lang="ts">

import { useAuthStore, useSchemaStore } from '../store/store';
import BreadcrumbComponent, { BreadcrumbItem } from './BreadcrumbComponent.vue';
export default {
    components: {BreadcrumbComponent},
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
                active: true,
            }
        ] as BreadcrumbItem[]
    }),
    setup() {
        return {
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
        }
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
        }
    }
}
</script>
<template>
    <BreadcrumbComponent :items="breadcrumbs"></BreadcrumbComponent>
    <div v-if="table !== null">
        <h1>Table: {{table.name}}</h1>
        <RouterLink to="/select">Select data</RouterLink>
        <table class="table table-striped">
            <thead>
                <tr>
                    <th><b>Column</b></th>
                    <th>Type</th>
                    <th>Comment</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="column in table.columns">
                    <td>{{ column.name }}</td>
                    <td>
                        {{ column.type }}
                        <span v-if="column.default">
                            [<b>{{ column.default }}</b>]
                        </span>
                    </td>
                    <td>{{ column.comment }}</td>
                </tr>
            </tbody>
        </table>
    </div>
    <div v-else="">
        <h1>No table selected</h1>
    </div>
</template>
