import { createI18n } from 'vue-i18n'
import en from './en'
import zhCn from './zh-cn'

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('aegis_locale') || 'zh-cn',
  fallbackLocale: 'en',
  messages: { en, 'zh-cn': zhCn }
})

export default i18n
