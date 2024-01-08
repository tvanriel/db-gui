<script lang="ts">
import { SqlHandlerStatementPlanMessage } from '../services/sql-service';
import { useSqlStore } from '../store/store';
import SqlProgress from './SqlProgress.vue';

export default {
    data: () => ({
        time: 0,
        interval: 0,
    }),
    mounted() {
        this.interval = setInterval(() => {
            if (this.plan === undefined) {
                this.time = 0;
            } else {
                this.time = Date.now() - this.plan.time;
            }
        }, 100);
    },
    unmounted() {
        if (this.interval !== undefined) {
            clearInterval(this.interval);
        }
    },

    setup() {
        return {
            sqlStore: useSqlStore(),
        }
    },
    computed: {
        plan() {
            return this.sqlStore.sqlCurrent
        },
        sqlProgress() {
            return this.sqlStore.sqlProgress
        },
        timeFmt(): string {
            const result = this.plan;

            if (result === null) {
                return "";
            }
            const i = this.time == 0 ? 0 : Math.floor(Math.log(this.time) / Math.log(1000));
            return `${(this.time / Math.pow(1000, i)).toFixed(2)} ${["ms", "s"][i]}`;
        }
    },
    components: { SqlProgress }
}
</script>


<template>
            <div class="card bg-dark mb-4" v-if="plan !== undefined">
            <h3 class="card-header">Executing Query</h3>
            <div class="card-body">
                <pre v-highlightjs><code language="mysql">{{ plan.sql }}</code></pre>
                <div class="spinner-border" role="status">
                    <span class="visually-hidden">Executing</span>
                </div>


                <p v-if="sqlProgress !== null">
                    <SqlProgress></SqlProgress>
                </p>
                <small>
                    ({{ timeFmt }})
                </small>
            </div>
        </div>
</template>