import { useRef } from 'react';
import { IconServer, IconUpload } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../../common';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigField } from '../fields/ConfigField';

const sanitizeNumber = (value, fallbackValue = 0) => {
  if (value === '' || value === null || value === undefined) {
    return fallbackValue;
  }
  const numeric = Number(value);
  return Number.isNaN(numeric) ? fallbackValue : numeric;
};

export function BasicConfigSection({ config, updateConfig, onUploadGeoCityDB, isUploadingGeoCityDB }) {
  const { t } = useTranslation();
  const fileInputRef = useRef(null);

  const handleGeoCityDBFileChange = async (event) => {
    const file = event.target.files?.[0];
    if (file && onUploadGeoCityDB) {
      await onUploadGeoCityDB(file);
    }
    event.target.value = '';
  };

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
      <div className="space-y-2">
        <ConfigField label={t('admin.system.labels.geocity_db')} value={config.geocity_db} />
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
          <input ref={fileInputRef} type="file" className="hidden" onChange={handleGeoCityDBFileChange} />
          <Button
            variant="outline"
            size="sm"
            icon={<IconUpload size={16} />}
            onClick={() => fileInputRef.current?.click()}
            loading={isUploadingGeoCityDB}
          >
            {t('admin.system.actions.uploadGeoCityDB')}
          </Button>
          <span className="text-xs text-neutral-500">{t('admin.system.hints.uploadGeoCityDB')}</span>
        </div>
      </div>
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
        label={t('admin.system.labels.trafficConcurrency')}
        type="number"
        value={config.asyncq.queues.traffic}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.asyncq.queues.traffic = sanitizeNumber(value, config.asyncq.queues.traffic);
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
