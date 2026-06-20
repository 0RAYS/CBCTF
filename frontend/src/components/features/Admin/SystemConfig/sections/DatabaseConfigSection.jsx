import { IconDatabase } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigField } from '../fields/ConfigField';

export function DatabaseConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();
  const postgres = config.gorm.postgres;

  return (
    <ConfigSection title={t('admin.system.sections.database')} icon={IconDatabase}>
      <ConfigField label={t('admin.system.labels.dbHost')} value={postgres.host} disabled />
      <ConfigField label={t('admin.system.labels.dbPort')} type="number" value={postgres.port} disabled />
      <ConfigField label={t('admin.system.labels.dbName')} value={postgres.db} disabled />
      <ConfigField label={t('admin.system.labels.dbUser')} value={postgres.user} disabled />
      <ConfigField label={t('admin.system.labels.dbPassword')} type="password" value={postgres.pwd} disabled />
      <ConfigField label={t('admin.system.labels.dbSSLMode')} value={postgres.sslmode ? 'true' : 'false'} disabled />
      <ConfigField label={t('admin.system.labels.dbOpen')} type="number" value={postgres.mxopen} disabled />
      <ConfigField label={t('admin.system.labels.dbIdle')} type="number" value={postgres.mxidle} disabled />
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
    </ConfigSection>
  );
}
