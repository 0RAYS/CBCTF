import { useRef } from 'react';
import { IconDatabase, IconUpload } from '@tabler/icons-react';
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
    <ConfigSection title={t('admin.system.sections.basic')} icon={IconDatabase}>
      <ConfigField
        label={t('admin.system.labels.host')}
        value={config.host}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.host = value;
          })
        }
      />
      <ConfigField label={t('admin.system.labels.path')} value={config.path} disabled />
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
        label={t('admin.system.labels.dbLogLevel')}
        type="select"
        value={config.gorm.log.level}
        options={['SILENT', 'INFO', 'WARNING', 'ERROR']}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.log.level = value;
          })
        }
      />
      {[
        ['victim', 'victimConcurrency'],
        ['traffic', 'trafficConcurrency'],
        ['generator', 'generatorConcurrency'],
        ['attachment', 'attachmentConcurrency'],
        ['email', 'emailConcurrency'],
        ['webhook', 'webhookConcurrency'],
        ['image', 'imageConcurrency'],
      ].map(([queue, label]) => (
        <ConfigField
          key={queue}
          label={t(`admin.system.labels.${label}`)}
          type="number"
          value={config.asyncq.queues[queue]}
          onChange={(value) =>
            updateConfig((draft) => {
              draft.asyncq.queues[queue] = sanitizeNumber(value, config.asyncq.queues[queue]);
            })
          }
        />
      ))}
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
    </ConfigSection>
  );
}
