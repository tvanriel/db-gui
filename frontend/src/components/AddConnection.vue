<script lang="ts">
import { useAuthStore, useSchemaStore } from '../store/store';
import BreadcrumbComponent, { BreadcrumbItem } from './BreadcrumbComponent.vue';

export default {
    data: () => ({
        server: "",
        username: "",
        password: "",
        breadcrumbs: [
            {
                text: 'Home',
                href: '/',
                active: false,
            },
            {
                text: "Add connection",
                href: "/add-connection",
                active: true,
            }
        ] as BreadcrumbItem[]
    }),

    mounted() {
        (this.$refs['server-input'] as HTMLInputElement).focus();
        this.authStore.unsetCurrent();
    },

    setup() {
        return {
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
        }
    },

    methods: {
        addConnection($event: Event) {
            $event.preventDefault();
            this.authStore.dismissAllErrors();
            this.authStore.authenticate(
                this.server,
                this.username,
                this.password,
            ).then(() => {
                if (this.authStore.errors.length === 0) {
                    this.$router.push('/schema');
                    this.schemaStore.loadDatabases(this.authStore.mustGetCurrent);
                }
            })
        }
    },
    
    components: { BreadcrumbComponent }
}
</script>

<template>
    <div class="row">
        <div class="col-sm-12">
            <BreadcrumbComponent :items="breadcrumbs"></BreadcrumbComponent>
        </div>
        <div class="col-sm-12 col-xl-6">
            <div class="card">
                <h5 class="card-header">
                    Add connection
                </h5>
                <div class="card-body">
                    <router-link to="/add-connection-script">Advanced authentication mode</router-link>
                    <form autocomplete="on" @submit="addConnection">
                        <label for="server" class="form-label">Server</label>
                        <input class="form-control" autofocus v-model="server" name="server" ref="server-input"/>
                        
                        <label for="username" class="form-label">Username</label>
                        <input class="form-control" v-model="username" name="username"/>
                        
                        <label for="password" class="form-label">Password</label>
                        <input class="form-control" v-model="password" name="password" type="password"/>
                        <div class="pt-4">
                            <button class="btn btn-primary" :disabled="authStore.isAuthenticating" @click="addConnection">
                                + Add Connection
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
    
</template>