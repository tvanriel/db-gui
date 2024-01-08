<script lang="ts">
import { PropType } from 'vue';
import { SqlResult } from '../../services/sql-service';

export default {
    props: {
        result: {
            type: Array,
            validator: (list: SqlResult) => {
                return list.filter((item) => {
                    return Array.isArray(item.result) && Array.isArray(item.columns) && typeof item.error === "string"
                })
            },
            required: true,
        } as unknown as PropType<SqlResult>,

    }
}
</script>

<template>
    <div>
        <div v-for="res in result">
            <span>{{ res?.result?.length }} Rows</span>
            <table class="table table-striped table-sm table-hover" style="width: initial;">
                <thead>
                    <tr>
                        <th v-for="col in res.columns">
                            <p class="mb-0">
                                {{ col.name }}
                            </p>
                            <span class="text-muted">
                                [{{ col.type }}]
                            </span>
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <template v-if="res.result === null || res.result.length < 1">
                        <tr>
                            <td :colspan="res.columns.length">
                                No rows
                            </td>
                        </tr>
                    </template>
                    <tr v-for="row in res.result">
                        <td v-for="item in row">
                            <i class="text-muted" v-if="item.null">NULL</i>
                            <span v-else>{{ item.value }}</span>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>