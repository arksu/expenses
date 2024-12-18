<script>

import api from '../api/apiService';
import { toast } from 'vue3-toastify';

export default {
    data() {
        return {
            amount: '',
            selectedCategory: '',
            comment : '',
            categories: []
        };
    },
    mounted() {
        this.getCategories();
    },

    methods: {
        getCategories() {
            api.getCategories().then(response => {
                this.categories = response.data;
            }).catch(error => {
                console.error('Error fetching categories:', error);
            });
        },
        createExpense() {
            api.createExpense({
                amount: this.amount,
                category_id: this.selectedCategory,
                comment: this.comment,
            }).then(() => {
                this.amount = '';
                this.selectedCategory = '';
                this.comment = '';
                toast('Расход добавлен');
            }).catch(error => {
                console.error('Error creating expense:', error);
            });
        }
    }

}
</script>

<template>
    <h1>Добавить расход</h1>
    <form class="form" @submit.prevent="createExpense">
      <label for="amount">Сумма:</label>
      <input id="amount" v-model="amount" type="number" style="width: 100px;" required />

      <label for="category">Категория:</label>
      <select id="category" v-model="selectedCategory">
        <option v-for="category in categories" :key="category.id" :value="category.id">
          {{ category.name }}
        </option>
      </select>

      <label for="comment">Комментарий</label>
      <input id="comment" v-model="comment" type="text" />

      <button type="submit" :disabled="!selectedCategory || !amount || amount < 0">Добавить расход</button>
    </form>
</template>

<style>
form {
    max-width: 200px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}
</style>