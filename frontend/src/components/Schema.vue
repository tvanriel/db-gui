<script lang="ts">
import { useAuthStore, useSchemaStore } from '../store/store';
import BreadcrumbComponent, { BreadcrumbItem } from './BreadcrumbComponent.vue';

export default {
    components: {BreadcrumbComponent},
    data: () => ({
        breadcrumbs: [
            {
                text: "Home",
                href: "/",
                active: false,
            },
            {
                text: "Schema",
                href: "/schema",
                active: true,
            }
        ] as BreadcrumbItem[],
        colours: [
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
],
    }),

    mounted() {
        if (this.authStore.hasCurrentConnection) {
            this.schemaStore.loadDatabases(this.authStore.mustGetCurrent);
        }
    },

    methods: {
        listTables(database: string) {
            this.schemaStore.selectDatabase(this.authStore.mustGetCurrent, database);
            this.$router.push("/list-tables");
        },
    },

    computed: {
        currentConnectionNickname: {
            get() {
                return this.authStore.currentConnection?.nickname;
            },
            set(to: string) {
                this.authStore.setNickname(this.authStore.mustGetCurrent, to);
            }
        },
        currentConnectionColour: {
            get() {
                return this.authStore.currentConnection?.colour;
            },
            set(to: string) {
                this.authStore.setColour(this.authStore.mustGetCurrent, to);
            }
        },
    },

    setup() {

        return {
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
        };
    },
}
</script>
<template>
<div v-if="authStore.currentConnection !== undefined">
    <BreadcrumbComponent :items="breadcrumbs"></BreadcrumbComponent>
    <div class="row">
        <div class="col-12 col-lg-6">
            <h1>Select database</h1>
            <p>
                You're logged in as <b>{{ authStore.currentConnection.name }}</b>
            </p>
        </div>
        <div class="col-12 col-lg-6">
            <fieldset>
        <legend>Connection settings</legend>
        <div class="form-group mb-2">
            <div class="input-group">
                <span class="input-group-text">
                    Nickname
                </span>
                <input type="text" v-model="currentConnectionNickname" class="form-control">
            </div>
        </div>
        <div class="form-group mb-2">
            <div class="input-group">
                <span class="input-group-text">
                    Colour
                </span>
                <select class="form-control" v-model="currentConnectionColour">
                    <option :value="colour" v-for="colour in colours">{{ colour }}</option>
                </select>
            </div>
        </div>
    </fieldset>
        </div>
    </div>

    <table class="table table-striped">
        <thead>
            <tr>
                <th>Schema name</th>
                <th>
                </th>
            </tr>
        </thead>
        <tbody v-if="schemaStore.databasesList[authStore.mustGetCurrent] !== undefined">
            <tr v-for="database in schemaStore.databasesList[authStore.mustGetCurrent]">
                <td colspan="2">
                    <button
                        class="btn btn-link"
                        @click="listTables(database.name)"
                    >{{ database.name }}</button>
                </td>
            </tr>
        </tbody>
    </table>
</div>
<div v-else="">
    <BreadcrumbComponent :items="breadcrumbs"></BreadcrumbComponent>
    <h1>No database selected</h1>
    <RouterLink to="/">Back to home</RouterLink>.
</div>
</template>