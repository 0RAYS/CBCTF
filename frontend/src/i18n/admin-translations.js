import i18n from '@/i18n/index';
import en from '@/i18n/locales/en.json?scope=admin';
import zhCN from '@/i18n/locales/zh-CN.json?scope=admin';

// Evaluated before either lazy admin layout renders, including direct deep links.
i18n.addResourceBundle('en', 'translation', en, true, true);
i18n.addResourceBundle('zh-CN', 'translation', zhCN, true, true);
