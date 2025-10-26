<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
// App.vue now uses Vue Router for navigation
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  background: #f8fafc;
}

#app {
  min-height: 100vh;
}
</style>