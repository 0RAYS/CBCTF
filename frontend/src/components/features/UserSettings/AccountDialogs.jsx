import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import ConfirmModal from '../../common/ConfirmModal';
import { Input } from '../../common';

export default function AccountDialogs({ dialog, onClose, user, onDeleteAccount, onEmailVerify }) {
  const { t } = useTranslation();
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const closeDelete = () => {
    onClose();
    setPassword('');
    setError('');
  };

  return (
    <>
      <ConfirmModal
        isOpen={dialog === 'delete'}
        onClose={closeDelete}
        onConfirm={() => {
          if (!password) {
            setError(t('auth.validation.passwordRequired'));
            return;
          }
          onDeleteAccount(password);
          closeDelete();
        }}
        title={t('user.settings.deleteAccount')}
        message={
          <div className="space-y-3">
            <p>{t('user.settings.deleteAccountConfirm')}</p>
            <Input
              type="password"
              value={password}
              onChange={(event) => {
                setPassword(event.target.value);
                setError('');
              }}
              placeholder={t('auth.placeholders.password')}
              error={error}
              autoFocus
              autoComplete="current-password"
            />
          </div>
        }
        confirmText={t('common.delete')}
        type="danger"
      />
      <ConfirmModal
        isOpen={dialog === 'email'}
        onClose={onClose}
        onConfirm={() => {
          onEmailVerify(user.email);
          onClose();
        }}
        title={t('user.settings.verifyEmail')}
        message={t('user.settings.verifyEmailHint')}
        confirmText={t('common.send')}
      />
    </>
  );
}
