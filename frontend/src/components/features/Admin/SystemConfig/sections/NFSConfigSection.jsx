import { IconFolder } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigField } from '../fields/ConfigField';

export function NFSConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.nfs')} icon={IconFolder}>
      <ConfigField
        label={t('admin.system.labels.nfsServer')}
        value={config.nfs.server}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.nfs.server = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.nfsPath')}
        value={config.nfs.path}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.nfs.path = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.nfsStorage')}
        value={config.nfs.storage}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.nfs.storage = value;
          })
        }
      />
    </ConfigSection>
  );
}
