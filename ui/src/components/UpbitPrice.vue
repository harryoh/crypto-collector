<template>
  <el-row style="margin-top:5px;">
    <el-card class="box-card">
      <div class="clearfix">
        <span style="padding-right:3px;"><strong>업비트</strong></span>
        <img v-if="!isLive" :src="warningImg" width="12">
      </div>
      <div>
        <el-table
          :data="upbitData"
          stripe
          style="width: 100%"
          size="medium"
          :row-class-name="tableRowClassName"
        >
          <el-table-column
            prop="Symbol"
            label="코인"
            min-width="30px"
          />
          <el-table-column
            prop="Price"
            label="가격"
            min-width="60px"
            :formatter="numberWithCommas"
          />
          <el-table-column
            prop="FixPremium"
            label="고정P"
            min-width="60px"
          />
          <el-table-column
            prop="CurrPremium"
            label="현재P"
            min-width="60px"
          />
          <el-table-column width="10px">
            <span class="dot" />
          </el-table-column>
          <el-table-column
            prop="Timestamp"
            label="시간"
            min-width="30px"
            :formatter="toSecAgo"
          />
        </el-table>
      </div>
    </el-card>
  </el-row>
</template>

<script>
import { mapGetters } from 'vuex'
import warningImg from '@/assets/warning.png'
import {
  tableRowClassName, numberWithCommas, getPremium, toSecAgo
} from '@/utils'

export default {
  data () {
    return {
      warningImg: warningImg,
      upbitData: [],
      ws: null,
      socketTimeout: null,
      isLive: false
    }
  },
  computed: {
    ...mapGetters([
      'currencyRate',
      'bybitPrice',
      'upbitPrice'
    ])
  },
  created () {
    this.start()
  },
  destroyed () {
    this.close()
  },
  methods: {
    tableRowClassName,
    numberWithCommas,
    getPremium,
    toSecAgo,
    heartbeat () {
      clearTimeout(this.socketTimeout)
      this.socketTimeout = setTimeout(() => {
        // console.log('Upbit Heartbeat Error')
        this.isLive = false
        this.start()
      }, 30000)
    },
    close () {
      clearTimeout(this.socketTimeout)
      this.ws && this.ws.close()
    },
    start () {
      this.close()
      this.getPrice()
    },
    getPrice () {
      const wsurl = 'wss://api.upbit.com/websocket/v1'
      this.ws = new WebSocket(wsurl)

      this.ws.binaryType = 'arraybuffer'
      this.ws.onopen = () => {
        this.ws.send('PING')
        this.ws.send(
          '[{"ticket":"UNIQUE_TICKET"},{"type":"trade","codes":["KRW-BTC","KRW-ETH","KRW-EOS","KRW-XRP"]}]'
        )
        this.isLive = true
      }

      this.ws.onmessage = this.setPrice

      this.ws.onerror = (err) => {
        console.error('Upbit Connection Error')
        console.error(err)
        this.isLive = false
        setTimeout(this.start, 5000)
      }
    },
    setPrice (evt) {
      const enc = new TextDecoder('utf-8')
      const arr = new Uint8Array(evt.data)
      // eslint-disable-next-line
      const info = JSON.parse(enc.decode(arr).replace(/|/g, ''))
      // eslint-disable-next-line no-prototype-builtins
      if (info.hasOwnProperty('error')) {
        console.error(info.error)
        return
      }
      // eslint-disable-next-line no-prototype-builtins
      if (info.hasOwnProperty('status')) {
        if (info.status === 'UP') {
          clearTimeout(this.socketTimeout)
          this.heartbeat()
        }
        return
      }

      const symbol = info.code.substring(4, 8)

      this.$store.dispatch('prices/setCoin', {
        key: 'upbitPrice',
        coin: symbol,
        value: {
          Symbol: symbol,
          Price: info.trade_price,
          Timestamp: parseInt(info.trade_timestamp / 1000)
        }
      })
      if (Object.keys(this.bybitPrice).length < 3) return
      this.upbitData = Object.values(this.upbitPrice).map(o => {
        const bybit = this.bybitPrice[o.Symbol]?.Price
        return {
          ...o,
          FixPremium: this.getPremium(bybit, o.Price, this.currencyRate.fixExchangeRate) + '%',
          CurrPremium: this.getPremium(bybit, o.Price, this.currencyRate.exchangeRate) + '%'
        }
      })
    }
  }
}
</script>
