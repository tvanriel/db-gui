<script lang="ts">
import { useAuthStore, useSchemaStore, useSelectStore } from '../store/store';

export default {

    setup() {
        return {
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
            selectStore: useSelectStore(),
        }
    },

    computed: {
        config() {
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            return this.selectStore.selectConfigFor(
                token, db.name, db.selectedTable,
            );
        },

        limit() {
            return this.config.limit.limit;
        },
        
        offset() {
            return this.config.limit.offset;
        },

        resultLength() {
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            
            const result = this.selectStore.resultSetFor(
                token, db.name, db.selectedTable,
            );

            if (result === undefined) return 0;
            return result[0].result.length;
        }
    },


    methods: {
        nextPage() {
            const to = this.offset + this.limit;
            const config = this.config;
            config.limit.offset = to;

            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);

            this.selectStore.updateLimit(
                token, db.name, db.selectedTable,
                config.limit,
            );

            this.selectStore.executeSelect(
                token, db.name, db.selectedTable,
            );
            window.scrollTo(0,0);
        },

        prevPage() {
            const to = Math.max(0, this.offset - this.limit);
            const config = this.config;
            config.limit.offset = to;

            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);

            this.selectStore.updateLimit(
                token, db.name, db.selectedTable,
                config.limit,
            );

            this.selectStore.executeSelect(
                token, db.name, db.selectedTable,
            );

            window.scrollTo(0,0);
        }
    }
}
</script>

<template>
    <ul class="pagination">
        <li class="page-item">
            <a @click="prevPage" style="cursor:pointer" :disabled="offset < 1" :class="{'page-link': true, 'disabled': offset < 1}">&lt;</a>
        </li>
        <li class="page-item">
            <a @click="nextPage" style="cursor:pointer" :disabled="resultLength < limit" :class="{'page-link': true, 'disabled': resultLength < limit}">&gt;</a>
        </li>

    </ul>
</template>
