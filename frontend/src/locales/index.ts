// 语言类型定义
export type LangType = 'zh' | 'en'

// 翻译字典类型定义
export interface TranslationMessages {
  saveConfig: string
  batchTest: string
  exportList: string
  clearList: string
  switchLang: string
  add: string
  protocol: string
  alias: string
  apiKey: string
  baseUrl: string
  selectProtocol: string
  selectOrInputBaseUrl: string
  status: string
  modelCount: string
  errorMsg: string
  operation: string
  test: string
  models: string
  copy: string
  chat: string
  selectModel: string
  modelList: string
  close: string
  send: string
  modelId: string
  testResult: string
  testPrompt: string
  cacheResults: string
  debugConsole: string
  selectApi: string
  selectDialogModel: string
  you: string
  delete: string
  success: string
  failed: string
  notTested: string
  countZh: string
  operationDelete: string
  startupPass: string
  startupPassTitle: string
  startupPassSettings: string
  promptSettings: string
  prompt: string
  newPassword: string
  confirmPassword: string
  clearEncryption: string
  confirm: string
  cancel: string
  placeholder: {
    alias: string
    apiKey: string
    selectBaseUrl: string
    chatInput: string
    testPrompt: string
    startupPass: string
    prompt: string
    newPassword: string
    confirmPassword: string
  }
  confirmClear: {
    title: string
    message: string
    confirm: string
    cancel: string
  }
  confirmClearEncryption: {
    title: string
    message: string
    confirm: string
    cancel: string
  }
  messages: {
    emptyList: string
    listCleared: string
    configSaved: string
    saveFailed: string
    added: string
    keys: string
    copySuccess: string
    decryptFailed: string
    copySuccessPlain: string
    exportFailed: string
    batchTestComplete: string
    inputApiKey: string
    inputAlias: string
    aliasExists: string
    multiKeyWarning: string
    encryptionCleared: string
  }
}

// 中文翻译
const zhCN: TranslationMessages = {
  saveConfig: '保存配置',
  batchTest: '批量检测',
  exportList: '导出列表',
  clearList: '清空列表',
  switchLang: '切换语言',
  add: '添加',
  protocol: '协议',
  alias: '别名',
  apiKey: 'API Key',
  baseUrl: 'Base URL',
  selectProtocol: '选择协议',
  selectOrInputBaseUrl: '选择或输入 Base URL',
  status: '状态',
  modelCount: '模型数',
  errorMsg: '错误信息',
  operation: '操作',
  test: '检测',
  models: '模型',
  copy: '复制',
  chat: '聊天',
  selectModel: '选择模型',
  modelList: '模型列表与智能度测试',
  close: '关闭',
  send: '发送',
  modelId: '模型 ID',
  testResult: '智能度测试结果',
  testPrompt: '测试提示词',
  cacheResults: '缓存记录',
  debugConsole: '调试窗口',
  selectApi: '选择 API',
  selectDialogModel: '选择对话模型',
  you: '你',
  delete: '删除',
  success: '成功',
  failed: '失败',
  notTested: '未测',
  countZh: '个',
  operationDelete: 'Delete',
  startupPass: '启动密码',
  startupPassTitle: '请输入启动密码',
  startupPassSettings: '设置启动密码',
  promptSettings: '提问设置',
  prompt: '智能提示词',
  newPassword: '新密码',
  confirmPassword: '确认密码',
  clearEncryption: '清除加密',
  confirm: '确定',
  cancel: '取消',
  placeholder: {
    alias: '输入别名',
    apiKey: '输入 API Key',
    selectBaseUrl: '选择或输入 Base URL',
    chatInput: '输入消息...',
    testPrompt: '自定义测试智能度的提示词',
    startupPass: '请输入启动密码',
    prompt: '请输入智能提示词...',
    newPassword: '请输入新密码',
    confirmPassword: '请再次输入密码确认'
  },
  confirmClear: {
    title: '确认清空',
    message: '确定要清空列表吗？此操作将删除所有已添加的 API Key，且不可恢复。',
    confirm: '确定',
    cancel: '取消',
  },
  confirmClearEncryption: {
    title: '清除加密',
    message: '确定要清除启动密码吗？这将使用机器码作为加密密码。',
    confirm: '确定',
    cancel: '取消',
  },
  messages: {
    emptyList: '列表为空，无需清空',
    listCleared: '列表已清空',
    configSaved: '配置已保存',
    saveFailed: '保存失败：',
    added: '添加了',
    keys: '个 Key',
    copySuccess: '已复制解密的 API Key 到剪贴板',
    decryptFailed: 'API Key 解密失败：',
    copySuccessPlain: '已复制到剪贴板',
    exportFailed: '导出失败：',
    batchTestComplete: '批量测试完成',
    inputApiKey: '请输入 API Key',
    inputAlias: '请输入别名',
    aliasExists: '该别名已存在，请使用其他别名以区分',
    multiKeyWarning: '目前带有别名的添加方式每次只建议添加单条，多条将使用相同别名',
    encryptionCleared: '启动密码已清除，现在使用机器码作为加密密码',
  }
}

// 英文翻译
const enUS: TranslationMessages = {
  saveConfig: 'Save Config',
  batchTest: 'Batch Test',
  exportList: 'Export List',
  clearList: 'Clear List',
  switchLang: 'Switch Language',
  add: 'Add',
  protocol: 'Protocol',
  alias: 'Alias',
  apiKey: 'API Key',
  baseUrl: 'Base URL',
  selectProtocol: 'Select Protocol',
  selectOrInputBaseUrl: 'Select or Input Base URL',
  status: 'Status',
  modelCount: 'Models',
  errorMsg: 'Error',
  operation: 'Operation',
  test: 'Test',
  models: 'Models',
  copy: 'Copy',
  chat: 'Chat',
  selectModel: 'Select Model',
  modelList: 'Model List & Intelligence Test',
  close: 'Close',
  send: 'Send',
  modelId: 'Model ID',
  testResult: 'Intelligence Test Result',
  testPrompt: 'Test Prompt',
  cacheResults: 'Cache Results',
  debugConsole: 'Debug Console',
  selectApi: 'Select API',
  selectDialogModel: 'Select Model',
  you: 'You',
  delete: 'Delete',
  success: 'Success',
  failed: 'Failed',
  notTested: 'N/A',
  countZh: '',
  operationDelete: 'Delete',
  startupPass: 'Startup Password',
  startupPassTitle: 'Please Enter Startup Password',
  startupPassSettings: 'Set Startup Password',
  promptSettings: 'Prompt Settings',
  prompt: 'Smart Prompt',
  newPassword: 'New Password',
  confirmPassword: 'Confirm Password',
  clearEncryption: 'Clear Encryption',
  confirm: 'Confirm',
  cancel: 'Cancel',
  placeholder: {
    alias: 'Enter alias',
    apiKey: 'Enter API Key',
    selectBaseUrl: 'Select or Input Base URL',
    chatInput: 'Type a message...',
    testPrompt: 'Custom prompt for testing intelligence',
    startupPass: 'Enter startup password',
    prompt: 'Enter smart prompt...',
    newPassword: 'Enter new password',
    confirmPassword: 'Confirm password'
  },
  confirmClear: {
    title: 'Confirm Clear',
    message: 'Are you sure you want to clear the list? This will delete all added API Keys and cannot be undone.',
    confirm: 'Confirm',
    cancel: 'Cancel',
  },
  confirmClearEncryption: {
    title: 'Clear Encryption',
    message: 'Are you sure you want to clear the startup password? This will use the machine code as the encryption password.',
    confirm: 'Confirm',
    cancel: 'Cancel',
  },
  messages: {
    emptyList: 'List is empty, nothing to clear',
    listCleared: 'List cleared',
    configSaved: 'Config saved',
    saveFailed: 'Save failed: ',
    added: 'Added',
    keys: ' key(s)',
    copySuccess: 'Decrypted API Key copied to clipboard',
    decryptFailed: 'API Key decryption failed: ',
    copySuccessPlain: 'Copied to clipboard',
    exportFailed: 'Export failed: ',
    batchTestComplete: 'Batch test completed',
    inputApiKey: 'Please enter API Key',
    inputAlias: 'Please enter alias',
    aliasExists: 'This alias already exists, please use another alias',
    multiKeyWarning: 'Adding with alias is recommended for single entry only, multiple entries will use the same alias',
    encryptionCleared: 'Startup password cleared, now using machine code as encryption password',
  }
}

// 翻译字典
export const translations: Record<LangType, TranslationMessages> = {
  zh: zhCN,
  en: enUS
}

// 创建翻译函数
export function createTranslator(lang: LangType) {
  return function t(key: string): string {
    const keys = key.split('.')
    let value: any = translations[lang]
    for (const k of keys) {
      value = value?.[k]
    }
    return value || key
  }
}

// 导出翻译消息类型
export type { TranslationMessages as TranslationMessagesType }
