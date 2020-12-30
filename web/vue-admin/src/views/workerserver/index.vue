<template>
  <div class="app-container">
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
      <el-table-column label="地址">
        <template slot-scope="scope">
          {{ scope.row.addr }}
        </template>
      </el-table-column>
      <el-table-column label="协议">
        <template slot-scope="scope">
          {{ scope.row.protocol }}
        </template>
      </el-table-column>
      <el-table-column label="扩展">
        <template slot-scope="scope">
          {{ scope.row.extra }}
        </template>
      </el-table-column>
      <el-table-column label="描述">
        <template slot-scope="scope">
          <span>{{ scope.row.description }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态">
        <template slot-scope="scope">{{ scope.row.invalid===0 ? '有效':'无效' }}</template>
      </el-table-column>
      <el-table-column label="添加时间">
        <template slot-scope="scope">{{ scope.row.createdAt | dateFilter }}</template>
      </el-table-column>
      <el-table-column label="操作">
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

    <pagination v-show="total>0" :total="total" :page.sync="search.page" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible" width="500px" top="100px">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="120px" style="width: 400px; margin-left:50px;">
        <el-form-item label="地址" prop="addr">
          <el-input v-model="temp.addr" />
        </el-form-item>
        <el-form-item label="协议" prop="protocol">
          <el-select v-model="temp.protocol" placeholder="请选择">
            <el-option
              v-for="item in protocolList"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="扩展字段">
          <el-input v-model="temp.extra" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="temp.description" />
        </el-form-item>
        <el-form-item label="责任人">
          <el-input v-model="temp.owner" />
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
import { getList, create, update, deleteAction } from '@/api/workServer'
import Pagination from '@/components/Pagination'
const moment = require('moment')

export default {
  filters: {
    dateFilter(val) {
      return moment(val).format('YYYY-MM-DD hh:mm:ss')
    }
  },
  components: {
    Pagination
  },
  data() {
    return {
      protocolList: [
        {
          value: 'HTTP',
          label: 'HTTP'
        },
        {
          value: 'FastCGI',
          label: 'FastCGI'
        },
        {
          value: 'CBNSQ',
          label: 'CBNSQ'
        }
      ],
      search: {
        page: 1
      },
      total: 0,
      list: [],
      listLoading: false,
      textMap: {
        update: '修改',
        create: '新增'
      },
      rules: {
        addr: [{ required: true, message: '必要选项', trigger: 'blur' }],
        protocol: [{ required: true, message: '必要选项', trigger: 'blur' }]
      },
      temp: {
        addr: '',
        protocol: '',
        extra: '',
        description: '',
        owner: '',
        invalid: 0 // 0有效1无效
      },
      dialogStatus: '',
      dialogFormVisible: false
    }
  },
  created() {
    this.getList()
  },
  methods: {
    resetTemp() {
      this.temp = {
        addr: '',
        protocol: '',
        extra: '',
        description: '',
        owner: '',
        invalid: 0
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
      getList({ 'page': this.search.page }).then(response => {
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
