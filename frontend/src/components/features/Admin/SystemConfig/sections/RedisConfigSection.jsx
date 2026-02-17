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

export function RedisConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.redis')} icon={IconDatabase}>
      <ConfigField
        label={t('admin.system.labels.redisHost')}
        value={config.redis.host}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.redis.host = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.redisPort')}
        type="number"
        value={config.redis.port}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.redis.port = sanitizeNumber(value, config.redis.port);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.redisPassword')}
        type="password"
        value={config.redis.pwd}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.redis.pwd = value;
          })
        }
      />
    </ConfigSection>
  );
}
