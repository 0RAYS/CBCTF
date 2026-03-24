import { IconDatabase } from '@tabler/icons-react';
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

export function DatabaseConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.database')} icon={IconDatabase}>
      <ConfigField
        label={t('admin.system.labels.dbHost')}
        value={config.gorm.postgres.host}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.postgres.host = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbPort')}
        type="number"
        value={config.gorm.postgres.port}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.postgres.port = sanitizeNumber(value, config.gorm.postgres.port);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbName')}
        value={config.gorm.postgres.db}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.postgres.db = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbUser')}
        value={config.gorm.postgres.user}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.postgres.user = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbPassword')}
        type="password"
        value={config.gorm.postgres.pwd}
        placeholder={t('common.leaveBlankToKeep')}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.postgres.pwd = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbIdle')}
        type="number"
        value={config.gorm.postgres.mxidle}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.postgres.mxidle = sanitizeNumber(value, config.gorm.postgres.mxidle);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbOpen')}
        type="number"
        value={config.gorm.postgres.mxopen}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.postgres.mxopen = sanitizeNumber(value, config.gorm.postgres.mxopen);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbLogLevel')}
        type="select"
        value={config.gorm.log.level}
        options={['INFO', 'WARNING', 'ERROR', 'SILENT']}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.log.level = value;
          })
        }
      />
    </ConfigSection>
  );
}
