<script lang="ts">
import { useAuthStore, useSchemaStore } from '../store/store';
import BreadcrumbComponent, { BreadcrumbItem } from './BreadcrumbComponent.vue';

export default {
    data: () => ({
        breadcrumb: [
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
                text: "List tables",
                active: true,
            }
        ] as BreadcrumbItem[]
    }),

    setup() {
        return {
            schemaStore: useSchemaStore(),
            authStore: useAuthStore(),
        }
    },

    mounted() {
        this.schemaStore.loadTables(this.authStore.mustGetCurrent, this.schemaStore.currentDatabase[this.authStore.mustGetCurrent] ??'');

    },

    methods: {
        
        selectTable(name: string) {
            this.schemaStore.selectTable(
                this.authStore.mustGetCurrent,
                this.schemaStore.currentDatabase[this.authStore.mustGetCurrent] ??'',
                name
            );
            this.$router.push("/describe-table");
        }
    },

    computed: {
        tables() {
            return this.schemaStore.databasesList[this.authStore.mustGetCurrent][this.schemaStore.currentDatabase[this.authStore.mustGetCurrent] ?? '']?.tableNames
        }
    },

    components: { BreadcrumbComponent }
}
</script>

<template>
    <div v-if="authStore.hasCurrentConnection && schemaStore.hasCurrentDatabase(authStore.mustGetCurrent)">
        <BreadcrumbComponent :items="breadcrumb">
        </BreadcrumbComponent>
        <h1>Tables in {{schemaStore.getCurrentDatabase(authStore.mustGetCurrent)?.name }}</h1>
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>Table Name</th>
                    <th>
                    </th>
                </tr>
            </thead>
            <tbody v-if="tables !== undefined">
                <tr v-for="table in tables">
                    <td colspan="2">
                        <button class="btn btn-link" @click="selectTable(table)">
                            {{ table }}
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    <div v-else="">
        <BreadcrumbComponent :items="breadcrumb"></BreadcrumbComponent>
        <h1>
            No database selected
        </h1>
    </div>
</template>
