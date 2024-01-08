<script lang="ts">
import { useAuthStore, useExportStore, useSchemaStore } from '../store/store';
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
                text: "Export",
                active: true,
            }
        ] as BreadcrumbItem[]
    }),

    setup() {
        return {
            authStore: useAuthStore(),
            exportStore: useExportStore(),
            schemaStore: useSchemaStore(),
        }
    },

    computed: {
    },

    methods: {

        setDatabase(event: Event) {
            const to = (event.target as HTMLSelectElement).value;
            this.exportStore.selectDatabase(this.authStore.mustGetCurrent, to);
        },

        submit(event: MouseEvent) {
            event.preventDefault();
            this.exportStore.startExport(this.authStore.current ??'');
        }
    },

    components: {Breadcrumbs},
}

</script>

<template>
    <div>
        <Breadcrumbs :items="breadcrumbs"></Breadcrumbs>
        <h3>Export database</h3>

        <div class="input-group">
            <span class="input-group-text">DB: </span>
            <select @change="setDatabase($event)" :value="exportStore.selectedDatabase.get(authStore.mustGetCurrent)" class="form-select">
                <option></option>
                <option v-for="database in schemaStore.databasesList[authStore.mustGetCurrent]">
                    {{ database.name }}
                </option>
            </select>
        </div>

        <button :class="{
                'btn btn-primary': true,
            }" @click="submit">Submit</button>
    </div>
</template>