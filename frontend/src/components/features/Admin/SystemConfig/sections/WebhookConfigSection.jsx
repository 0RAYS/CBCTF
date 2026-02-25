import { IconServer } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigListField } from '../fields/ConfigListField.jsx';

export function WebhookConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.webhook')} icon={IconServer}>
      <ConfigListField
        label={t('admin.system.labels.webhookWhitelist')}
        items={config.webhook.whitelist || []}
        onAdd={() =>
          updateConfig((draft) => {
            draft.webhook.whitelist.push('');
          })
        }
        onUpdate={(index, value) =>
          updateConfig((draft) => {
            draft.webhook.whitelist[index] = value;
          })
        }
        onRemove={(index) =>
          updateConfig((draft) => {
            draft.webhook.whitelist.splice(index, 1);
          })
        }
      />
    </ConfigSection>
  );
}
