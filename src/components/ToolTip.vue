<!-- Tooltip.vue -->
<template>
    <div class="tooltip" :style="{ top: position.y + 'px', left: position.x + 'px' }" v-if="showTooltip">
      {{ text }}
    </div>
  </template>
  
  <script>
  export default {
    props: {
      text: String,
      showTooltip: Boolean,
      position: Object,
    },
  };
  </script>
  
  <style scoped>
  .tooltip {
    position: fixed;
    background-color: white;
    border: 1px solid black;
    padding: 5px;
    z-index: 1;
    white-space: normal;
  }
  </style>
  