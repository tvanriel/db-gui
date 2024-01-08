<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'
import Page from './components/Page.vue'
import PageHeader from './components/PageHeader.vue'
import { watch } from 'vue';

</script>

<template>
  <div class="container-fluid min-vh-100 px-0">
    <div class="row gx-0">
      <div class="col-2">
        <Sidebar></Sidebar>
      </div>
      <div class="col-10">
        <PageHeader></PageHeader>
        <Page></Page>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
