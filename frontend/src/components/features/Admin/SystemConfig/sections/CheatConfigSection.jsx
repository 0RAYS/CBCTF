import { IconServer } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigListField } from '../fields/ConfigListField';

export function CheatConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.cheat')} icon={IconServer}>
      <ConfigListField
        label={t('admin.system.labels.cheatIpWhitelist')}
        items={config.cheat.ip_whitelist || []}
        onAdd={() =>
          updateConfig((draft) => {
            draft.cheat.ip_whitelist.push('');
          })
        }
        onUpdate={(index, value) =>
          updateConfig((draft) => {
            draft.cheat.ip_whitelist[index] = value;
          })
        }
        onRemove={(index) =>
          updateConfig((draft) => {
            draft.cheat.ip_whitelist.splice(index, 1);
          })
        }
      />
    </ConfigSection>
  );
}
