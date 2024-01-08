<script lang="ts">
import { useAuthStore, useSchemaStore } from '../store/store';


export default {
    methods: {
        setDatabase(to: string) {
            this.schemaStore.selectDatabase(this.authStore.mustGetCurrent, to);
        },
        selectTable(table: string) {
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.schemaStore.selectTable(token, db.name, table);
            this.$router.push("/select");
            window.scrollTo(0,0);
        },
        describeTable(table: string) {
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.schemaStore.selectTable(token, db.name, table);
            this.$router.push("/describe-table");
            window.scrollTo(0,0);
        },
    },

    setup() {
        return {
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
        }
    },
    computed:{ 
        currentConnection() {
            return this.authStore.currentConnection;
        },
        hasCurrentConnection() {
            return this.authStore.hasCurrentConnection;
        },
        databases() {
            return this.schemaStore.databasesList[this.authStore.mustGetCurrent];
        },
        currentDatabase: {
            get() {
                return this.schemaStore.currentDatabase[this.authStore.mustGetCurrent];
            },
            set(to: string) {
                this.schemaStore.selectDatabase(this.authStore.mustGetCurrent, to);
                this.schemaStore.loadTables(this.authStore.mustGetCurrent, to);
            }
        }
    },
}
</script>

<template>
    <div class="min-vh-100 bg-dark w-100 border-right-1" id="sidebar">
        <h1 class="text-center w-100 py-3"><span @click="$router.push('/')" style="cursor: pointer">DB-GUI</span></h1>


        <!-- SQL Editor button -->
        <div class="px-2">
            <template v-if="currentConnection !== undefined">
                <div class="btn-group mb-1 w-100">
                    <router-link class="btn btn-primary" to="/sql-editor">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                            <path d="M27 15 L27 30 2 30 2 5 17 5 M30 6 L26 2 9 19 7 25 13 23 Z M22 6 L26 10 Z M9 19 L13 23 Z" />
                        </svg>
                        SQL Editor
                    </router-link>
                    <router-link class="btn btn-primary" to="/import">
                        <svg id="i-import" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                            <path d="M28 22 L28 30 4 30 4 22 M16 4 L16 24 M8 16 L16 24 24 16" />
                        </svg>
                        Import
                    </router-link>
                </div>
                <div class="btn-group w-100">
                    <router-link class="btn btn-primary" to="/create-table">
                        <svg id="i-plus" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                            <path d="M16 2 L16 30 M2 16 L30 16" />
                        </svg>
                        Create table
                    </router-link>
                    <router-link class="btn btn-primary" to="/export">
                        <svg id="i-export" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                            <path d="M28 22 L28 30 4 30 4 22 M16 4 L16 24 M8 12 L16 4 24 12" />
                        </svg>
                        Export
                    </router-link>
                </div>

            </template>
            <template v-else="">
                <div class="btn-group mb-1 w-100">
                    <button class="btn btn-disabled" disabled>
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                            <path d="M27 15 L27 30 2 30 2 5 17 5 M30 6 L26 2 9 19 7 25 13 23 Z M22 6 L26 10 Z M9 19 L13 23 Z" />
                        </svg>
                        SQL Editor
                    </button>
                    <button class="btn btn-disabled" disabled>
                        <svg id="i-import" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                            <path d="M28 22 L28 30 4 30 4 22 M16 4 L16 24 M8 16 L16 24 24 16" />
                        </svg>
                        Import
                    </button>
                </div>
                <div class="btn-group w-100">
                    <button class="btn btn-disabled" disabled>
                        <svg id="i-plus" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                            <path d="M16 2 L16 30 M2 16 L30 16" />
                        </svg>
                        Create table
                    </button>
                    <button class="btn btn-disabled" disabled>
                        <svg id="i-export" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                            <path d="M28 22 L28 30 4 30 4 22 M16 4 L16 24 M8 12 L16 4 24 12" />
                        </svg>
                        Export
                    </button>
                </div>

            </template>
        </div>
        <hr>
        <div class="px-2">
            <div class="input-group">
            <span class="input-group-text">DB: </span>
            <select  v-model="currentDatabase" v-if="hasCurrentConnection" class="form-select">
                <option></option>
                <option v-for="database in schemaStore.databasesList[authStore.mustGetCurrent]">{{database.name}}</option>
            </select>
            <select class="disabled form-control" disabled v-else></select>
        </div>
        </div>
        <hr>
        <ul v-if="authStore.hasCurrentConnection" class="list-group overflow-hidden list-group-flush" style="">
            <template v-if="schemaStore.hasCurrentDatabase(authStore.mustGetCurrent)">
            
            <li v-for="(table, num) in schemaStore.mustGetCurrentDatabase(authStore.mustGetCurrent).tableNames" class="list-group-item list-group-item-action p-0 px-2 border-0" :style="{
                cursor: 'pointer',
                height: '24px',
                'background-color': num % 2 === 0 ? 'var(--bs-black)' : 'var(--bs-gray-900)'}" @click.stop="describeTable(table)">
                <span class="nobr">
                    <span @click.stop="selectTable(table)">
                        <svg 
                            alt="Select" 
                            class="text-info"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 32 32"
                            width="21"
                            height="21" 
                            fill="none" 
                            stroke="currentcolor" 
                            stroke-width="2"
                        >
                            <path d="M4 10 L4 28 28 28 28 10 M2 4 L2 10 30 10 30 4 Z M12 15 L20 15" />
                        </svg>
                    </span>
                    <span>&nbsp;<span class="d-inline-block">{{ table }}</span></span>
                </span>
            </li>
        </template>
        </ul>
        
        
    </div>
</template>