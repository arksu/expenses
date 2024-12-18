<script>
import api from '../api/apiService';
import { toast } from 'vue3-toastify';

export default {
  data() {
    return {
      expenses: [],
      currentPage: 0,
      pageSize: 10,
      allDataLoaded: false
    };
  },
  mounted() {
    this.fetchExpenses();
  },
  methods: {
    fetchExpenses() {
      api.getExpenses(this.currentPage, this.pageSize).then(response => {
        const newExpenses = response.data;
        if (newExpenses.length < this.pageSize) {
          this.allDataLoaded = true;
        }
        this.expenses.push(...newExpenses);
      }).catch(error => {
        console.error('Error fetching expenses:', error);
      });
    },
    loadMore() {
      this.currentPage++;
      this.fetchExpenses();
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      const options = { day: 'numeric', month: 'long', year: 'numeric' };
      return date.toLocaleDateString('ru-RU', options);
    }
  }
};
</script>

<template>
    <div>
      <h1>Expenses</h1>
      <ul>
        <li v-for="expense in expenses" :key="expense.id">
        <strong>{{ expense.amount }}</strong> - {{ expense.category }} ({{ formatDate(expense.date) }})
      </li>
      </ul>
      <button @click="loadMore" v-if="!allDataLoaded">Загрузить еще</button>
    </div>
  </template>
