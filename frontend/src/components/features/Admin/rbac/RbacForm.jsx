import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import { FormField, Input, Select, Textarea } from '../../../common';

export default function RbacForm({ domain, mode, item, form, setForm, roles }) {
  const { t } = useTranslation();
  const id = useId();
  const prefix = `admin.rbac.${domain}`;
  if (mode === 'delete') {
    return (
      <p className="text-neutral-300">
        {t(`${prefix}.modal.deletePrompt`)} <span className="text-white font-semibold">{item?.name}</span>?{' '}
        {t(`${prefix}.modal.deleteWarning`)}
      </p>
    );
  }
  return (
    <div className="space-y-4">
      <Input
        label={t(`${prefix}.form.name`)}
        value={form.name}
        onChange={(e) => setForm({ ...form, name: e.target.value })}
        placeholder={t(`${prefix}.form.namePlaceholder`)}
        required={mode === 'create'}
        disabled={mode === 'edit' && item?.default}
      />
      <FormField label={t(`${prefix}.form.description`)} htmlFor={`${id}-description`}>
        <Textarea
          id={`${id}-description`}
          value={form.description}
          onChange={(e) => setForm({ ...form, description: e.target.value })}
          placeholder={t(`${prefix}.form.descriptionPlaceholder`)}
          rows={3}
          fullWidth
        />
      </FormField>
      {roles && (
        <FormField label={t(`${prefix}.form.role`)} htmlFor={`${id}-role`}>
          <Select
            id={`${id}-role`}
            value={form.role_id}
            onChange={(e) => setForm({ ...form, role_id: parseInt(e.target.value) || '' })}
            options={roles.map((role) => ({ value: role.id, label: role.name }))}
            placeholder={t(`${prefix}.form.rolePlaceholder`)}
            fullWidth
          />
        </FormField>
      )}
    </div>
  );
}
