import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../common';
import { BasicConfigSection } from './sections/BasicConfigSection';
import { GinConfigSection } from './sections/GinConfigSection';
import { DatabaseConfigSection } from './sections/DatabaseConfigSection';
import { RedisConfigSection } from './sections/RedisConfigSection';
import { K8sConfigSection } from './sections/K8sConfigSection';
import { NFSConfigSection } from './sections/NFSConfigSection';
import { CheatConfigSection } from './sections/CheatConfigSection';
import { WebhookConfigSection } from './sections/WebhookConfigSection.jsx';

const tabs = [
  { key: 'basic', i18nKey: 'admin.system.sections.basic', Component: BasicConfigSection },
  { key: 'gin', i18nKey: 'admin.system.sections.gin', Component: GinConfigSection },
  { key: 'database', i18nKey: 'admin.system.sections.database', Component: DatabaseConfigSection },
  { key: 'redis', i18nKey: 'admin.system.sections.redis', Component: RedisConfigSection },
  { key: 'k8s', i18nKey: 'admin.system.sections.k8s', Component: K8sConfigSection },
  { key: 'nfs', i18nKey: 'admin.system.sections.nfs', Component: NFSConfigSection },
  { key: 'webhook', i18nKey: 'admin.system.sections.webhook', Component: WebhookConfigSection },
  { key: 'cheat', i18nKey: 'admin.system.sections.cheat', Component: CheatConfigSection },
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

      <div className="flex flex-wrap border-b border-neutral-700">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            className={`px-4 py-2 text-sm font-medium transition-colors ${
              activeTab === tab.key
                ? 'text-blue-400 border-b-2 border-blue-400'
                : 'text-neutral-400 hover:text-neutral-300'
            }`}
            onClick={() => setActiveTab(tab.key)}
          >
            {t(tab.i18nKey)}
          </button>
        ))}
      </div>

      {ActiveComponent && <ActiveComponent config={config} updateConfig={updateConfig} />}
    </div>
  );
}
