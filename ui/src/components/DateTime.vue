<template>
  <div>
    {{ datetime }}
  </div>
</template>

<script>
export default {
  name: 'DateTime',
  data () {
    return {
      datetime: '',
      timeInterval: null
    }
  },
  created () {
    this.timeInterval = setInterval(() => {
      this.datetime = new Date().toLocaleString('ko-KR')
    }, 1000)
  },
  destroyed () {
    clearInterval(this.timeInterval)
  }
}
</script>

<style lang="scss" scoped>
  .premium-container {
    padding: 6px 6px;
  }
</style>
