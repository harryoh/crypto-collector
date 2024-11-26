<template>
  <div class="premium-container">
    <currency-rate />
    <bybit-price />
    <upbit-price />
    <bithumb-price />
    <!-- <alarm-rule /> -->
  </div>
</template>

<script>
import { getPrice } from '@/api/price'
import CurrencyRate from '@/components/CurrencyRate'
import BybitPrice from '@/components/BybitPrice'
import UpbitPrice from '@/components/UpbitPrice'
import BithumbPrice from '@/components/BithumbPrice'
// import AlarmRule from '@/components/AlarmRule'

export default {
  name: 'PremiumMonitor',
  components: { CurrencyRate, UpbitPrice, BithumbPrice, BybitPrice },
  data () {
    return {
      loop: true
    }
  },
  created () {
    this.fetchData()
  },
  destroyed () {
    this.loop = false
  },
  methods: {
    fetchData () {
      getPrice().then(response => {
        const currency = Number(response.value.replace(/[^0-9.-]+/g, ''))
        this.$store.dispatch('prices/setCurrency', { exchangeRate: currency })
        const prices = response.data
        for (const item in prices) {
          switch (item) {
            case 'Currency':
              // eslint-disable-next-line no-case-declarations
              const currency = prices[item].Price.find(x => x.Symbol === 'USD_KRW' || x.Symbol === 'USDKRW' || x.Symbol === 'KRWUSD')
              if (currency) {
                currency.exchangeRate = currency.Price
                this.$store.dispatch('prices/setCurrency', currency)
              }
              break
            case 'BybitPrice':
              if (prices[item].Price.length) {
                this.$store.dispatch('prices/setPrice', {
                  key: 'bybitPrice',
                  value: {
                    BTC: prices[item].Price.find(x => x.Symbol === 'BTC') || 0,
                    ETH: prices[item].Price.find(x => x.Symbol === 'ETH') || 0,
                    EOS: prices[item].Price.find(x => x.Symbol === 'EOS') || 0,
                    XRP: prices[item].Price.find(x => x.Symbol === 'XRP') || 0
                  }
                })
              }
              break
            case 'Rules':
              this.$store.dispatch('alarmRule/setRules', prices[item])
              break
          }
        }
        this.loop && setTimeout(this.fetchData, 60000)
      }).catch(err => {
        console.error(err)
        this.loop && setTimeout(this.fetchData, 60000)
      })
    },
    toTimeStrFull (val) {
      if (!val) return
      const iso = new Date((Number(val) + (9 * 60 * 60)) * 1000).toISOString()
      return iso.slice(0, 19).replace('T', ' ')
    },
    timeStrToSec (timestr) {
      const sec = timestr.split(':').reduce((acc, time) => (60 * acc) + +time)
      return sec
    }
  }
}
</script>

<style lang="scss">
  .premium-container {
    .el-table {
      .warning-row .dot {
        background-color: red;
      }
      th,td {
        padding: 2px 0;
      }
      .cell {
        padding-left: 1px;
        padding-right: 1px;
        font-size: 15px;
      }
    }
    .box-card {
      .el-card__body {
        padding: 10px;
        min-width: 10px;
      }
    }
    .el-button {
      padding: 5px 6px
    }
    .dot {
      height: 5px;
      width: 5px;
      background-color: green;
      border-radius: 50%;
      display: inline-block;
    }
  }
</style>

<style lang="scss" scoped>
  .premium-container {
    padding: 6px 6px;
  }
</style>
