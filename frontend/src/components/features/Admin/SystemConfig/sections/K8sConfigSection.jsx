import { IconCloud } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigField } from '../fields/ConfigField';
import { FrpServerField } from '../fields/FrpServerField';

export function K8sConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.k8s')} icon={IconCloud}>
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
        label={t('admin.system.labels.k8sCapture')}
        value={config.k8s.capture}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.k8s.capture = value;
          })
        }
      />
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
      <ConfigField
        label={t('admin.system.labels.frpc')}
        value={config.k8s.frp.frpc}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.k8s.frp.frpc = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.nginx')}
        value={config.k8s.frp.nginx}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.k8s.frp.nginx = value;
          })
        }
      />
      <FrpServerField config={config} updateConfig={updateConfig} />
    </ConfigSection>
  );
}
