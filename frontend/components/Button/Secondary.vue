c<!-- components/Button/Primary.vue -->
<template>
  <button
    :type="type"
    :class="[
      'bg-white text-blue-500 font-semibold px-8 py-3 rounded',
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
