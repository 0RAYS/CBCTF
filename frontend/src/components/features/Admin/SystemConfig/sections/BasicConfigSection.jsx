import { IconServer } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigField } from '../fields/ConfigField';

const sanitizeNumber = (value, fallbackValue = 0) => {
  if (value === '' || value === null || value === undefined) {
    return fallbackValue;
  }
  const numeric = Number(value);
  return Number.isNaN(numeric) ? fallbackValue : numeric;
};

export function BasicConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.basic')} icon={IconServer}>
      <ConfigField
        label={t('admin.system.labels.host')}
        value={config.host}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.host = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.path')}
        value={config.path}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.path = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.geocity_db')}
        value={config.geocity_db}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.geocity_db = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.logLevel')}
        type="select"
        value={config.log.level}
        options={['DEBUG', 'INFO', 'WARNING', 'ERROR']}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.log.level = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.saveLogs')}
        type="boolean"
        value={config.log.save}
        options={[
          { value: 'true', label: t('common.yes') },
          { value: 'false', label: t('common.no') },
        ]}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.log.save = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.workerLog')}
        type="select"
        value={config.asyncq.log.level}
        options={['DEBUG', 'INFO', 'WARNING', 'ERROR']}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.log.level = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.concurrency')}
        type="number"
        value={config.asyncq.concurrency}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.concurrency = sanitizeNumber(value, config.asyncq.concurrency);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.victimConcurrency')}
        type="number"
        value={config.asyncq.queues.victim}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.queues.victim = sanitizeNumber(value, config.asyncq.queues.victim);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.generatorConcurrency')}
        type="number"
        value={config.asyncq.queues.generator}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.queues.generator = sanitizeNumber(value, config.asyncq.queues.generator);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.attachmentConcurrency')}
        type="number"
        value={config.asyncq.queues.attachment}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.queues.attachment = sanitizeNumber(value, config.asyncq.queues.attachment);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.emailConcurrency')}
        type="number"
        value={config.asyncq.queues.email}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.queues.email = sanitizeNumber(value, config.asyncq.queues.email);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.webhookConcurrency')}
        type="number"
        value={config.asyncq.queues.webhook}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.queues.webhook = sanitizeNumber(value, config.asyncq.queues.webhook);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.imageConcurrency')}
        type="number"
        value={config.asyncq.queues.image}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.queues.image = sanitizeNumber(value, config.asyncq.queues.image);
          })
        }
      />
    </ConfigSection>
  );
}
