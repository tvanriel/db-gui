<script lang="ts">
import { useAuthStore } from '../store/store';


export default {
    setup() {
        return {
            authStore: useAuthStore(),
        }
    },
    computed: {

        selectedToken(): string|undefined {
            return this.authStore.current;
        },

        connections() {
            return this.authStore.getConnections;
        },

        isOnAddConnectionPage() {
            return this.$route.path.toLocaleLowerCase() === '/add-connection'
        },
    },

    methods: {
        addConnection() {
            this.$router.push('/add-connection');
        },

        selectConnection(token: string) {
            this.authStore.selectConnection(token);
            this.$router.push(`/schema`)
        }
    }
}
</script>

<template>
    <div class="w-100 pt-1">
        <div class="nav nav-tabs justify-content-between px-1">
            <div>
                <li class="nav-item d-inline-block border-0" v-for="([token, connection], id) in authStore.getConnections">
                    <button 
                    @click="selectConnection(token)" 
                    :class="{'nav-link': true, 'active bg-dark': authStore.isSelected(token)}">
                        <template v-if="connection.nickname === ''">{{connection.name}}</template>
                        <template v-else="">{{ connection.nickname }}</template>
                        <span class="text-muted">
                            [<span :class="`text-${connection.colour}`">{{ id + 1}}</span>]
                        </span>
                    </button>
                </li>
            </div>
            <div>
                <li class="nav-item d-inline-block">
                    <button
                        :class="{'nav-link border-0 float-right': true, 'bg-dark active': isOnAddConnectionPage}" @click="addConnection">
                        + Add connection
                    </button>
                </li> 
            </div>
        </div>
    </div>
</template>