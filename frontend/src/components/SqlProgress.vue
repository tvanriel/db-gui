<script lang="ts">
import { useSqlStore } from '../store/store';

export default {
    setup() {
        return {
            sqlStore: useSqlStore(),
        }
    },
    computed: {
        progress() {
            return this.sqlStore.sqlProgress;
        },
        progressPct(): string {
            if (this.progress !== undefined && this.progress.type === "sql/handler/progress_executing") {
                return `${(this.progress.executed / this.progress.total * 100).toFixed(0)}%`;
            }
            return '';
        },
        progressBytes(): string {
            if (this.progress === undefined || this.progress.type !== "sql/handler/progress_parsing") {
                return '';
            }
            const size = this.progress.progress as number;
            const i = size == 0 ? 0 : Math.floor(Math.log(size) / Math.log(1000));
            return `${(size / Math.pow(1000, i)).toFixed(0)} ${['B', 'kB', 'MB', 'GB', 'TB'][i]}`;
        },
    }
}
</script>

<template>
    <div v-if="progress !== undefined">
        <div v-if="progress.type === 'sql/handler/progress_parsing'">
            <div class="progress">
                <div class="progress-bar progress-bar-striped bg-primary text-white" style="width:100%;">
                    Parsing: {{ progressBytes }}
                </div>
            </div>
        </div>
        <div v-if="progress.type === 'sql/handler/progress_executing'">
            <div class="progress">
                <div class="progress-bar progress-bar-striped bg-success text-white" :style="{
                    width: progressPct
                }">
                    Executing: {{ progressPct }}
                </div>
            </div>
        </div>
    </div>
</template>