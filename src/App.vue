<template>
      <router-view :key="$route.fullPath"/>
  </template>
  
  <script>  
  export default {
  };
  </script>
  
  <style>
  body {
    margin: 0;
    padding: 0;
    overflow-x: hidden;
    font-size: 17px;
  }
  </style>