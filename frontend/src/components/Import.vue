<script lang="ts">
import { useAuthStore, useImportStore, useSchemaStore } from '../store/store';
import Breadcrumbs, { BreadcrumbItem } from './BreadcrumbComponent.vue';


export default {
    data: () => ({
        file: null,
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
                text: "Import",
                active: true,
            }
        ] as BreadcrumbItem[]
    }),

    computed: {
        currentDatabase() {
            return this.schemaStore.getCurrentDatabase(this.authStore.mustGetCurrent)
        },
        databases() {
            return this.schemaStore.databasesList[this.authStore.mustGetCurrent]!
        },
        isImporting() {
            return this.importStore.isImporting[this.authStore.mustGetCurrent]
        }
    },

    setup() {
        return {
            schemaStore: useSchemaStore(),
            authStore: useAuthStore(),
            importStore: useImportStore(),
        }
    },

    methods: {

        setDatabase(event: Event) {
            const target: null|EventTarget = event.target;
            if (target === null) {
                return;
            }
            this.schemaStore.selectDatabase(this.authStore.mustGetCurrent, (target as HTMLSelectElement).value);
        },

        submit() {
            const fd = new FormData();
            const db = this.schemaStore.getCurrentDatabase(this.authStore.mustGetCurrent);
            if (db === undefined) {
                this.importStore.addError(this.authStore.mustGetCurrent, "No database selected");
                return;
            }
            fd.append("databaseName", db.name);
            const files = (this.$refs.file as HTMLInputElement).files;

            if (files === null) return;
            if (files.length === 0) return;


            fd.append("file", files[0]);

            this.importStore.runImport(this.authStore.mustGetCurrent, fd);
        }
    },
    
    components: {Breadcrumbs},
}
</script>

<template>
    <div v-if="authStore.hasCurrentConnection">
        <Breadcrumbs :items="breadcrumbs"></Breadcrumbs>
        <h3>Import file into database</h3>
        <div class="alert alert-info"  v-if="importStore.isImporting[authStore.mustGetCurrent] !== undefined">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Executing</span>
            </div>
            Import is running, please wait.
        </div>
        <div v-if="importStore.results[authStore.mustGetCurrent] !== undefined">
            <div class="alert alert-danger mb-3" v-for="error in importStore.results[authStore.mustGetCurrent].errors">
                {{ error.error }}
                <pre v-highlightjs class="mb-0" v-if="error.sql"><code class="mysql">{{ error.sql }}</code></pre>            
            </div>

            <div class="alert alert-success">
                Executed {{ importStore.results[authStore.mustGetCurrent].executed }} queries.
            </div>
        </div>

        <div class="input-group">
            <span class="input-group-text">DB: </span>
            <select @change="setDatabase($event)" :value="currentDatabase?.name" class="form-select">
                <option></option>
                <option v-for="database in databases" :value="database.name">{{database.name}}</option>
            </select>
        </div>

        <label for="formFile" class="form-label">Import file</label>
        <input class="form-control" accept=".sql.gz,.sql" type="file" ref="file">
        <div class="my-3">
            <button :class="{
                'btn btn-primary': true,
                'disabled': isImporting
            }" :disabled="isImporting" @click="submit">Submit</button>
        </div>
    </div>

</template>
