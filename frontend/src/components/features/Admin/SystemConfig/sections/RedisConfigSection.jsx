import { IconServer } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '@/components/features/Admin/SystemConfig/fields/ConfigSection';
import { ConfigField } from '@/components/features/Admin/SystemConfig/fields/ConfigField';

export function RedisConfigSection({ config }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.redis')} icon={IconServer}>
      <ConfigField label={t('admin.system.labels.redisHost')} value={config.redis.host} disabled />
      <ConfigField label={t('admin.system.labels.redisPort')} type="number" value={config.redis.port} disabled />
      <ConfigField label={t('admin.system.labels.redisPassword')} type="password" value={config.redis.pwd} disabled />
    </ConfigSection>
  );
}
