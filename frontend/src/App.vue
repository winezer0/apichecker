<template>
  <div class="app-container">
    <!-- 顶部工具栏与操作 -->
    <div class="toolbar">
      <div class="toolbar-actions" style="margin-bottom: 15px; display: flex; justify-content: space-between; align-items: center;">
        <div>
          <el-button type="primary" @click="saveConfig">{{ t('saveConfig') }}</el-button>
          <el-button type="success" @click="batchTest" :loading="testing">{{ t('batchTest') }}</el-button>
          <el-button @click="exportList">{{ t('exportList') }}</el-button>
          <el-button type="danger" @click="clearList">{{ t('clearList') }}</el-button>
          <el-button @click="showPromptSettings">{{ t('promptSettings') }}</el-button>
        </div>
        <div>
          <el-button @click="toggleLanguage">{{ t('switchLang') }}</el-button>
          <el-button @click="showStartupPassSettings">{{ t('startupPass') }}</el-button>
        </div>
      </div>
      <el-form :inline="true" class="toolbar-form" @submit.prevent>
        <el-form-item :label="t('protocol')">
          <el-select v-model="addForm.protocol" :placeholder="t('selectProtocol')" style="width: 120px;">
            <el-option v-for="p in protocols" :key="p" :label="p" :value="p"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('alias')">
          <el-input v-model="addForm.alias" :placeholder="t('placeholder.alias')" style="width: 150px;"></el-input>
        </el-form-item>
        <el-form-item :label="t('apiKey')">
          <el-input v-model="addForm.apiKey" :placeholder="t('placeholder.apiKey')" style="width: 200px;"></el-input>
        </el-form-item>
        <el-form-item :label="t('baseUrl')">
          <el-select v-model="addForm.baseUrl" :placeholder="t('placeholder.selectBaseUrl')" style="width: 250px;" allow-create filterable default-first-option popper-class="baseurl-select">
            <el-option
              v-for="url in commonBaseUrls"
              :key="url"
              :label="url"
              :value="url">
              <span style="text-align: left; display: inline-block; width: 100%;">{{ url }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="addKey">{{ t('add') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 数据表格 -->
    <div class="table-container">
      <el-table 
        :data="keysList" 
        style="width: 100%" 
        border 
        highlight-current-row 
        @current-change="handleRowChange"
        @cell-dblclick="handleCellDblClick">
        <el-table-column prop="protocol" :label="t('protocol')" width="100"></el-table-column>
        <el-table-column prop="alias" :label="t('alias')" width="150"></el-table-column>
        <el-table-column prop="apiKey" :label="t('apiKey') + ' (' + t('copy') + ')'" width="220">
          <template #default="scope">
            <span>{{ maskKey(scope.row.apiKey) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="baseUrl" :label="t('baseUrl') + ' (' + t('copy') + ')'" min-width="200">
          <template #default="scope">
            <span>{{ scope.row.baseUrl }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('status')" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.success)">
              {{ getStatusText(scope.row.success, scope.row.errorMsg) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('modelCount')" width="120">
          <template #default="scope">
            <el-button v-if="scope.row.models && scope.row.models.length > 0" link type="primary" @click="showModels(scope.row)">
              {{ scope.row.models.length }} {{ isZh ? t('countZh') : '' }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('operation')" width="120">
          <template #default="scope">
            <el-button link type="primary" @click="testSingle(scope.row)" :loading="scope.row.testing">{{ t('test') }}</el-button>
            <el-button link type="danger" @click="deleteRow(scope.$index)">{{ isZh ? t('delete') : t('operationDelete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 模型列表独立组件 -->
    <ModelList
      v-model:visible="modelsModalVisible"
      v-model:prompt="currentPrompt"
      :models="currentModels"
      :protocol="currentProtocol"
      :apiKey="currentApiKey"
      :baseUrl="currentBaseUrl"
      :currentLang="currentLang"
    />

    <!-- 底部聊天调试框 -->
    <div class="chat-container">
      <div class="chat-header">
        <el-select v-model="selectedApiIndex" :placeholder="t('selectApi')" style="width: 200px; margin-right: 10px;" size="small" @change="onApiSelected">
          <el-option 
            v-for="(item, index) in validKeysList" 
            :key="index" 
            :label="item.alias ? `${item.alias} (${maskKey(item.apiKey)})` : maskKey(item.apiKey)" 
            :value="index">
          </el-option>
        </el-select>
        
        <el-select v-model="selectedModel" :placeholder="t('selectDialogModel')" style="width: 200px;" size="small" :disabled="!selectedApiIndex && selectedApiIndex !== 0">
          <el-option 
            v-for="(model, idx) in availableModelsForChat" 
            :key="idx" 
            :label="model" 
            :value="model">
          </el-option>
        </el-select>
      </div>
      <div class="chat-history" ref="chatHistoryRef">
        <div v-for="(msg, index) in chatHistory" :key="index" :class="['chat-msg', msg.role]">
          <div class="msg-role">{{ msg.role === 'user' ? t('you') : 'AI' }}</div>
          <div class="msg-content markdown-body" v-html="renderMarkdown(msg.content)"></div>
        </div>
      </div>
      <div class="chat-input">
        <el-input
          v-model="chatInput"
          class="chat-textarea"
          type="textarea"
          :rows="2"
          :placeholder="t('placeholder.chatInput')"
          @keydown.enter.prevent.exact="sendChat"
          @keydown.enter.shift.exact="chatInput += '\n'"
        ></el-input>
        <el-button type="primary" @click="sendChat" :loading="chatting">{{ t('send') }}</el-button>
      </div>
    </div>

    <!-- 启动密码验证弹框 -->
    <el-dialog
      v-model="startupPassModalVisible"
      :title="t('startupPassTitle')"
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
    >
      <el-form>
        <el-form-item :label="t('startupPass')">
          <el-input
            v-model="startupPass"
            type="password"
            :placeholder="t('placeholder.startupPass')"
            @keydown.enter="verifyStartupPass"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="verifyStartupPass">{{ t('confirm') }}</el-button>
        <el-button @click="cancelStartupPass">{{ t('cancel') }}</el-button>
      </template>
    </el-dialog>

    <!-- 启动密码设置弹框 -->
    <el-dialog
      v-model="startupPassSettingsModalVisible"
      :title="t('startupPassSettings')"
      width="400px"
    >
      <el-form>
        <el-form-item :label="t('newPassword')">
          <el-input
            v-model="newStartupPass"
            type="password"
            :placeholder="t('placeholder.newPassword')"
            @keydown.enter="saveStartupPass"
          ></el-input>
        </el-form-item>
        <el-form-item :label="t('confirmPassword')">
          <el-input
            v-model="newStartupPassConfirm"
            type="password"
            :placeholder="t('placeholder.confirmPassword')"
            @keydown.enter="saveStartupPass"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="clearStartupPass" type="warning">{{ t('clearEncryption') }}</el-button>
        <el-button @click="saveStartupPass" type="primary">{{ t('confirm') }}</el-button>
        <el-button @click="cancelStartupPassSettings">{{ t('cancel') }}</el-button>
      </template>
    </el-dialog>

    <!-- 提问设置弹框 -->
    <el-dialog
      v-model="promptModalVisible"
      :title="t('promptSettings')"
      width="600px"
    >
      <el-form>
        <el-form-item :label="t('prompt')">
          <el-input
            v-model="currentPrompt"
            type="textarea"
            :rows="8"
            :placeholder="t('placeholder.prompt')"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="savePromptSettings" type="primary">{{ t('confirm') }}</el-button>
        <el-button @click="cancelPromptSettings">{{ t('cancel') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, nextTick, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { marked } from 'marked'
import {
  CheckKey,
  BatchCheckKeys,
  Chat,
  GetAppConfig,
  GetDefaultBaseURLs,
  GetProtocols,
  SaveAppConfig,
  ExportEncryptedList,
  DecryptAPIKey,
  VerifyStartupPass,
  LoadConfigWithPassword,
  GetStartupPassStatus,
  GetPrompt,
  SetPrompt,
  SetStartupPass,
  ClearStartupPass,
  GetModelCache,
  SaveModelCache,
  TestModel
} from '../wailsjs/go/main/App'
import { LangType, translations, createTranslator } from './locales'
import { appconfig } from '../wailsjs/go/models'
import ModelList from './components/ModelList.vue'
import { Quit } from '../wailsjs/runtime'

interface KeyRow {
  protocol: string
  alias: string
  apiKey: string
  baseUrl: string
  success: boolean
  models: string[]
  modelCount: number
  errorMsg: string
  testing?: boolean
}

interface AppConfigPayload {
  keysList: KeyRow[]
  chatState: {
    selectedApiIndex: number
    selectedModel: string
  }
  defaults?: Record<string, string>
  startupPass?: string
  prompt?: string
}

const protocols = ref<string[]>([])
const defaultBaseUrls = ref<Record<string, string>>({})
const addForm = reactive({ protocol: '', alias: '', apiKey: '', baseUrl: '' })
const keysList = ref<KeyRow[]>([])
const testing = ref(false)
const isHydrating = ref(true)
const modelsModalVisible = ref(false)
const currentModels = ref<string[]>([])
const startupPassModalVisible = ref(false)
const startupPass = ref('')
const startupPassConfirm = ref('')
const startupPassSettingsModalVisible = ref(false)
const newStartupPass = ref('')
const newStartupPassConfirm = ref('')
const promptModalVisible = ref(false)
const currentPrompt = ref('')
const modelCache = ref<appconfig.ModelCache>({ apiModels: {} })

// 常用 BaseURL 列表
const commonBaseUrls = ref<string[]>([
  'https://api.openai.com/v1',
  'https://api.anthropic.com/v1',
  'https://dashscope.aliyuncs.com/compatible-mode/v1'
])

// 语言配置
const currentLang = ref<LangType>('zh')

// 创建响应式翻译函数
const t = (key: string) => {
  const keys = key.split('.')
  let value: any = translations[currentLang.value]
  for (const k of keys) {
    value = value?.[k]
  }
  return value || key
}

const currentProtocol = ref('')
const currentApiKey = ref('')
const currentBaseUrl = ref('')
const selectedApiIndex = ref<number | null>(null)
const selectedModel = ref<string>('')
const availableModelsForChat = ref<string[]>([])
const chatHistory = ref<{ role: string, content: string }[]>([])
const chatInput = ref('')
const chatting = ref(false)
const chatHistoryRef = ref<HTMLElement | null>(null)

const validKeysList = computed(() => keysList.value.filter(k => k.success))

const getDefaultBaseUrl = (protocol: string) => {
  return defaultBaseUrls.value[protocol] || ''
}

const getEffectiveBaseUrl = (protocol: string, baseUrl: string) => {
  const trimmedBaseUrl = baseUrl.trim()
  return trimmedBaseUrl || getDefaultBaseUrl(protocol)
}

const buildConfigPayload = (): appconfig.AppConfig => {
  return new appconfig.AppConfig({
    keysList: keysList.value.map(row => ({
      ...row,
      baseUrl: getEffectiveBaseUrl(row.protocol, row.baseUrl),
      modelCount: row.models.length,
      models: [...row.models],
      testing: false
    })),
    chatState: {
      selectedApiIndex: selectedApiIndex.value ?? -1,
      selectedModel: selectedModel.value
    },
    defaults: { ...defaultBaseUrls.value },
    prompt: currentPrompt.value
  })
}

const applyConfigPayload = (payload: AppConfigPayload | appconfig.AppConfig) => {
  defaultBaseUrls.value = {
    ...defaultBaseUrls.value,
    ...(payload.defaults || {})
  }

  addForm.protocol = protocols.value[0] || 'openai'
  addForm.alias = ''
  addForm.apiKey = ''
  addForm.baseUrl = getDefaultBaseUrl(addForm.protocol)

  keysList.value = (payload.keysList || []).map(row => ({
    protocol: row.protocol,
    alias: row.alias,
    apiKey: row.apiKey,
    baseUrl: getEffectiveBaseUrl(row.protocol, row.baseUrl || ''),
    success: row.success || false,
    models: row.models || [],
    modelCount: row.modelCount || 0,
    errorMsg: row.errorMsg || '',
    testing: false
  }))

  selectedApiIndex.value = payload.chatState?.selectedApiIndex >= 0 ? payload.chatState.selectedApiIndex : null
  selectedModel.value = payload.chatState?.selectedModel || ''

  // 处理启动密码 - 不自动填充，保持为空
  startupPass.value = ''
  
  // 加载提示词
  currentPrompt.value = payload.prompt || ''
}

const loadInitialState = async () => {
  try {
    isHydrating.value = true
    
    // 加载协议列表和默认 Base URL
    const [protocolRes, defaultRes] = await Promise.all([GetProtocols(), GetDefaultBaseURLs()])
    protocols.value = protocolRes || []
    defaultBaseUrls.value = defaultRes || {}

    addForm.protocol = protocols.value[0] || 'openai'
    addForm.baseUrl = getDefaultBaseUrl(addForm.protocol)

    // 直接显示启动密码输入框，不加载配置
    // 配置将在密码验证成功后加载
    startupPassModalVisible.value = true
  } catch (err) {
    console.error("Failed to initialize app:", err)
    ElMessage.error('初始化失败，请检查后端服务')
  } finally {
    isHydrating.value = false
  }
}

const switchLanguage = () => {
  currentLang.value = currentLang.value === 'zh' ? 'en' : 'zh'
}

const toggleLanguage = () => {
  switchLanguage()
}

const verifyStartupPass = async () => {
  try {
    const isValid = await VerifyStartupPass(startupPass.value)
    if (isValid) {
      // 密码验证成功，加载配置
      await LoadConfigWithPassword(startupPass.value)
      startupPassModalVisible.value = false
      startupPass.value = ''
      startupPassConfirm.value = ''
      await loadConfigAfterVerify()
    } else {
      ElMessage.error('密码错误，请重新输入')
      startupPass.value = ''
    }
  } catch (err: any) {
    console.error("Verification failed:", err)
    ElMessage.error('密码验证失败：' + (err?.toString?.() || 'unknown error'))
  }
}

const loadConfigAfterVerify = async () => {
  try {
    // 加载配置（使用启动密码）
    const savedConfig = await GetAppConfig()
    applyConfigPayload(savedConfig)
    
    // 加载模型缓存
    try {
      const cachedModel = await GetModelCache()
      if (cachedModel) {
        modelCache.value = cachedModel
      }
    } catch (err: any) {
      console.error("Failed to load model cache:", err)
    }
    
    // 加载提示词
    try {
      const prompt = await GetPrompt()
      if (prompt) {
        currentPrompt.value = prompt
      }
    } catch (err: any) {
      console.error("Failed to get prompt:", err)
    }
    
    if (selectedApiIndex.value !== null) {
      onApiSelected(true)
    }
  } catch (err: any) {
    console.error("Failed to load config after verification:", err)
    ElMessage.warning(`配置文件加载失败：${err?.toString?.() || 'unknown error'}`)
  }
}

const cancelStartupPass = () => {
  Quit()
}

const showStartupPassSettings = () => {
  newStartupPass.value = ''
  newStartupPassConfirm.value = ''
  startupPassSettingsModalVisible.value = true
}

const saveStartupPass = async () => {
  if (newStartupPass.value === '') {
    ElMessage.warning('请输入启动密码')
    return
  }
  if (newStartupPass.value !== newStartupPassConfirm.value) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  try {
    await SetStartupPass(newStartupPass.value)
    ElMessage.success('启动密码设置成功')
    startupPassSettingsModalVisible.value = false
    newStartupPass.value = ''
    newStartupPassConfirm.value = ''
  } catch (err: any) {
    console.error("Failed to set startup pass:", err)
    ElMessage.error('设置启动密码失败：' + (err?.toString?.() || 'unknown error'))
  }
}

// 清除启动密码
const clearStartupPass = async () => {
  try {
    await ElMessageBox.confirm(
      t('confirmClearEncryption.message'),
      t('confirmClearEncryption.title'),
      {
        confirmButtonText: t('confirmClearEncryption.confirm'),
        cancelButtonText: t('confirmClearEncryption.cancel'),
        type: 'warning',
      }
    )
    
    await ClearStartupPass()
    ElMessage.success(t('messages.encryptionCleared'))
    startupPassSettingsModalVisible.value = false
    newStartupPass.value = ''
    newStartupPassConfirm.value = ''
  } catch {
    // 用户取消操作，不做任何处理
  }
}

const cancelStartupPassSettings = () => {
  newStartupPass.value = ''
  newStartupPassConfirm.value = ''
  startupPassSettingsModalVisible.value = false
}

const handleBeforeUnload = () => {
  void saveConfig()
}

onMounted(async () => {
  await loadInitialState()
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onBeforeUnmount(async () => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
  // 关闭前自动保存
  await saveConfig()
})

watch(
  () => addForm.protocol,
  (newProtocol, oldProtocol) => {
    const oldDefault = getDefaultBaseUrl(oldProtocol || '')
    if (!addForm.baseUrl.trim() || addForm.baseUrl === oldDefault) {
      addForm.baseUrl = getDefaultBaseUrl(newProtocol)
    }
  }
)

watch(
  validKeysList,
  list => {
    if (selectedApiIndex.value !== null && selectedApiIndex.value >= list.length) {
      selectedApiIndex.value = list.length > 0 ? 0 : null
      onApiSelected(true)
    }
  },
  { deep: true }
)

const addKey = async () => {
  if (!addForm.apiKey) {
    ElMessage.warning('请输入 API Key')
    return
  }
  if (!addForm.alias) {
    ElMessage.warning('请输入别名')
    return
  }
  // 别名查重
  if (keysList.value.some(k => k.alias === addForm.alias)) {
    ElMessage.warning('该别名已存在，请使用其他别名以区分')
    return
  }

  // 支持多行粘贴
  const keys = addForm.apiKey.split('\n').map(k => k.trim()).filter(k => k)
  
  if (keys.length > 1) {
    ElMessage.warning('目前带有别名的添加方式每次只建议添加单条，多条将使用相同别名')
  }

  // 先存储明文，保存时由后端加密
  keys.forEach((k, index) => {
    let finalAlias = addForm.alias
    if (keys.length > 1) {
      finalAlias = `${addForm.alias}-${index + 1}`
    }

    keysList.value.push({
      protocol: addForm.protocol,
      alias: finalAlias,
      apiKey: k,
      baseUrl: getEffectiveBaseUrl(addForm.protocol, addForm.baseUrl),
      success: false,
      models: [],
      modelCount: 0,
      errorMsg: '',
      testing: false
    })
  })
  
  addForm.apiKey = ''
  addForm.alias = ''
  addForm.baseUrl = getDefaultBaseUrl(addForm.protocol)
  ElMessage.success(`添加了 ${keys.length} 个 Key`)
  
  // 添加后自动保存
  await saveConfig()
}

const showPromptSettings = () => {
  currentPrompt.value = ''
  void GetPrompt().then((prompt: string) => {
    currentPrompt.value = prompt
    promptModalVisible.value = true
  }).catch((err: any) => {
    console.error("Failed to get prompt:", err)
    ElMessage.error('获取提示词失败：' + (err?.toString?.() || 'unknown error'))
  })
}

const savePromptSettings = async () => {
  if (!currentPrompt.value.trim()) {
    ElMessage.warning('提示词不能为空')
    return
  }

  try {
    await SetPrompt(currentPrompt.value.trim())
    await saveConfig()
    ElMessage.success('提示词已保存')
    promptModalVisible.value = false
  } catch (err: unknown) {
    console.error("Failed to save prompt:", err)
    const errorMsg = err instanceof Error ? err.message : String(err)
    ElMessage.error('保存失败：' + errorMsg)
  }
}

const cancelPromptSettings = () => {
  promptModalVisible.value = false
  currentPrompt.value = ''
}

const saveConfig = async () => {
  try {
    const updatedConfig = await SaveAppConfig(buildConfigPayload())
    applyConfigPayload(updatedConfig)
    ElMessage.success('配置已保存')
  } catch (err: unknown) {
    console.error("Failed to save config:", err)
    const errorMsg = err instanceof Error ? err.message : String(err)
    ElMessage.error('保存失败：' + errorMsg)
  }
}

const clearList = async () => {
  if (keysList.value.length === 0) {
    ElMessage.info(t('messages.emptyList'))
    return
  }
  
  try {
    await ElMessageBox.confirm(
      t('confirmClear.message'),
      t('confirmClear.title'),
      {
        confirmButtonText: t('confirmClear.confirm'),
        cancelButtonText: t('confirmClear.cancel'),
        type: 'warning',
      }
    )
    
    keysList.value = []
    selectedApiIndex.value = null
    selectedModel.value = ''
    chatHistory.value = []
    ElMessage.success(t('messages.listCleared'))
  } catch {
    // 用户取消操作，不做任何处理
  }
}

const exportList = async () => {
  if (keysList.value.length === 0) return
  try {
    const encryptedData = await ExportEncryptedList()
    if (!encryptedData || encryptedData.length === 0) return
    
    let csv = 'Protocol,Alias,API Key,Base URL,Status,Model Count,Error\n'
    encryptedData.forEach(row => {
      csv += `${row.protocol},${row.alias},${row.apiKey},${row.baseUrl},${row.status},${row.modelCount},"${row.errorMsg.replace(/"/g, '""')}"\n`
    })
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'api_keys_export.csv'
    link.click()
  } catch (err: unknown) {
    console.error("Export failed:", err)
    const errorMsg = err instanceof Error ? err.message : String(err)
    ElMessage.error('导出失败：' + errorMsg)
  }
}

const deleteRow = (index: number) => {
  const row = keysList.value[index]
  keysList.value.splice(index, 1)
  if (validKeysList.value.length === 0) {
    selectedApiIndex.value = null
    selectedModel.value = ''
  }
  
  // 从模型缓存中删除对应的记录
  if (row) {
    const apiKey = `${row.protocol}@${row.alias}`
    if (modelCache.value.apiModels && modelCache.value.apiModels[apiKey]) {
      delete modelCache.value.apiModels[apiKey]
      void SaveModelCache(modelCache.value)
    }
  }
  
  // 保存配置
  void saveConfig()
}

const copyKey = async (text: string) => {
  if (!text) return
  try {
    // 尝试解密 API Key
    const decrypted = await DecryptAPIKey(text)
    await navigator.clipboard.writeText(decrypted)
    ElMessage.success('已复制解密的 API Key 到剪贴板')
  } catch (err: unknown) {
    // 解密失败，可能是明文或非加密数据
    if (text.startsWith('aesgcm:')) {
      const errorMsg = err instanceof Error ? err.message : String(err)
      ElMessage.error('API Key 解密失败：' + errorMsg)
    } else {
      // 可能是明文，直接复制
      await navigator.clipboard.writeText(text)
      ElMessage.success('已复制到剪贴板')
    }
  }
}

const handleCellDblClick = (row: any, column: any) => {
  if (column.property === 'apiKey') {
    copyKey(row.apiKey)
  } else if (column.property === 'baseUrl') {
    copyKey(row.baseUrl)
  }
}

const maskKey = (key: string) => {
  if (!key) return ''
  if (key.length <= 8) return '***'
  return key.substring(0, 4) + '...' + key.substring(key.length - 4)
}

const getStatusType = (success: boolean) => {
  if (success) return 'success'
  return 'danger'
}

const isZh = computed(() => currentLang.value === 'zh')

const getStatusText = (success: boolean, errorMsg: string) => {
  if (success) return isZh ? t('success') : 'Success'
  if (errorMsg) return isZh ? t('failed') : 'Failed'
  return isZh ? t('notTested') : 'N/A'
}

const testSingle = async (row: KeyRow) => {
  // 移除检测按钮的转圈动作
  // row.testing = true
  row.success = false
  row.errorMsg = ''
  try {
    row.baseUrl = getEffectiveBaseUrl(row.protocol, row.baseUrl)
    const res = await CheckKey(row.protocol, row.alias, row.apiKey, row.baseUrl)
    updateRowWithResult(row, res)
    
    // 更新模型缓存
    await updateModelCache(row)
  } catch (err: any) {
    row.success = false
    row.errorMsg = err.toString()
  } finally {
    // 移除检测按钮的转圈动作
    // row.testing = false
    
    // 保存配置，确保检测结果持久化
    await saveConfig()
  }
}

const batchTest = async () => {
  if (keysList.value.length === 0) return
  testing.value = true
  
  for (const p of protocols.value) {
    const pKeys = keysList.value.filter(k => k.protocol === p)
    if (pKeys.length > 0) {
      const aliases = pKeys.map(k => k.alias)
      const keys = pKeys.map(k => k.apiKey)
      const firstBaseUrl = getEffectiveBaseUrl(p, pKeys[0].baseUrl || '')
      try {
        const res = await BatchCheckKeys(p, aliases, keys, firstBaseUrl)
        if (res && res.results) {
          res.results.forEach((r: any) => {
            const row = keysList.value.find(k => k.apiKey === r.key && k.protocol === p)
            if (row) updateRowWithResult(row, r)
          })
        }
      } catch (err) {
        console.error("Batch test error:", err)
      }
    }
  }
  testing.value = false
  ElMessage.success('批量检测完成')
  
  // 更新模型缓存
  for (const row of keysList.value) {
    await updateModelCache(row)
  }
  
  // 保存配置，确保检测结果持久化
  await saveConfig()
}

const updateRowWithResult = (row: KeyRow, res: any) => {
  row.success = res.isValid
  row.models = res.models || []
  row.errorMsg = res.errorMsg || ''
  row.modelCount = row.models.length
  
  if (selectedApiIndex.value !== null && validKeysList.value[selectedApiIndex.value]?.apiKey === row.apiKey) {
    onApiSelected(true)
  }
}

// 更新模型缓存
const updateModelCache = async (row: KeyRow) => {
  if (row.models.length === 0) return
  
  try {
    const apiKey = `${row.protocol}@${row.alias}`
    if (!modelCache.value.apiModels) {
      modelCache.value.apiModels = {}
    }
    if (!modelCache.value.apiModels[apiKey]) {
      modelCache.value.apiModels[apiKey] = {}
    }
    
    // 测试每个模型的智能度并缓存响应
    for (const model of row.models) {
      try {
        const response = await TestModel(row.protocol, row.apiKey, row.baseUrl, model, currentPrompt.value || '请用一句话解释量子纠缠')
        modelCache.value.apiModels[apiKey][model] = response
      } catch (err) {
        console.error(`Failed to test model ${model}:`, err)
        // 保存失败的测试记录
        modelCache.value.apiModels[apiKey][model] = `Error: ${err?.toString?.() || 'Unknown error'}`
      }
    }
    
    // 保存模型缓存
    await SaveModelCache(modelCache.value)
  } catch (err: any) {
    console.error("Failed to update model cache:", err)
  }
}

const showModels = async (row: KeyRow) => {
  // 确保提示词已加载
  if (!currentPrompt.value) {
    try {
      const prompt = await GetPrompt()
      if (prompt) {
        currentPrompt.value = prompt
      }
    } catch (err: any) {
      console.error("Failed to get prompt:", err)
    }
  }
  
  currentModels.value = row.models
  currentProtocol.value = row.protocol
  currentApiKey.value = row.apiKey
  currentBaseUrl.value = row.baseUrl
  modelsModalVisible.value = true
}

const handleRowChange = (row: KeyRow | undefined) => {
  if (row && row.success) {
    const idx = validKeysList.value.findIndex(k => k.apiKey === row.apiKey && k.protocol === row.protocol)
    if (idx !== -1) {
      selectedApiIndex.value = idx
      onApiSelected()
    }
  }
}

const onApiSelected = (keepCurrentModel = false) => {
  if (selectedApiIndex.value === null) {
    availableModelsForChat.value = []
    selectedModel.value = ''
    return
  }
  const api = validKeysList.value[selectedApiIndex.value]
  if (!api) {
    selectedApiIndex.value = null
    availableModelsForChat.value = []
    selectedModel.value = ''
    return
  }

  availableModelsForChat.value = api.models || []
  
  if (keepCurrentModel && selectedModel.value && availableModelsForChat.value.includes(selectedModel.value)) {
    return
  }
  if (availableModelsForChat.value.length > 0) {
    selectedModel.value = availableModelsForChat.value[0]
  } else {
    if (api.protocol === 'anthropic') selectedModel.value = "claude-3-haiku-20240307"
    else selectedModel.value = "gpt-3.5-turbo"
    
    availableModelsForChat.value = [selectedModel.value]
  }
}

const renderMarkdown = (content: string) => {
  return marked(content)
}

const sendChat = async () => {
  if (!chatInput.value.trim()) return
  if (selectedApiIndex.value === null || !selectedModel.value) {
    ElMessage.warning('请先选择有效的 API 和 对话模型')
    return
  }
  
  const userMsg = chatInput.value
  chatInput.value = ''
  chatHistory.value.push({ role: 'user', content: userMsg })
  
  scrollToBottom()
  
  chatting.value = true
  try {
    const api = validKeysList.value[selectedApiIndex.value]
    const reply = await Chat(
      api.protocol, 
      api.apiKey, 
      getEffectiveBaseUrl(api.protocol, api.baseUrl), 
      selectedModel.value, 
      userMsg
    )
    chatHistory.value.push({ role: 'assistant', content: reply })
  } catch (err: any) {
    chatHistory.value.push({ role: 'assistant', content: `**Error:** ${err.toString()}` })
  } finally {
    chatting.value = false
    scrollToBottom()
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (chatHistoryRef.value) {
      chatHistoryRef.value.scrollTop = chatHistoryRef.value.scrollHeight
    }
  })
}
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 20px;
  box-sizing: border-box;
  text-align: left;
}

.toolbar {
  margin-bottom: 20px;
  background: #f5f7fa;
  padding: 15px;
  border-radius: 8px;
}

.table-container {
  flex: 1;
  overflow: auto;
  margin-bottom: 20px;
}

.chat-container {
  height: 300px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  background: #fff;
  text-align: left;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  flex-wrap: wrap;
  gap: 10px;
  padding: 10px 15px;
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
}

.chat-history {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
  text-align: left;
}

.chat-msg {
  margin-bottom: 15px;
  text-align: left;
}

.chat-msg.user {
  text-align: right;
}

.chat-msg.user .msg-content {
  background: #c6e2ff;
  color: #303133;
  display: inline-block;
}

.chat-msg.assistant .msg-content {
  background: #f4f4f5;
  color: #303133;
  display: inline-block;
}

.msg-role {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.msg-content {
  padding: 8px 12px;
  border-radius: 4px;
  max-width: 80%;
  text-align: left;
}

.chat-input {
  padding: 10px;
  border-top: 1px solid #e4e7ed;
  display: flex;
  align-items: flex-end;
  gap: 10px;
}

.chat-textarea {
  flex: 1;
}

:global(body) {
  margin: 0;
  background-color: #ffffff;
}

.chat-input :deep(.el-textarea__wrapper) {
  width: 100%;
}

.chat-input :deep(.el-textarea) {
  width: 100%;
}

.chat-input :deep(.el-textarea__inner) {
  text-align: left !important;
  padding-left: 12px;
}

.baseurl-select :deep(.el-select-dropdown__item) {
  text-align: left !important;
  padding-left: 20px !important;
}

:global(.markdown-body p) {
  margin: 0;
}

:global(.markdown-body) {
  color: #303133 !important;
}
</style>
