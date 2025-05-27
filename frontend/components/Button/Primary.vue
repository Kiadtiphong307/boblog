<!-- components/Button/Primary.vue -->
<template>
    <button
      :type="type"
      :class="[
        'px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition',
        customClass,
      ]"
      @click="$emit('click', $event)"
    >
      <slot>Submit</slot>
    </button> 
  </template>
  
  <script setup>
  defineProps({
    type: {
      type: String,
      default: 'button',
    },
    customClass: {
      type: String,
      default: '',
    },
  })
  </script>
  