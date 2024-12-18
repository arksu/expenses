<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="nav">
        <RouterLink to="/">Расход</RouterLink>
        <RouterLink to="/expense">Посл</RouterLink>
        <RouterLink to="/category">Категории</RouterLink>
        <RouterLink to="/about">?</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
</template>

<style scoped>

.nav {
  size: 12px;
}

a {
  padding-left: 5px;
  padding-right: 5px;
}

</style>
