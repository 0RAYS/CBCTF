import { IconPlus, IconX } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, Input } from '../../../../common';
import { ConfigField } from './ConfigField';

export function ExternalNetworkField({ externalNetworks, updateConfig }) {
  const { t } = useTranslation();
  const networks = externalNetworks || { enabled: false, interfaces: [] };
  const interfaces = networks.interfaces || [];

  const addInterface = () => {
    updateConfig((draft) => {
      if (!draft.k8s.external_networks) {
        draft.k8s.external_networks = { enabled: false, interfaces: [] };
      }
      draft.k8s.external_networks.interfaces.push({ interface: '', cidr: '', gateway: '' });
    });
  };

  const removeInterface = (index) => {
    updateConfig((draft) => {
      draft.k8s.external_networks.interfaces.splice(index, 1);
    });
  };

  const updateInterface = (index, field, value) => {
    updateConfig((draft) => {
      draft.k8s.external_networks.interfaces[index][field] = value;
    });
  };

  return (
    <div className="space-y-2">
      <ConfigField
        label={t('admin.system.labels.externalNetworkEnabled')}
        type="boolean"
        value={networks.enabled}
        options={[
          { value: 'true', label: t('common.yes') },
          { value: 'false', label: t('common.no') },
        ]}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.k8s.external_networks.enabled = value;
          })
        }
      />

      <div className="flex items-center justify-between">
        <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.externalNetworkInterfaces')}</span>
        <Button size="icon" variant="ghost" aria-label={t('common.add')} onClick={addInterface}>
          <IconPlus size={14} />
        </Button>
      </div>

      <div className="space-y-2">
        {interfaces.map((item, index) => (
          <div key={`external-network-${index}`} className="space-y-2 border border-neutral-300/10 rounded p-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-mono text-neutral-400">
                {t('admin.system.labels.externalNetwork')} #{index + 1}
              </span>
              <Button
                size="icon"
                variant="ghost"
                aria-label={t('common.remove')}
                onClick={() => removeInterface(index)}
              >
                <IconX size={14} />
              </Button>
            </div>
            <div className="grid gap-2 md:grid-cols-3">
              <div className="space-y-1">
                <span className="text-xs font-mono text-neutral-400">
                  {t('admin.system.labels.externalNetworkInterface')}
                </span>
                <Input
                  size="sm"
                  value={item.interface || ''}
                  placeholder="ens192"
                  onChange={(event) => updateInterface(index, 'interface', event.target.value)}
                />
              </div>
              <div className="space-y-1">
                <span className="text-xs font-mono text-neutral-400">
                  {t('admin.system.labels.externalNetworkCidr')}
                </span>
                <Input
                  size="sm"
                  value={item.cidr || ''}
                  placeholder="192.168.0.0/24"
                  onChange={(event) => updateInterface(index, 'cidr', event.target.value)}
                />
              </div>
              <div className="space-y-1">
                <span className="text-xs font-mono text-neutral-400">
                  {t('admin.system.labels.externalNetworkGateway')}
                </span>
                <Input
                  size="sm"
                  value={item.gateway || ''}
                  placeholder="192.168.0.1"
                  onChange={(event) => updateInterface(index, 'gateway', event.target.value)}
                />
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
