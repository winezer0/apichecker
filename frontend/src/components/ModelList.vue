<template>
  <el-dialog :model-value="visible" @update:model-value="$emit('update:visible', $event)" :title="t('modelList')" width="900px" destroy-on-close>
    <div class="model-list-header">
      <div>
        <el-button type="primary" @click="batchTestModels" :loading="batchTesting">{{ t('batchTest') }}</el-button>
        <el-button @click="exportModelList" style="margin-left: 10px;">{{ t('exportList') }}</el-button>
        <el-button @click="cacheModelResults" style="margin-left: 10px;">{{ t('cacheResults') }}</el-button>
      </div>
      <div class="test-prompt-input">
        <span>{{ t('testPrompt') }}：</span>
        <el-input
          v-model="testPrompt"
          type="textarea"
          :rows="1"
          readonly
          style="width: 300px"
        ></el-input>
      </div>
    </div>
    
    <el-table :data="tableData" border v-loading="loading">
      <el-table-column prop="name" :label="t('modelId')" width="220" />
      <el-table-column prop="result" :label="t('testResult')">
        <template #default="scope">
          <span v-if="scope.row.errorMsg" class="error-msg">{{ scope.row.errorMsg }}</span>
          <span v-else>{{ scope.row.result }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('operation')" width="120" align="center">
        <template #default="scope">
          <el-button link type="primary" @click="testSingleModel(scope.row)" :loading="scope.row.testing">{{ t('test') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { TestModel, GetCachedModelResult, GetModelCache, SaveModelCache } from '../../wailsjs/go/main/App'
import { translations, LangType } from '../locales'

// 从父组件获取当前语言
const currentLang = computed(() => props.currentLang)

const t = (key: string) => {
  const keys = key.split('.')
  let value: any = translations[currentLang.value]
  for (const k of keys) {
    value = value?.[k]
  }
  return value || key
}

const props = defineProps({
  visible: Boolean,
  models: {
    type: Array as () => string[],
    default: () => []
  },
  protocol: { type: String, default: '' },
  apiKey: { type: String, default: '' },
  baseUrl: { type: String, default: '' },
  prompt: { type: String, default: '' },
  currentLang: {
    type: String as () => 'zh' | 'en',
    default: 'zh'
  }
})

const emit = defineEmits(['update:visible', 'update:prompt'])

interface ModelRow {
  name: string
  result: string
  errorMsg: string
  testing: boolean
}

const tableData = ref<ModelRow[]>([])
const testPrompt = ref('')
const batchTesting = ref(false)
const loading = ref(false)
const modelCache = ref<{ apiModels: Record<string, Record<string, string>> }>({ apiModels: {} })

// 监听 prompt prop 的变化
watch(() => props.prompt, (newVal) => {
  testPrompt.value = newVal || ''
})

watch(() => props.visible, async (newVal) => {
  if (newVal) {
    loading.value = true
    // 直接使用 props.prompt，不再使用默认值
    testPrompt.value = props.prompt || ''
    tableData.value = props.models.map(m => ({
      name: m,
      result: '',
      errorMsg: '',
      testing: false
    }))
    
    // 加载模型缓存
    try {
      const cache = await GetModelCache()
      modelCache.value = cache
    } catch (err) {
      console.error("Failed to load model cache:", err)
    }
    
    // Load cached results for each model
    for (let row of tableData.value) {
      try {
        const apiKey = `${props.protocol}@${props.apiKey.substring(0, 8)}` // 使用 apiKey 的前 8 个字符作为标识
        if (modelCache.value.apiModels && modelCache.value.apiModels[apiKey] && modelCache.value.apiModels[apiKey][row.name]) {
          row.result = modelCache.value.apiModels[apiKey][row.name]
        }
      } catch (err) {
        console.error("Failed to get cache:", err)
      }
    }
    loading.value = false
  }
})

const testSingleModel = async (row: ModelRow) => {
  row.testing = true
  try {
    const res = await TestModel(props.protocol, props.apiKey, props.baseUrl, row.name, testPrompt.value)
    row.result = res
    row.errorMsg = '' // 测试成功时清空错误消息
    
    // 保存模型缓存
    await saveModelCache(row.name, res)
  } catch (err: any) {
    row.errorMsg = typeof err === 'string' ? err : (err.message || err.toString() || '未知错误')
    // 保存错误消息到缓存
    await saveModelCache(row.name, row.errorMsg)
  } finally {
    row.testing = false
  }
}

const batchTestModels = async () => {
  batchTesting.value = true
  for (let row of tableData.value) {
    // skip if already tested and we don't want to retest? Actually, usually batch test forces retest
    // but to be safe, let's just run them all
    await testSingleModel(row)
  }
  batchTesting.value = false
  ElMessage.success(t('messages.batchTestComplete'))
}

const exportModelList = () => {
  if (tableData.value.length === 0) return
  let csv = 'Model ID,Test Result,Error Message\n'
  tableData.value.forEach(row => {
    csv += `${row.name},"${row.result.replace(/"/g, '""')}","${row.errorMsg.replace(/"/g, '""')}"\n`
  })
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `model_test_results_${props.protocol}_${props.apiKey.substring(0, 8)}.csv`
  link.click()
}

// 保存模型缓存
const saveModelCache = async (modelName: string, result: string) => {
  try {
    const apiKey = `${props.protocol}@${props.apiKey.substring(0, 8)}` // 使用 apiKey 的前 8 个字符作为标识
    if (!modelCache.value.apiModels) {
      modelCache.value.apiModels = {}
    }
    if (!modelCache.value.apiModels[apiKey]) {
      modelCache.value.apiModels[apiKey] = {}
    }
    modelCache.value.apiModels[apiKey][modelName] = result
    await SaveModelCache(modelCache.value)
  } catch (err) {
    console.error("Failed to save model cache:", err)
  }
}

// 缓存当前模型结果
const cacheModelResults = async () => {
  try {
    const apiKey = `${props.protocol}@${props.apiKey.substring(0, 8)}` // 使用 apiKey 的前 8 个字符作为标识
    if (!modelCache.value.apiModels) {
      modelCache.value.apiModels = {}
    }
    if (!modelCache.value.apiModels[apiKey]) {
      modelCache.value.apiModels[apiKey] = {}
    }
    
    // 缓存所有模型结果
    for (const row of tableData.value) {
      modelCache.value.apiModels[apiKey][row.name] = row.errorMsg || row.result
    }
    
    await SaveModelCache(modelCache.value)
    ElMessage.success('模型结果已缓存')
  } catch (err) {
    console.error("Failed to cache model results:", err)
    ElMessage.error('缓存模型结果失败')
  }
}
</script>

<style scoped>
.model-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.test-prompt-input {
  display: flex;
  align-items: center;
}

.error-msg {
  color: #f56c6c;
}
</style>
