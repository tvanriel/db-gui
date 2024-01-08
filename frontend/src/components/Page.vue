<script lang="ts">
import { useAuthStore } from '../store/store';
export default {
    setup() {
        const authStore = useAuthStore();
        return {authStore};
    },
}
</script>
<template>
    <div class="w-100" id="page">
       <div class="container-fluid">
            <div class="row">
                <div class="col-12 pt-3">
                    <div class="alert alert-dismissible alert-danger" v-for="(error, i) in authStore.errors">
                        <button type="button" class="btn-close" data-bs-dismiss="alert" @click="authStore.dismissError(i)"></button>
                        <p style="white-space: pre;" class="mb-0">{{error}}</p>
                    </div>
                    <router-view></router-view>
                </div>
            </div>
       </div>
    </div>
</template>