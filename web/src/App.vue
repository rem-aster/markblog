<template>
    <Navbar />
    <br class="m-5" />
    <RouterView />
</template>

<script setup lang="ts">
import Navbar from "@/components/Navbar.vue";
import { RouterView } from 'vue-router'
</script>
