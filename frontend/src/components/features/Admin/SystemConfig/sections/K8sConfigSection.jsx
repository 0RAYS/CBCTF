import { IconCloudComputing } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigField } from '../fields/ConfigField';
import { ConfigListField } from '../fields/ConfigListField';
import { FrpServerField } from '../fields/FrpServerField';
import { Input } from '../../../../common';

const sanitizeNumber = (value, fallbackValue = 0) => {
  if (value === '' || value === null || value === undefined) {
    return fallbackValue;
  }
  const numeric = Number(value);
  return Number.isNaN(numeric) ? fallbackValue : numeric;
};

export function K8sConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.k8s')} icon={IconCloudComputing}>
      <ConfigField
        label={t('admin.system.labels.k8sConfig')}
        value={config.k8s.config}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.k8s.config = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.k8sNamespace')}
        value={config.k8s.namespace}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.k8s.namespace = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.k8sTcpdump')}
        value={config.k8s.tcpdump}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.k8s.tcpdump = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.k8sWorker')}
        type="number"
        value={config.k8s.generator_worker}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.k8s.generator_worker = sanitizeNumber(value, config.k8s.generator_worker);
          })
        }
      />

      <div className="space-y-1">
        <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.externalNetwork')} CIDR</span>
        <div className="space-y-2">
          <Input
            size="sm"
            value={config.k8s.external_network.cidr}
            onChange={(event) =>
              updateConfig((draft) => {
                draft.k8s.external_network.cidr = event.target.value;
              })
            }
            placeholder={t('admin.system.k8s.cidr')}
          />
          <span className="text-xs font-mono text-neutral-400">{t('admin.system.k8s.gateway')}</span>
          <Input
            size="sm"
            value={config.k8s.external_network.gateway}
            onChange={(event) =>
              updateConfig((draft) => {
                draft.k8s.external_network.gateway = event.target.value;
              })
            }
            placeholder={t('admin.system.k8s.gateway')}
          />
          <span className="text-xs font-mono text-neutral-400">{t('admin.system.k8s.interface')}</span>
          <Input
            size="sm"
            value={config.k8s.external_network.interface}
            onChange={(event) =>
              updateConfig((draft) => {
                draft.k8s.external_network.interface = event.target.value;
              })
            }
            placeholder={t('admin.system.k8s.interface')}
          />
          <ConfigListField
            label={t('admin.system.k8s.excludeIps')}
            items={config.k8s.external_network.exclude_ips || []}
            onAdd={() =>
              updateConfig((draft) => {
                draft.k8s.external_network.exclude_ips.push('');
              })
            }
            onUpdate={(index, value) =>
              updateConfig((draft) => {
                draft.k8s.external_network.exclude_ips[index] = value;
              })
            }
            onRemove={(index) =>
              updateConfig((draft) => {
                draft.k8s.external_network.exclude_ips.splice(index, 1);
              })
            }
          />
        </div>

        <div className="space-y-1 pt-2">
          <div className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-0.5 rounded">
            <span className="text-neutral-400">{t('admin.system.k8s.cidr')}:</span> {config.k8s.external_network.cidr}
          </div>
          <div className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-0.5 rounded">
            <span className="text-neutral-400">{t('admin.system.k8s.gateway')}:</span>{' '}
            {config.k8s.external_network.gateway}
          </div>
          <div className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-0.5 rounded">
            <span className="text-neutral-400">{t('admin.system.k8s.interface')}:</span>{' '}
            {config.k8s.external_network.interface}
          </div>
          <div className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-0.5 rounded">
            <span className="text-neutral-400">{t('admin.system.k8s.excludeIps')}:</span>{' '}
            {(config.k8s.external_network.exclude_ips || []).join(', ')}
          </div>
        </div>
      </div>

      <div className="border-b border-neutral-300/20" />

      <div className="space-y-2">
        <ConfigField
          label={t('admin.system.labels.frpOn')}
          type="boolean"
          value={config.k8s.frp.on}
          options={[
            { value: 'true', label: t('common.yes') },
            { value: 'false', label: t('common.no') },
          ]}
          onChange={(value) =>
            updateConfig((draft) => {
              draft.k8s.frp.on = value;
            })
          }
        />
        <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.frpc')}</span>
        <Input
          size="sm"
          value={config.k8s.frp.frpc}
          onChange={(event) =>
            updateConfig((draft) => {
              draft.k8s.frp.frpc = event.target.value;
            })
          }
          placeholder={t('admin.system.k8s.frp')}
        />
        <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.nginx')}</span>
        <Input
          size="sm"
          value={config.k8s.frp.nginx}
          onChange={(event) =>
            updateConfig((draft) => {
              draft.k8s.frp.nginx = event.target.value;
            })
          }
          placeholder={t('admin.system.k8s.nginx')}
        />
        <FrpServerField frpsList={config.k8s.frp.frps || []} updateConfig={updateConfig} />

        <div className="space-y-1 pt-2">
          <div className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-0.5 rounded">
            <span className="text-neutral-400">{t('admin.system.k8s.frp')}:</span> {config.k8s.frp.frpc}
          </div>
          <div className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-0.5 rounded">
            <span className="text-neutral-400">{t('admin.system.k8s.nginx')}:</span> {config.k8s.frp.nginx}
          </div>
          <div className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-0.5 rounded">
            <span className="text-neutral-400">{t('admin.system.k8s.enabled')}:</span>{' '}
            {config.k8s.frp.on ? t('common.yes') : t('common.no')}
          </div>
          {(config.k8s.frp.frps || []).map((frps, index) => (
            <div key={index} className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-0.5 rounded ml-2">
              <div>
                {t('admin.system.labels.frps')}: {frps.host}
              </div>
              <div>
                {t('admin.system.k8s.port')}: {frps.port}
              </div>
              <div>
                {t('admin.system.k8s.token')}: {frps.token}
              </div>
              {frps.allowed && (
                <div>
                  {t('admin.system.k8s.allowedPorts')}:{' '}
                  {frps.allowed
                    .map((port) => {
                      if (!port) {
                        return '';
                      }
                      const exclude =
                        Array.isArray(port.exclude) && port.exclude.length > 0
                          ? ` (ex ${port.exclude.join(', ')})`
                          : '';
                      return `${port.from}-${port.to}${exclude}`;
                    })
                    .filter(Boolean)
                    .join(', ')}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </ConfigSection>
  );
}
