<template>
  <el-row>
    <el-card class="box-card">
      <div>
        <el-table
          :data="[currencyRate]"
          size="medium"
          style="width: 100%"
        >
          <el-table-column
            prop="Symbol"
            label=""
            min-width="60px"
          />
          <el-table-column
            prop="fixExchangeRate"
            label="고정환율"
            min-width="50px"
          />
          <el-table-column
            prop="exchangeRate"
            label="현재환율"
            min-width="50px"
            :formatter="simpleFloat"
          />
          <el-table-column
            prop="Timestamp"
            label="시간"
            min-width="60px"
            :formatter="toTimeStr"
          />
        </el-table>
      </div>
    </el-card>
  </el-row>
</template>

<script>

import { mapGetters } from 'vuex'
import { toTimeStr, simpleFloat } from '@/utils'

export default {
  computed: {
    ...mapGetters([
      'currencyRate'
    ])
  },
  methods: {
    toTimeStr,
    simpleFloat
  }
}
</script>
