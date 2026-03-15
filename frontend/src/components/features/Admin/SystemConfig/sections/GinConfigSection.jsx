import { IconServer } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigField } from '../fields/ConfigField';
import { ConfigListField } from '../fields/ConfigListField';

const sanitizeNumber = (value, fallbackValue = 0) => {
  if (value === '' || value === null || value === undefined) {
    return fallbackValue;
  }
  const numeric = Number(value);
  return Number.isNaN(numeric) ? fallbackValue : numeric;
};

export function GinConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();

  return (
    <ConfigSection title={t('admin.system.sections.gin')} icon={IconServer}>
      <ConfigField
        label={t('admin.system.labels.ginHost')}
        value={config.gin.host}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gin.host = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.ginMode')}
        type="select"
        value={config.gin.mode}
        options={['debug', 'test', 'release']}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gin.mode = value;
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.ginPort')}
        type="number"
        value={config.gin.port}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gin.port = sanitizeNumber(value, config.gin.port);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.ginUploadMax')}
        type="number"
        value={config.gin.upload.max}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gin.upload.max = sanitizeNumber(value, config.gin.upload.max);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.ginGlobalRate')}
        type="number"
        value={config.gin.ratelimit.global}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gin.ratelimit.global = sanitizeNumber(value, config.gin.ratelimit.global);
          })
        }
      />
      <ConfigListField
        label={t('admin.system.labels.ginProxies')}
        items={config.gin.proxies || []}
        onAdd={() =>
          updateConfig((draft) => {
            draft.gin.proxies.push('');
          })
        }
        onUpdate={(index, value) =>
          updateConfig((draft) => {
            draft.gin.proxies[index] = value;
          })
        }
        onRemove={(index) =>
          updateConfig((draft) => {
            draft.gin.proxies.splice(index, 1);
          })
        }
      />
      <ConfigListField
        label={t('admin.system.labels.ginLogWhitelist')}
        items={config.gin.log.whitelist || []}
        onAdd={() =>
          updateConfig((draft) => {
            draft.gin.log.whitelist.push('');
          })
        }
        onUpdate={(index, value) =>
          updateConfig((draft) => {
            draft.gin.log.whitelist[index] = value;
          })
        }
        onRemove={(index) =>
          updateConfig((draft) => {
            draft.gin.log.whitelist.splice(index, 1);
          })
        }
      />
      <ConfigListField
        label={t('admin.system.labels.ginRateWhitelist')}
        items={config.gin.ratelimit.whitelist || []}
        onAdd={() =>
          updateConfig((draft) => {
            draft.gin.ratelimit.whitelist.push('');
          })
        }
        onUpdate={(index, value) =>
          updateConfig((draft) => {
            draft.gin.ratelimit.whitelist[index] = value;
          })
        }
        onRemove={(index) =>
          updateConfig((draft) => {
            draft.gin.ratelimit.whitelist.splice(index, 1);
          })
        }
      />
      <ConfigListField
        label={t('admin.system.labels.ginCORS')}
        items={config.gin.cors || []}
        onAdd={() =>
          updateConfig((draft) => {
            draft.gin.cors.push('');
          })
        }
        onUpdate={(index, value) =>
          updateConfig((draft) => {
            draft.gin.cors[index] = value;
          })
        }
        onRemove={(index) =>
          updateConfig((draft) => {
            draft.gin.cors.splice(index, 1);
          })
        }
      />
      <ConfigField
        label={t('admin.system.labels.ginJWTSecret')}
        type="password"
        value={config.gin.jwt.secret}
        placeholder={t('common.leaveBlankToKeep')}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.gin.jwt.secret = value;
          })
        }
      />
      <ConfigListField
        label={t('admin.system.labels.metricsWhitelist')}
        items={config.gin.metrics.whitelist || []}
        onAdd={() =>
          updateConfig((draft) => {
            draft.gin.metrics.whitelist.push('');
          })
        }
        onUpdate={(index, value) =>
          updateConfig((draft) => {
            draft.gin.metrics.whitelist[index] = value;
          })
        }
        onRemove={(index) =>
          updateConfig((draft) => {
            draft.gin.metrics.whitelist.splice(index, 1);
          })
        }
      />
    </ConfigSection>
  );
}
