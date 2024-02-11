<template>
  <el-row style="margin-top:5px;">
    <el-card class="box-card">
      <div class="clearfix">
        <span style="padding-right:3px;"><strong>바이빗</strong></span>
        <img v-if="!isLive" :src="warningImg" width="12">
      </div>
      <div>
        <el-table
          :data="bybitData"
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
            prop="FundingRate"
            label="펀딩비용"
            min-width="60px"
          />
          <el-table-column
            prop="PredictedFundingRate"
            label="N펀딩비용"
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
  tableRowClassName, numberWithCommas, toSecAgo
} from '@/utils'

export default {
  data () {
    return {
      warningImg: warningImg,
      bybitData: [],
      ws: null,
      socketTimeout: null,
      isLive: false
    }
  },
  computed: {
    ...mapGetters([
      'bybitPrice'
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
    toSecAgo,
    heartbeat () {
      clearTimeout(this.socketTimeout)
      this.ws.send('{"op":"ping"}')
      this.socketTimeout = setTimeout(() => {
        // console.log('Bybit Heartbeat Error')
        this.isLive = false
        this.start()
      }, 10000)
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
      const wsurl = 'wss://stream.bybit.com/v5/public/inverse'
      this.ws = new WebSocket(wsurl)

      this.ws.onopen = () => {
        this.ws.send('{"op": "subscribe", "args": ["publicTrade.BTCUSD", "publicTrade.ETHUSD", "publicTrade.EOSUSD", "publicTrade.XRPUSD"]}')
        this.heartbeat()
        this.isLive = true
      }

      this.ws.onmessage = this.setPrice

      this.ws.onerror = (err) => {
        console.error('Bybit Connection Error')
        console.error(err)
        this.isLive = false
        setTimeout(this.start, 5000)
      }
    },
    setPrice (evt) {
      const res = JSON.parse(evt.data)
      // eslint-disable-next-line no-prototype-builtins
      if (res.hasOwnProperty('ret_msg')) {
        if (res.success === true && res.ret_msg === 'pong') {
          clearTimeout(this.socketTimeout)
          setTimeout(this.heartbeat, 30000)
        }
        return
      }

      if (!res.data) return
      const info = res.data[res.data.length - 1]
      const symbol = info.s.substring(0, 3)

      this.$store.dispatch('prices/setCoin', {
        key: 'bybitPrice',
        coin: symbol,
        value: {
          ...this.bybitPrice[symbol],
          Symbol: symbol,
          Price: info.p,
          Timestamp: parseInt(info.T / 1000)
        }
      })
      this.bybitData = Object.values(this.bybitPrice)
    }
  }
}
</script>
