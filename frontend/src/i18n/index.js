import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import en from './locales/en.json';
import zhCN from './locales/zh-CN.json';

const supportedLanguages = ['en', 'zh-CN'];

const getInitialLanguage = () => {
  if (typeof window === 'undefined') return 'en';
  const stored = window.localStorage.getItem('language');
  if (stored && supportedLanguages.includes(stored)) return stored;
  const browser = window.navigator.language || 'en';
  if (supportedLanguages.includes(browser)) return browser;
  if (browser.toLowerCase().startsWith('zh')) return 'zh-CN';
  return 'en';
};

i18n.use(initReactI18next).init({
  resources: {
    en: { translation: en },
    'zh-CN': { translation: zhCN },
  },
  lng: getInitialLanguage(),
  fallbackLng: 'en',
  interpolation: { escapeValue: false },
});

const LANG_ATTR_MAP = {
  'zh-CN': 'zh-Hans',
  en: 'en',
};

// 初始化时立即设置正确的 lang 属性
if (typeof window !== 'undefined') {
  const initialLang = getInitialLanguage();
  document.documentElement.lang = LANG_ATTR_MAP[initialLang] || initialLang;
}

export const setLanguage = (lng) => {
  if (!supportedLanguages.includes(lng)) return;
  i18n.changeLanguage(lng);
  if (typeof window !== 'undefined') {
    window.localStorage.setItem('language', lng);
    document.documentElement.lang = LANG_ATTR_MAP[lng] || lng;
  }
};

export default i18n;
