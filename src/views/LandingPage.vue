<template>
  <div class="landing-page">
    <div class="middle-strip">
      <SearchForm></SearchForm>
    </div>
  </div>
</template>

<script>
import SearchForm from '../components/SearchForm.vue';
export default {
  components: {
  SearchForm,
},
}
</script>
  
<style scoped>
.landing-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
}

.middle-strip {
  height: 15vh;
  width: 100vw;
  background-color: rgb(28, 48, 79);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>