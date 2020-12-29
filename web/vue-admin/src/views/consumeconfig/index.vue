<template>
  <div class="app-container">
    <el-input v-model="search.topic" placeholder="队列名" style="width: 200px;" />
    <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-search" @click="handleFilter">
      搜索
    </el-button>
    <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
      新增
    </el-button>

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
          {{ scope.row.id }}
        </template>
      </el-table-column>
      <el-table-column label="队列名">
        <template slot-scope="scope">
          {{ scope.row.topic }}
        </template>
      </el-table-column>
      <el-table-column label="通道名">
        <template slot-scope="scope">
          {{ scope.row.channel }}
        </template>
      </el-table-column>
      <el-table-column label="描述">
        <template slot-scope="scope">
          <span>{{ scope.row.description }}</span>
        </template>
      </el-table-column>
      <el-table-column label="责任人">
        <template slot-scope="scope">{{ scope.row.owner }}</template>
      </el-table-column>
      <el-table-column label="积压报警阈值">
        <template slot-scope="scope">{{ scope.row.monitorThreshold }}</template>
      </el-table-column>
      <el-table-column label="并发量">
        <template slot-scope="scope">{{ scope.row.handleNum }}</template>
      </el-table-column>
      <el-table-column label="maxInFlight">
        <template slot-scope="scope">{{ scope.row.maxInFlight }}</template>
      </el-table-column>
      <el-table-column label="是否重新入队" width="120px">
        <template slot-scope="scope">
          {{ scope.row.isRequeue | booleanFilter }}
        </template></el-table-column>
      <el-table-column label="超时时间">
        <template slot-scope="scope">{{ scope.row.timeoutDial }}</template>
      </el-table-column>
      <el-table-column label="读超时时间" width="120px">
        <template slot-scope="scope">{{ scope.row.timeoutRead }}</template>
      </el-table-column>
      <el-table-column label="写超时时间" width="120px">
        <template slot-scope="scope">{{ scope.row.timeoutWrite }}</template>
      </el-table-column>
      <el-table-column label="是否有效">
        <template slot-scope="scope">
          {{ scope.row.invalid === 0 ?'有效':'无效' }}
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180px">
        <template slot-scope="scope">{{ scope.row.createdAt | dateFilter }}</template>
      </el-table-column>
      <el-table-column label="操作" width="250">
        <template slot-scope="{row}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            修改
          </el-button>
          <el-button type="primary" size="mini" @click="toWorkPage(row.id)">
            work机
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

    <pagination v-show="total>0" :total="total" :page.sync="search.page" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible" width="500px" top="100px">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="120px" style="width: 400px; margin-left:50px;">
        <el-form-item label="Topic" prop="topic">
          <el-input v-model="temp.topic" />
        </el-form-item>
        <el-form-item label="Channel" prop="channel">
          <el-input v-model="temp.channel" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="temp.description" />
        </el-form-item>
        <el-form-item label="责任人">
          <el-input v-model="temp.owner" />
        </el-form-item>
        <el-form-item label="积压报警阈值">
          <el-input v-model="temp.monitorThreshold" />
        </el-form-item>
        <el-form-item label="并发量">
          <el-input v-model="temp.handleNum" />
        </el-form-item>
        <el-form-item label="maxInFlight">
          <el-input v-model="temp.maxInFlight" />
        </el-form-item>
        <el-form-item label="是否重新入队">
          <el-switch v-model="temp.isRequeue" />
        </el-form-item>
        <el-form-item label="超时时间(秒)">
          <el-input v-model="temp.timeoutDial" />
        </el-form-item>
        <el-form-item label="读超时时间(秒)">
          <el-input v-model="temp.timeoutRead" />
        </el-form-item>
        <el-form-item label="写超时时间(秒)">
          <el-input v-model="temp.timeoutWrite" />
        </el-form-item>
        <el-form-item label="是否有效">
          <el-radio-group v-model="temp.invalid">
            <el-radio :label="0">有效</el-radio>
            <el-radio :label="1">无效</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
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
import { getList, create, update, deleteAction } from '@/api/consume'
import Pagination from '@/components/Pagination'
const moment = require('moment')

export default {
  filters: {
    booleanFilter(val) {
      return val ? '是' : '否'
    },
    dateFilter(val) {
      return moment(val).format('YYYY-MM-DD hh:mm:ss')
    }
  },
  components: {
    Pagination
  },
  data() {
    return {
      list: [],
      total: 0,
      listLoading: false,
      search: {
        topic: '',
        page: 1
      },
      textMap: {
        update: '修改',
        create: '新增'
      },
      rules: {
        topic: [{ required: true, message: '必要选项', trigger: 'blur' }],
        channel: [{ required: true, message: '必要选项', trigger: 'blur' }]
      },
      temp: {
        topic: '',
        channel: '',
        description: '',
        owner: '',
        monitorThreshold: 50000,
        handleNum: 2,
        maxInFlight: 2,
        isRequeue: false,
        timeoutDial: 3590,
        timeoutRead: 3590,
        timeoutWrite: 3590,
        invalid: 0 // 0有效 1无效
      },
      dialogStatus: '',
      dialogFormVisible: false
    }
  },
  created() {
    this.getList()
    console.log(this.$route)
  },
  methods: {
    toWorkPage(id) {
      this.$router.push({
        path: `/consumeServerMap/${id}`
      })
    },
    resetTemp() {
      this.temp = {
        topic: '',
        channel: '',
        description: '',
        owner: '',
        monitorThreshold: 50000,
        handleNum: 2,
        maxInFlight: 2,
        isRequeue: false,
        timeoutDial: 3590,
        timeoutRead: 3590,
        timeoutWrite: 3590,
        invalid: 0 // 0有效 1无效
      }
    },
    handleFilter() {
      this.getList()
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row)

      this.$delete(this.temp, 'createdAt')
      this.$delete(this.temp, 'updatedAt')

      console.log(JSON.stringify(this.temp))

      this.dialogStatus = 'update'
      this.dialogFormVisible = true
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
              title: 'Success',
              message: '更新成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    getList() {
      this.listLoading = true
      getList(this.search).then(response => {
        this.listLoading = false
        this.list = response.result.result
        this.search.page = response.result.page
        this.total = response.result.total
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          console.log(JSON.stringify(this.temp))
          create(this.temp).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '新增成功',
              message: '新增成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    deleteAction(id) {
      deleteAction({ id }).then(() => {
        this.dialogFormVisible = false
        this.$notify({
          title: '删除成功',
          message: '删除成功',
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
