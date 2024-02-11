<template>
  <div id="app">
    <header>
      <date-time class="datetime" />
      <button onclick="location.reload()">reload</button>
    </header>
    <router-view/>
  </div>
</template>

<script>
import DateTime from '@/components/DateTime'
export default {
  components: {
    DateTime
  }
}
</script>
<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

.datetime {
  display: inline-block;
  padding: 3px 20px 0 0;
  height: 100%;
  font-size: 13px;
  color: #434446;
  vertical-align: text-bottom;
}

</style>
