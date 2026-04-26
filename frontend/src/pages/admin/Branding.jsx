import { useEffect, useRef, useState } from 'react';
import { useDispatch } from 'react-redux';
import { motion } from 'motion/react';
import { Button, Card, Input, Textarea } from '../../components/common';
import { getAdminBranding, updateAdminBranding, uploadBrandingLogo } from '../../api/admin/branding';
import { toast } from '../../utils/toast';
import { useImmerState } from '../../hooks/useImmerState';
import { mergeBranding } from '../../config/branding';
import { fetchBranding, setBranding } from '../../store/branding';
import { useTranslation } from 'react-i18next';

function LocalizedField({ label, value, onChange, placeholderZh, placeholderEn, multiline = false, rows = 3 }) {
  const FieldComponent = multiline ? Textarea : Input;

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <FieldComponent
        label={`${label} · 中文`}
        value={value?.zh_cn || ''}
        placeholder={placeholderZh}
        onChange={(event) => onChange('zh_cn', event.target.value)}
        rows={multiline ? rows : undefined}
      />
      <FieldComponent
        label={`${label} · English`}
        value={value?.en || ''}
        placeholder={placeholderEn}
        onChange={(event) => onChange('en', event.target.value)}
        rows={multiline ? rows : undefined}
      />
    </div>
  );
}

function BrandingSettings() {
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const [branding, updateBrandingState] = useImmerState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef(null);

  const fetchCurrentBranding = async () => {
    setLoading(true);
    try {
      const response = await getAdminBranding();
      if (response.code === 200) {
        const normalized = mergeBranding(response.data);
        updateBrandingState(() => normalized);
        dispatch(setBranding(normalized));
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.branding.toast.fetchFailed') });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCurrentBranding();
  }, []);

  const updateLocalizedField = (path, language, nextValue) => {
    updateBrandingState((draft) => {
      let target = draft;
      for (const key of path) target = target[key];
      target[language] = nextValue;
    });
  };

  const handleSave = async () => {
    if (!branding) return;
    setSaving(true);
    try {
      const payload = {
        site_name: branding.site_name,
        admin_name: branding.admin_name,
        browser_title: branding.browser_title,
        browser_description: branding.browser_description,
        footer_copyright: branding.footer_copyright,
        home_logo_alt: branding.home_logo_alt,
        home: branding.home,
      };
      const response = await updateAdminBranding(payload);
      if (response.code === 200) {
        const normalized = mergeBranding(response.data);
        updateBrandingState(() => normalized);
        dispatch(setBranding(normalized));
        toast.success({ description: t('admin.branding.toast.updateSuccess') });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.branding.toast.updateFailed') });
    } finally {
      setSaving(false);
    }
  };

  const handleLogoChange = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file) return;

    setUploading(true);
    try {
      const response = await uploadBrandingLogo(file);
      if (response.code === 200) {
        toast.success({ description: t('admin.branding.toast.logoUpdated') });
        await fetchCurrentBranding();
        dispatch(fetchBranding());
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.branding.toast.logoUpdateFailed') });
    } finally {
      setUploading(false);
    }
  };

  if (loading || !branding) {
    return <div className="p-4 text-neutral-400">{t('common.loading')}</div>;
  }

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: 16 }} animate={{ opacity: 1, y: 0 }} className="space-y-6">
        <Card padding="lg" className="space-y-6">
          <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4">
            <div>
              <h2 className="text-lg font-mono text-neutral-50">{t('admin.branding.sections.identity')}</h2>
              <p className="text-sm text-neutral-400 mt-2">{t('admin.branding.identityHint')}</p>
            </div>
            <div className="flex gap-3">
              <Button variant="outline" onClick={() => fileInputRef.current?.click()} loading={uploading}>
                {t('admin.branding.actions.uploadLogo')}
              </Button>
              <Button variant="primary" onClick={handleSave} loading={saving}>
                {t('common.save')}
              </Button>
            </div>
          </div>

          <div className="grid grid-cols-1 xl:grid-cols-[320px_minmax(0,1fr)] gap-6">
            <div className="border border-neutral-300/20 rounded-md bg-black/20 p-4">
              <div className="text-sm text-neutral-400 mb-3">{t('admin.branding.labels.homeLogo')}</div>
              <div className="aspect-square rounded-md border border-neutral-300/20 bg-neutral-950 flex items-center justify-center overflow-hidden">
                <img
                  src={branding.home_logo}
                  alt={branding.home_logo_alt?.zh_cn || 'logo'}
                  className="w-full h-full object-contain"
                />
              </div>
              <div className="text-xs text-neutral-500 mt-3 break-all">{branding.home_logo}</div>
            </div>

            <div className="space-y-5">
              <LocalizedField
                label={t('admin.branding.labels.siteName')}
                value={branding.site_name}
                placeholderZh="深潜 CTF"
                placeholderEn="DEEP DIVE CTF"
                onChange={(language, value) => updateLocalizedField(['site_name'], language, value)}
              />
              <LocalizedField
                label={t('admin.branding.labels.adminName')}
                value={branding.admin_name}
                placeholderZh="深潜管理台"
                placeholderEn="DEEP DIVE Admin"
                onChange={(language, value) => updateLocalizedField(['admin_name'], language, value)}
              />
              <LocalizedField
                label={t('admin.branding.labels.browserTitle')}
                value={branding.browser_title}
                placeholderZh="深潜 CTF"
                placeholderEn="DEEP DIVE CTF"
                onChange={(language, value) => updateLocalizedField(['browser_title'], language, value)}
              />
              <LocalizedField
                label={t('admin.branding.labels.browserDescription')}
                value={branding.browser_description}
                placeholderZh="深潜 CTF 网络安全竞赛平台"
                placeholderEn="DEEP DIVE CTF competition platform"
                multiline
                rows={3}
                onChange={(language, value) => updateLocalizedField(['browser_description'], language, value)}
              />
              <LocalizedField
                label={t('admin.branding.labels.footerCopyright')}
                value={branding.footer_copyright}
                placeholderZh="© 2025 深潜 CTF"
                placeholderEn="© 2025 DEEP DIVE CTF"
                onChange={(language, value) => updateLocalizedField(['footer_copyright'], language, value)}
              />
              <LocalizedField
                label={t('admin.branding.labels.homeLogoAlt')}
                value={branding.home_logo_alt}
                placeholderZh="深潜 CTF 首页 Logo"
                placeholderEn="DEEP DIVE CTF home logo"
                onChange={(language, value) => updateLocalizedField(['home_logo_alt'], language, value)}
              />
            </div>
          </div>
        </Card>

        <Card padding="lg" className="space-y-5">
          <h2 className="text-lg font-mono text-neutral-50">{t('admin.branding.sections.homeHero')}</h2>
          <LocalizedField
            label={t('admin.branding.labels.heroTitlePrefix')}
            value={branding.home.hero.title_prefix}
            placeholderZh="深入探索"
            placeholderEn="Dive Deep into the"
            onChange={(language, value) => updateLocalizedField(['home', 'hero', 'title_prefix'], language, value)}
          />
          <LocalizedField
            label={t('admin.branding.labels.heroTitleHighlight')}
            value={branding.home.hero.title_highlight}
            placeholderZh="网络安全"
            placeholderEn="Cyber Security"
            onChange={(language, value) => updateLocalizedField(['home', 'hero', 'title_highlight'], language, value)}
          />
          <LocalizedField
            label={t('admin.branding.labels.heroTitleSuffix')}
            value={branding.home.hero.title_suffix}
            placeholderZh="挑战"
            placeholderEn="Challenge"
            onChange={(language, value) => updateLocalizedField(['home', 'hero', 'title_suffix'], language, value)}
          />
          <LocalizedField
            label={t('admin.branding.labels.heroSubtitle')}
            value={branding.home.hero.subtitle}
            placeholderZh="加入深潜 CTF 社区, 在真实场景中对抗、练习与成长"
            placeholderEn="Join the elite community of hackers..."
            multiline
            rows={4}
            onChange={(language, value) => updateLocalizedField(['home', 'hero', 'subtitle'], language, value)}
          />
          <LocalizedField
            label={t('admin.branding.labels.heroPrimaryAction')}
            value={branding.home.hero.primary_action}
            placeholderZh="立即参赛"
            placeholderEn="START HACKING"
            onChange={(language, value) => updateLocalizedField(['home', 'hero', 'primary_action'], language, value)}
          />
          <LocalizedField
            label={t('admin.branding.labels.heroSecondaryAction')}
            value={branding.home.hero.secondary_action}
            placeholderZh="了解更多"
            placeholderEn="LEARN MORE"
            onChange={(language, value) => updateLocalizedField(['home', 'hero', 'secondary_action'], language, value)}
          />
        </Card>

        <Card padding="lg" className="space-y-6">
          <h2 className="text-lg font-mono text-neutral-50">{t('admin.branding.sections.homeSections')}</h2>

          {[
            ['challenge_types', t('admin.branding.sections.challengeTypes')],
            ['upcoming', t('admin.branding.sections.upcoming')],
            ['leaderboard', t('admin.branding.sections.leaderboard')],
          ].map(([sectionKey, sectionLabel]) => (
            <div key={sectionKey} className="space-y-5 border border-neutral-300/15 rounded-md p-4 bg-black/15">
              <div className="text-sm font-mono text-neutral-300">{sectionLabel}</div>
              <LocalizedField
                label={t('admin.branding.labels.sectionTitlePrefix')}
                value={branding.home[sectionKey].title_prefix}
                placeholderZh="标题前缀"
                placeholderEn="Title prefix"
                onChange={(language, value) =>
                  updateLocalizedField(['home', sectionKey, 'title_prefix'], language, value)
                }
              />
              <LocalizedField
                label={t('admin.branding.labels.sectionTitleHighlight')}
                value={branding.home[sectionKey].title_highlight}
                placeholderZh="标题高亮"
                placeholderEn="Title highlight"
                onChange={(language, value) =>
                  updateLocalizedField(['home', sectionKey, 'title_highlight'], language, value)
                }
              />
              <LocalizedField
                label={t('admin.branding.labels.sectionSubtitle')}
                value={branding.home[sectionKey].subtitle}
                placeholderZh="区块说明"
                placeholderEn="Section subtitle"
                multiline
                rows={3}
                onChange={(language, value) => updateLocalizedField(['home', sectionKey, 'subtitle'], language, value)}
              />
              {sectionKey !== 'challenge_types' && (
                <LocalizedField
                  label={t('admin.branding.labels.sectionAction')}
                  value={branding.home[sectionKey].action}
                  placeholderZh="按钮文案"
                  placeholderEn="CTA label"
                  onChange={(language, value) => updateLocalizedField(['home', sectionKey, 'action'], language, value)}
                />
              )}
            </div>
          ))}
        </Card>
      </motion.div>

      <input
        ref={fileInputRef}
        type="file"
        className="hidden"
		accept="image/png,image/jpeg,image/jpg,image/gif"
        onChange={handleLogoChange}
      />
    </div>
  );
}

export default BrandingSettings;
