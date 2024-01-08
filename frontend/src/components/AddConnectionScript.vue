<script lang="ts">
import { useAuthStore, useSchemaStore } from '../store/store';
import BreadcrumbComponent, { BreadcrumbItem } from './BreadcrumbComponent.vue';

export default {
    data: () => ({
        script: "",
        name: "",
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
            this.authStore.authenticateScript(
                this.name,
                this.script,
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
                    <router-link to="/add-connection">Back</router-link>
                    <form autocomplete="on" @submit="addConnection">
                        <label for="name" class="form-label">Name</label>
                        <input class="form-control" autofocus v-model="name" name="name" ref="server-input"/>
                        <label for="script" class="form-label">Script</label>
                        <textarea name="script" id="" cols="30" rows="10" class="form-control text-monospace" style="white-space: pre" v-model="script"></textarea>
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