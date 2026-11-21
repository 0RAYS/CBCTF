import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Tabs } from '@/components/common';
import { BasicConfigSection } from '@/components/features/Admin/SystemConfig/sections/BasicConfigSection';
import { GinConfigSection } from '@/components/features/Admin/SystemConfig/sections/GinConfigSection';
import { DatabaseConfigSection } from '@/components/features/Admin/SystemConfig/sections/DatabaseConfigSection';
import { RedisConfigSection } from '@/components/features/Admin/SystemConfig/sections/RedisConfigSection';
import { K8sConfigSection } from '@/components/features/Admin/SystemConfig/sections/K8sConfigSection';
import { CheatConfigSection } from '@/components/features/Admin/SystemConfig/sections/CheatConfigSection';
import { WebhookConfigSection } from '@/components/features/Admin/SystemConfig/sections/WebhookConfigSection.jsx';
import { RegistrationConfigSection } from '@/components/features/Admin/SystemConfig/sections/RegistrationConfigSection';

const tabs = [
  { key: 'basic', i18nKey: 'admin.system.sections.basic', Component: BasicConfigSection },
  { key: 'gin', i18nKey: 'admin.system.sections.gin', Component: GinConfigSection },
  { key: 'database', i18nKey: 'admin.system.sections.database', Component: DatabaseConfigSection },
  { key: 'redis', i18nKey: 'admin.system.sections.redis', Component: RedisConfigSection },
  { key: 'k8s', i18nKey: 'admin.system.sections.k8s', Component: K8sConfigSection },
  { key: 'webhook', i18nKey: 'admin.system.sections.webhook', Component: WebhookConfigSection },
  { key: 'cheat', i18nKey: 'admin.system.sections.cheat', Component: CheatConfigSection },
  { key: 'registration', i18nKey: 'admin.system.sections.registration', Component: RegistrationConfigSection },
];

export function SystemConfigForm({
  config,
  updateConfig,
  onUpdate,
  onUploadGeoCityDB,
  onRestart,
  isUpdating,
  isUploadingGeoCityDB,
  isRestarting,
}) {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState('basic');

  if (!config) {
    return null;
  }

  const activeTabDef = tabs.find((tab) => tab.key === activeTab);
  const ActiveComponent = activeTabDef?.Component;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-end gap-2">
        <Button variant="primary" size="sm" onClick={onUpdate} loading={isUpdating} disabled={isUpdating}>
          {t('admin.system.actions.update')}
        </Button>
        <Button variant="danger" size="sm" onClick={onRestart} disabled={isRestarting}>
          {t('admin.system.actions.restart')}
        </Button>
      </div>

      <Tabs
        variant="compact"
        wrapperClassName={null}
        value={activeTab}
        onChange={setActiveTab}
        items={tabs.map((tab) => ({ key: tab.key, label: t(tab.i18nKey) }))}
      />

      {ActiveComponent && (
        <ActiveComponent
          config={config}
          updateConfig={updateConfig}
          onUploadGeoCityDB={onUploadGeoCityDB}
          isUploadingGeoCityDB={isUploadingGeoCityDB}
        />
      )}
    </div>
  );
}
