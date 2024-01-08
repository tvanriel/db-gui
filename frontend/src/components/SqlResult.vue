<script lang="ts">
import { SqlHandlerErrorMessage, SqlHandlerExecuteResultMessage, SqlHandlerQueryResultMessage } from '../services/sql-service';
import { useSqlStore } from '../store/store';
import ResultSet from './ResultSets/SimpleResultSet.vue';
export default {
    props: {
        i: {
            type: Number,
            required: true,
        },
    },
    components: {ResultSet},
    setup() {
        return {
            sqlStore: useSqlStore(),
        }
    },
    computed: {
        result() {
            const result = this.sqlStore.getSqlResult[this.i];
            if (result !== null) {
                return result;
            }
            return null;
        },
        time(): string {
            const result = this.result;
            if (result === null) return '';
            if (result.type === "sql/handler/error") {return "";}
            if (result.type === "sql/handler/statement_error") {return "";}
            if (result.type === "sql/handler/plan") {return "";}
            const size = result.microseconds as number;
            const i = size == 0 ? 0 : Math.floor(Math.log(size) / Math.log(1000));
            return `${(size / Math.pow(1000, i)).toFixed(2)} ${['μs', 'ms', 's'][i]}`;
        }
    }
}
</script>

<template>
    <template v-if="result !== null || result === undefined">
    <div v-if="result.type === 'sql/handler/query'">
        <pre v-highlightjs><code language="mysql">{{ result.sql }}</code></pre>
        <ResultSet :result="result.results" v-if="result.results !== undefined"></ResultSet>
        <small>({{ time }})</small>
    </div>
    <div v-if="result.type === 'sql/handler/exec'">
        <pre v-highlightjs><code language="mysql">{{ result.sql }}</code></pre>
        <small>({{ time }})</small>
        <p>
            Rows affected: {{ result.rowsAffected }}
        </p>
        <p v-if="result.lastInsertId !== 0">
            Last insert ID: {{ result.lastInsertId }}
        </p>
    </div>
    <div v-if="result.type === 'sql/handler/error'">
        <div class="alert alert-danger">
            {{ result.message }}
        </div>
    </div>

    <div v-if="result.type === 'sql/handler/statement_error'">        
        <pre v-highlightjs><code language="mysql">{{ result.params.sql }}</code></pre>
        <div class="alert alert-danger">
            {{ result.message }}
        </div>
    </div>
</template>
</template>