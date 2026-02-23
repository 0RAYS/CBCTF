import { useEffect, useState } from 'react';
import { IconUserPlus } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { ConfigSection } from '../fields/ConfigSection';
import { ConfigField } from '../fields/ConfigField';
import { getGroupList } from '../../../../../api/admin/rbac';
import { Select } from '../../../../common';

export function RegistrationConfigSection({ config, updateConfig }) {
  const { t } = useTranslation();
  const [groups, setGroups] = useState([]);

  useEffect(() => {
    getGroupList({ limit: 100, offset: 0 }).then((res) => {
      if (res.code === 200) setGroups(res.data.groups);
    });
  }, []);

  const groupOptions = [
    { value: '0', label: t('admin.system.labels.registrationDefaultGroupNone') },
    ...groups.map((g) => ({ value: String(g.id), label: g.name })),
  ];

  return (
    <ConfigSection title={t('admin.system.sections.registration')} icon={IconUserPlus}>
      <ConfigField
        label={t('admin.system.labels.registrationEnabled')}
        type="boolean"
        value={config.registration.enabled}
        options={[
          { value: 'true', label: t('common.yes') },
          { value: 'false', label: t('common.no') },
        ]}
        onChange={(value) =>
          updateConfig((draft) => {
            draft.registration.enabled = value;
          })
        }
      />
      <div className="space-y-1">
        <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.registrationDefaultGroup')}</span>
        <Select
          size="sm"
          value={String(config.registration.default_group ?? 0)}
          onChange={(e) =>
            updateConfig((draft) => {
              draft.registration.default_group = Number(e.target.value);
            })
          }
          options={groupOptions}
        />
      </div>
    </ConfigSection>
  );
}
