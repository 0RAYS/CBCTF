import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Tabs } from '../../../common';
import { BasicConfigSection } from './sections/BasicConfigSection';
import { GinConfigSection } from './sections/GinConfigSection';
import { DatabaseConfigSection } from './sections/DatabaseConfigSection';
import { RedisConfigSection } from './sections/RedisConfigSection';
import { K8sConfigSection } from './sections/K8sConfigSection';
import { NFSConfigSection } from './sections/NFSConfigSection';
import { CheatConfigSection } from './sections/CheatConfigSection';
import { WebhookConfigSection } from './sections/WebhookConfigSection.jsx';
import { RegistrationConfigSection } from './sections/RegistrationConfigSection';

const tabs = [
  { key: 'basic', i18nKey: 'admin.system.sections.basic', Component: BasicConfigSection },
  { key: 'gin', i18nKey: 'admin.system.sections.gin', Component: GinConfigSection },
  { key: 'database', i18nKey: 'admin.system.sections.database', Component: DatabaseConfigSection },
  { key: 'redis', i18nKey: 'admin.system.sections.redis', Component: RedisConfigSection },
  { key: 'k8s', i18nKey: 'admin.system.sections.k8s', Component: K8sConfigSection },
  { key: 'nfs', i18nKey: 'admin.system.sections.nfs', Component: NFSConfigSection },
  { key: 'webhook', i18nKey: 'admin.system.sections.webhook', Component: WebhookConfigSection },
  { key: 'cheat', i18nKey: 'admin.system.sections.cheat', Component: CheatConfigSection },
  { key: 'registration', i18nKey: 'admin.system.sections.registration', Component: RegistrationConfigSection },
];

export function SystemConfigForm({ config, updateConfig, onUpdate, onRestart, isUpdating, isRestarting }) {
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

      {ActiveComponent && <ActiveComponent config={config} updateConfig={updateConfig} />}
    </div>
  );
}
