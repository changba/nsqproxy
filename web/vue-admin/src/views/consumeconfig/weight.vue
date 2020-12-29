<template>
  <div class="app-container">
    <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
      新增
    </el-button>

    <div style="margin-top:20px;">
      <el-divider />
      <span>{{ fatherInfo.topic }} / {{ fatherInfo.channel }}</span>
      <el-divider />
    </div>

    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
      class="common-table"
    >
      <el-table-column label="ID">
        <template slot-scope="scope">
          {{ scope.row.serverid }}
        </template>
      </el-table-column>
      <el-table-column label="地址">
        <template slot-scope="scope">
          {{ scope.row.workServer.addr }}
        </template>
      </el-table-column>
      <el-table-column label="描述">
        <template slot-scope="scope">
          {{ scope.row.workServer.description }}
        </template>
      </el-table-column>
      <el-table-column label="权重">
        <template slot-scope="scope">
          {{ scope.row.weight }}
        </template>
      </el-table-column>
      <el-table-column label="关联状态">
        <template slot-scope="scope">
          {{ scope.row.invalid | statusFilter }}
        </template>
      </el-table-column>
      <el-table-column label="服务器状态">
        <template slot-scope="scope">
          {{ scope.row.workServer.invalid | statusFilter }}
        </template>
      </el-table-column>
      <el-table-column label="添加时间" width="200">
        <template slot-scope="scope">
          {{ scope.row.createdAt | format }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template slot-scope="{row}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            修改
          </el-button>
          <el-popconfirm
            confirm-button-text="好的"
            cancel-button-text="不用了"
            icon="el-icon-info"
            icon-color="red"
            title="确定删除吗？"
            @onConfirm="deleteAction(row.id)"
          >
            <el-button slot="reference" type="danger" size="mini" style="margin-left:10px;">
              删除
            </el-button>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :before-close="handleClose" :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible" width="500px" top="100px">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="120px" style="width: 400px; margin-left:50px;">
        <el-form-item label="选择机器" prop="serverid">
          <el-select v-model="temp.serverid" :multiple="dialogStatus==='create'" placeholder="请选择">
            <el-option
              v-for="item in machine"
              :key="item.id"
              :label="item.addr"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="权重(正整数)" prop="weight">
          <el-input v-model="temp.weight" />
        </el-form-item>
        <el-form-item label="是否有效">
          <el-radio-group v-model="temp.invalid">
            <el-radio :label="0">有效</el-radio>
            <el-radio :label="1">无效</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="handleCancel">
          取消
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          确认
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getList, create, update, deleteAction } from '@/api/consumeWorkServerMap'
import { getAll } from '@/api/workServer'
const moment = require('moment')

export default {
  filters: {
    statusFilter(val) {
      return val === 0 ? '有效' : '无效'
    },
    format(val) {
      return moment(val).format('YYYY-MM-DD hh:mm:ss')
    }
  },
  data() {
    const checkWeight = (rule, value, callback) => {
      if (value === '') {
        return callback(new Error('权重不能为空'))
      } else {
        if (!isNaN(value) && value < 1) {
          callback(new Error('请输入大于0的数字值'))
        } else {
          callback()
        }
      }
    }

    return {
      machine: [],
      list: [],
      total: 0,
      listLoading: false,
      textMap: {
        update: '修改',
        create: '新增'
      },
      rules: {
        weight: [{ validator: checkWeight, trigger: 'blur' }],
        serverid: [{ required: true, message: '必要选项', trigger: 'blur' }]
      },
      temp: {
        serverid: '',
        weight: 1,
        invalid: 0
      },
      dialogStatus: '',
      dialogFormVisible: false,
      consumeid: '',
      fatherInfo: {}
    }
  },
  created() {
    if (this.$route.params.id) {
      this.consumeid = this.$route.params.id
      this.getList()
      this.getMachine()
    }
  },
  methods: {
    handleCancel() {
      this.dialogFormVisible = false
      this.resetTemp()
    },
    handleClose(done) {
      this.resetTemp()
      done()
    },
    resetTemp() {
      this.temp = {
        serverid: '',
        weight: 1,
        invalid: 0
      }
    },
    handleFilter() {
      this.getList()
    },
    handleCreate() {
      this.resetTemp()
      console.log(JSON.stringify(this.temp))
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true

      this.$delete(this.temp, 'createdAt')
      this.$delete(this.temp, 'updatedAt')
      this.$delete(this.temp, 'workServer')
      // console.log(JSON.stringify(this.temp))

      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)

          update(tempData).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '更新成功',
              // message: '更新成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
            this.resetTemp()
          })
        }
      })
    },
    getList() {
      this.listLoading = true
      getList({ id: this.consumeid }).then(response => {
        this.listLoading = false
        this.list = response.result.serverList
        this.fatherInfo = {
          'topic': response.result.topic,
          'channel': response.result.channel,
          'description': response.result.description
        }
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.temp.serverid = this.temp.serverid.join(',')
          this.temp = Object.assign(this.temp, { 'consumeid': this.consumeid })
          create(this.temp).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '新增成功',
              // message: '新增成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
          }).catch(error => {
            this.dialogFormVisible = false
            console.log(error)
          })
        }
      })
    },
    getMachine() {
      getAll().then((response) => {
        this.machine = response.result
      })
    },
    deleteAction(id) {
      deleteAction({ id }).then(() => {
        this.dialogFormVisible = false
        this.$notify({
          title: '删除成功',
          // message: '删除成功',
          type: 'success',
          duration: 2000
        })
        this.getList()
      })
    }
  }
}
</script>
<style lang="scss" scoped>
.common-table{
  margin-top: 10px;
}
</style>
