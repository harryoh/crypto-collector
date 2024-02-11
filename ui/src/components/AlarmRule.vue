<template>
  <el-row style="margin-top:5px;">
    <el-card class="box-card">
      <div class="clearfix">
        <strong>알람</strong>
      </div>
      <div>
        <el-table
          :data="alarmRules"
          size="medium"
          style="width: 100%"
        >
          <el-table-column
            prop="Exchange"
            label="거래소"
            min-width="60px"
          />
          <el-table-column
            prop="Symbol"
            label="코인"
            min-width="60px"
          />
          <el-table-column
            prop="Use"
            label="사용여부"
            min-width="60px"
            :formatter="boolToStr"
          />
          <el-table-column
            prop="AlarmMin"
            label="최소P"
            min-width="60px"
          />
          <el-table-column
            prop="AlarmMax"
            label="최대P"
            min-width="60px"
          />
          <el-table-column
            label=""
            min-width="40px"
          >
            <template v-slot:default="scope">
              <el-button
                type="text"
                size="medium"
                @click="alarmEdit(scope.$index, scope.row)"
              >수정</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <el-dialog title="Alarm Edit" :visible.sync="alarmFormVisible" width="60%">
      <el-form :model="alarmForm">
        <el-form-item label="사용">
          <el-radio-group v-model="alarmForm.Use">
            <el-radio :label="true">Y</el-radio>
            <el-radio :label="false">N</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-select v-model="alarmForm.Exchange" placeholder="거래소">
            <el-option
              v-for="item in exchangeList"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="alarmForm.Symbol" placeholder="코인">
            <el-option
              v-for="item in coinList"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="최소P">
          <el-input-number
            v-model="alarmForm.AlarmMin"
            size="mini"
            :precision="1"
            :step="0.1"
          />
        </el-form-item>
        <el-form-item label="최대P">
          <el-input-number
            v-model="alarmForm.AlarmMax"
            size="mini"
            :precision="1"
            :step="0.1"
          />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="alarmFormVisible = false">취소</el-button>
        <el-button type="primary" @click="onAlarmSubmit">저장</el-button>
      </span>
    </el-dialog>
  </el-row>
</template>

<script>

import { mapGetters } from 'vuex'
import { Message } from 'element-ui'
import { updateAlarm } from '@/api/price'

export default {
  data () {
    return {
      alarmFormVisible: false,
      alarmForm: {},
      exchangeList: [{
        value: 'all',
        label: 'All'
      }, {
        value: 'bithumb',
        label: 'Bithumb'
      }, {
        value: 'upbit',
        label: 'Upbit'
      }],
      coinList: [{
        value: 'all',
        label: 'All'
      }, {
        value: 'BTC',
        label: 'BTC'
      }, {
        value: 'ETH',
        label: 'ETH'
      }, {
        value: 'XRP',
        label: 'XRP'
      }, {
        value: 'EOS',
        label: 'EOS'
      }]
    }
  },
  computed: {
    ...mapGetters([
      'alarmRules'
    ])
  },
  methods: {
    boolToStr (row, col, val) {
      return (val) ? 'Y' : 'N'
    },
    alarmEdit (index, row) {
      this.alarmForm = this.alarmRules[index]
      this.alarmForm.index = index
      this.alarmFormVisible = true
    },
    onAlarmSubmit () {
      if (this.alarmForm.AlarmMin > this.alarmForm.AlarmMax) {
        Message({
          message: '최소값이 최대값보다 더 큽니다.',
          type: 'error',
          duration: 5 * 1000
        })
        return false
      }
      this.alarmRules[this.alarmForm.index] = this.alarmForm
      updateAlarm(this.alarmRules).then(response => {
        this.alarmFormVisible = false
      })
    }
  }
}
</script>
