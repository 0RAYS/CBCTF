import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import { FormField, Input, Textarea } from '../../../common';
import Checkbox from '../../../common/Checkbox';

export default function UserForm({ mode, selectedUser, editForm, setEditForm }) {
  const { t } = useTranslation();
  const descriptionId = useId();
  if (mode === 'delete') {
    return (
      <div className="space-y-2 text-neutral-300">
        <p>
          {t('admin.users.modal.deletePrompt')} <span className="text-white font-semibold">{selectedUser?.name}</span>?
        </p>
        <p className="text-red-400 text-sm">{t('admin.users.modal.deleteWarning')}</p>
      </div>
    );
  }
  return (
    <div className="space-y-4">
      <Input
        label={t('admin.users.form.username')}
        value={editForm.name}
        onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
        placeholder={t('admin.users.form.usernamePlaceholder')}
        required={mode === 'create'}
      />
      <Input
        label={t('admin.users.form.email')}
        type="email"
        value={editForm.email}
        onChange={(e) => setEditForm({ ...editForm, email: e.target.value })}
        placeholder={t('admin.users.form.emailPlaceholder')}
        required={mode === 'create'}
      />
      <div>
        <Input
          label={t('admin.users.form.password')}
          type="password"
          value={editForm.password}
          onChange={(e) => setEditForm({ ...editForm, password: e.target.value })}
          placeholder={t(`admin.users.form.${mode === 'edit' ? 'passwordEditPlaceholder' : 'passwordPlaceholder'}`)}
          required={mode === 'create'}
        />
        {mode === 'edit' && <p className="mt-1 text-xs text-neutral-500">{t('admin.users.form.passwordEditHint')}</p>}
      </div>
      <FormField label={t('admin.users.form.description')} htmlFor={descriptionId}>
        <Textarea
          id={descriptionId}
          value={editForm.description}
          onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
          placeholder={t('admin.users.form.descriptionPlaceholder')}
          rows={3}
          fullWidth
        />
      </FormField>
      <div className="flex flex-col gap-2">
        {['verified', 'banned', 'hidden'].map((field) => (
          <Checkbox
            key={field}
            label={t(`admin.users.status.${field}`)}
            checked={editForm[field]}
            onChange={(e) => setEditForm({ ...editForm, [field]: e.target.checked })}
          />
        ))}
      </div>
    </div>
  );
}
