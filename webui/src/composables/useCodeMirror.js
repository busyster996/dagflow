/**
 * CodeMirror 6 编辑器管理 Composable
 * 替代Monaco Editor，体积更小，性能更好
 */

import { ref, onUnmounted, nextTick } from 'vue'
import { EditorView, keymap, lineNumbers, highlightActiveLineGutter, highlightSpecialChars, drawSelection, dropCursor, rectangularSelection, crosshairCursor, highlightActiveLine } from '@codemirror/view'
import { EditorState, Compartment } from '@codemirror/state'
import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands'
import { indentOnInput, bracketMatching, foldGutter, foldKeymap } from '@codemirror/language'
import { yaml } from '@codemirror/lang-yaml'
import { oneDark } from '@codemirror/theme-one-dark'
import { highlightSelectionMatches, searchKeymap } from '@codemirror/search'
import { closeBrackets, closeBracketsKeymap, autocompletion, completionKeymap } from '@codemirror/autocomplete'
import { lintKeymap } from '@codemirror/lint'

/**
 * CodeMirror 编辑器管理
 * @param {Object} options - 配置选项
 * @param {string} options.language - 编辑器语言，默认 'yaml'
 * @param {string} options.theme - 编辑器主题，'light' 或 'dark'，默认 'dark'
 * @param {boolean} options.readOnly - 是否只读，默认 false
 * @param {string} options.initialValue - 初始值
 * @param {Object} options.editorOptions - 额外配置
 * @returns {Object} 编辑器管理对象
 */
export function useCodeMirror(options = {}) {
  const {
    language = 'yaml',
    theme = 'dark',
    readOnly = false,
    initialValue = '',
    editorOptions = {},
  } = options

  let editorView = null
  let isInitialized = false
  
  // 创建可配置的 compartments
  const languageConf = new Compartment()
  const readOnlyConf = new Compartment()
  const themeConf = new Compartment()

  /**
   * 获取语言扩展
   */
  const getLanguageExtension = (lang) => {
    switch (lang) {
      case 'yaml':
        return yaml()
      case 'json':
        // JSON使用JavaScript语言模式
        return yaml() // 暂时使用yaml，后续可以添加json支持
      default:
        return yaml()
    }
  }

  /**
   * 获取主题扩展
   */
  const getThemeExtension = (themeName) => {
    return themeName === 'dark' ? oneDark : []
  }

  /**
   * 基础扩展集合
   */
  const getBaseExtensions = () => [
    lineNumbers(),
    highlightActiveLineGutter(),
    highlightSpecialChars(),
    history(),
    foldGutter(),
    drawSelection(),
    dropCursor(),
    EditorState.allowMultipleSelections.of(true),
    indentOnInput(),
    bracketMatching(),
    closeBrackets(),
    autocompletion(),
    rectangularSelection(),
    crosshairCursor(),
    highlightActiveLine(),
    highlightSelectionMatches(),
    keymap.of([
      ...closeBracketsKeymap,
      ...defaultKeymap,
      ...searchKeymap,
      ...historyKeymap,
      ...foldKeymap,
      ...completionKeymap,
      ...lintKeymap,
      indentWithTab,
    ]),
  ]

  /**
   * 初始化编辑器
   * @param {HTMLElement} container - 编辑器容器元素
   * @param {string} value - 初始内容
   */
  const initEditor = async (container, value = initialValue) => {
    if (!container) {
      console.warn('CodeMirror editor container not found')
      return
    }

    // 清理已存在的编辑器
    if (editorView) {
      editorView.destroy()
      editorView = null
    }

    await nextTick()

    try {
      const state = EditorState.create({
        doc: value || '',
        extensions: [
          ...getBaseExtensions(),
          languageConf.of(getLanguageExtension(language)),
          readOnlyConf.of(EditorView.editable.of(!readOnly)),
          themeConf.of(getThemeExtension(theme)),
          EditorView.lineWrapping, // 自动换行
          ...(editorOptions.extensions || []),
        ],
      })

      editorView = new EditorView({
        state,
        parent: container,
      })

      isInitialized = true
    } catch (error) {
      console.error('Failed to initialize CodeMirror Editor:', error)
    }
  }

  /**
   * 获取编辑器内容
   * @returns {string} 编辑器内容
   */
  const getValue = () => {
    if (!editorView) return ''
    return editorView.state.doc.toString()
  }

  /**
   * 设置编辑器内容
   * @param {string} value - 新内容
   */
  const setValue = (value) => {
    if (!editorView) return
    
    const transaction = editorView.state.update({
      changes: {
        from: 0,
        to: editorView.state.doc.length,
        insert: value || '',
      },
    })
    
    editorView.dispatch(transaction)
  }

  /**
   * 追加内容
   * @param {string} text - 要追加的文本
   */
  const appendValue = (text) => {
    if (!editorView) return

    const currentValue = getValue()
    const newValue = currentValue + (currentValue ? '\n' : '') + text
    setValue(newValue)

    // 滚动到底部
    scrollToBottom()
  }

  /**
   * 清空内容
   */
  const clear = () => {
    setValue('')
  }

  /**
   * 设置只读状态
   * @param {boolean} readonly - 是否只读
   */
  const setReadOnly = (readonly) => {
    if (!editorView) return
    
    editorView.dispatch({
      effects: readOnlyConf.reconfigure(EditorView.editable.of(!readonly)),
    })
  }

  /**
   * 设置语言
   * @param {string} lang - 语言类型
   */
  const setLanguage = (lang) => {
    if (!editorView) return
    
    editorView.dispatch({
      effects: languageConf.reconfigure(getLanguageExtension(lang)),
    })
  }

  /**
   * 设置主题
   * @param {string} newTheme - 主题名称 'light' 或 'dark'
   */
  const setTheme = (newTheme) => {
    if (!editorView) return
    
    editorView.dispatch({
      effects: themeConf.reconfigure(getThemeExtension(newTheme)),
    })
  }

  /**
   * 格式化文档（CodeMirror需要手动实现，这里提供基础的缩进整理）
   */
  const formatDocument = () => {
    if (!editorView) return
    // CodeMirror 6 没有内置的格式化功能
    // 可以根据需要添加自定义格式化逻辑
    console.log('CodeMirror formatDocument - not implemented')
  }

  /**
   * 聚焦编辑器
   */
  const focus = () => {
    if (!editorView) return
    editorView.focus()
  }

  /**
   * 调整大小（CodeMirror 6自动处理）
   */
  const layout = () => {
    if (!editorView) return
    // CodeMirror 6 会自动调整大小
    editorView.requestMeasure()
  }

  /**
   * 滚动到顶部
   */
  const scrollToTop = () => {
    if (!editorView) return
    
    editorView.dispatch({
      effects: EditorView.scrollIntoView(0, { y: 'start' }),
    })
  }

  /**
   * 滚动到底部
   */
  const scrollToBottom = () => {
    if (!editorView) return
    
    const lastLine = editorView.state.doc.length
    editorView.dispatch({
      effects: EditorView.scrollIntoView(lastLine, { y: 'end' }),
    })
  }

  /**
   * 获取选中的文本
   * @returns {string} 选中的文本
   */
  const getSelectedText = () => {
    if (!editorView) return ''
    
    const selection = editorView.state.selection.main
    return editorView.state.doc.sliceString(selection.from, selection.to)
  }

  /**
   * 插入文本到当前光标位置
   * @param {string} text - 要插入的文本
   */
  const insertText = (text) => {
    if (!editorView) return
    
    const selection = editorView.state.selection.main
    editorView.dispatch({
      changes: {
        from: selection.from,
        to: selection.to,
        insert: text,
      },
      selection: { anchor: selection.from + text.length },
    })
  }

  /**
   * 监听内容变化
   * @param {Function} callback - 回调函数
   * @returns {Function} 取消监听的函数
   */
  const onContentChange = (callback) => {
    if (!editorView) return () => {}

    const updateListener = EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        callback(getValue())
      }
    })

    editorView.dispatch({
      effects: EditorState.appendConfig.of(updateListener),
    })

    return () => {
      // CodeMirror 6 的监听器会在视图销毁时自动清理
    }
  }

  /**
   * 销毁编辑器
   */
  const dispose = () => {
    if (editorView) {
      editorView.destroy()
      editorView = null
      isInitialized = false
    }
  }

  /**
   * 获取编辑器实例
   * @returns {EditorView} 编辑器实例
   */
  const getEditor = () => {
    return editorView
  }

  /**
   * 检查是否已初始化
   * @returns {boolean} 是否已初始化
   */
  const isEditorInitialized = () => {
    return isInitialized
  }

  // 组件卸载时自动销毁编辑器
  onUnmounted(() => {
    dispose()
  })

  return {
    initEditor,
    getValue,
    setValue,
    appendValue,
    clear,
    setReadOnly,
    setLanguage,
    setTheme,
    formatDocument,
    focus,
    layout,
    scrollToTop,
    scrollToBottom,
    getSelectedText,
    insertText,
    onContentChange,
    dispose,
    getEditor,
    isEditorInitialized,
  }
}