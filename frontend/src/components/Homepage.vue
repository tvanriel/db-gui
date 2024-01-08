<script lang="ts">
import { Token } from '../store/auth/auth-store';
import { useAuthStore, useSchemaStore } from '../store/store';
import BreadcrumbComponent from './BreadcrumbComponent.vue';
export default {
    components: {
        BreadcrumbComponent,
    },
    data: () => ({
        items: [
            {
                text: "Home",
                active: true,
                href: "/",
            }
        ]
    }),

    methods: {
        selectConnection(token: Token) {
            this.authStore.selectConnection(token);
            this.$router.push("/schema");
        },
    },

    mounted() {
        if (this.authStore.hasCurrentConnection) {
            if (this.schemaStore.hasCurrentDatabase(this.authStore.mustGetCurrent)) {
                this.schemaStore.deselectDatabase(this.authStore.mustGetCurrent);
            }

        }
        this.authStore.unsetCurrent();
    },

    setup() {
        return {
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
        }
    },
}
</script>

<template>
    <div>
        <BreadcrumbComponent :items="items"></BreadcrumbComponent>
        <h1>Database GUI</h1>
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>ID</th>
                    <th colspan="5">Connection</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="([token, connection], id) in authStore.getConnections">
                    <td>
                        <span class="text-muted">[<span :class="`text-${connection.colour}`">{{ id + 1 }}</span>]</span>
                    </td>
                    <td colspan="4">
                        <button class="btn btn-link" @click="selectConnection(token)">
                            <template v-if="connection.nickname === ''">{{connection.name}}</template>
                            <template v-else="">{{ connection.nickname }} </template>
                        </button>
                    </td>
                    <td>
                        <div class="btn-group"  style="float: right;">
                            <button class="btn btn-info" v-if="connection.config !== null" @click="authStore.reconnect(token)">Reconnect</button>
                            <button class="btn btn-danger" @click="authStore.removeConnection(token)">
                                <svg id="i-trash" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="20" height="20" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
        <path d="M28 6 L6 6 8 30 24 30 26 6 4 6 M16 12 L16 24 M21 12 L20 24 M11 12 L12 24 M12 6 L13 2 19 2 20 6" />
    </svg>Delete
                            </button>
                        </div>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>