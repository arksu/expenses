<script>
import api from '../api/apiService';

export default {
  data() {
    return {
      categories: []
    };
  },
  mounted() {
    this.fetchCategories();
  },
  methods: {
    fetchCategories() {
      api.getCategories().then(response => {
        this.categories = response.data;
      }).catch(error => {
        console.error('Error fetching categories:', error);
      });
    }   
  }
};
</script>

<template>
    <div>
      <h1>Categories</h1>
      <ul>
        <li v-for="category in categories" :key="category.id">{{ category.name }}</li>
      </ul>
    </div>
  </template>