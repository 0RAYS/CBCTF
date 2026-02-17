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
        value={config.gorm.mysql.host}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.mysql.host = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbPort')}
        type="number"
        value={config.gorm.mysql.port}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.mysql.port = sanitizeNumber(value, config.gorm.mysql.port);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbName')}
        value={config.gorm.mysql.db}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.mysql.db = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbUser')}
        value={config.gorm.mysql.user}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.mysql.user = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbPassword')}
        type="password"
        value={config.gorm.mysql.pwd}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.mysql.pwd = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbIdle')}
        type="number"
        value={config.gorm.mysql.mxidle}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.mysql.mxidle = sanitizeNumber(value, config.gorm.mysql.mxidle);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.dbOpen')}
        type="number"
        value={config.gorm.mysql.mxopen}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gorm.mysql.mxopen = sanitizeNumber(value, config.gorm.mysql.mxopen);
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
