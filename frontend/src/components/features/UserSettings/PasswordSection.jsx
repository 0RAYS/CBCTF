import { motion } from 'motion/react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../common';

const emptyPasswords = { currentPassword: '', newPassword: '', confirmPassword: '' };

export default function PasswordSection({ user, onPasswordChange, onDelete }) {
  const { t } = useTranslation();
  const [values, setValues] = useState(emptyPasswords);
  const [errors, setErrors] = useState(emptyPasswords);
  const handleBlur = ({ target: { name, value } }) => {
    if (name === 'newPassword') {
      setErrors((prev) => ({
        ...prev,
        newPassword: value && value.length < 6 ? t('auth.validation.passwordMin') : '',
      }));
    }
    if (name === 'confirmPassword' && values.newPassword) {
      setErrors((prev) => ({
        ...prev,
        confirmPassword: value !== values.newPassword ? t('auth.validation.passwordMismatch') : '',
      }));
    }
  };
  const handleSubmit = (e) => {
    e.preventDefault();
    const nextErrors = {
      currentPassword: user.hasNoPwd
        ? ''
        : !values.currentPassword
          ? t('user.settings.validation.currentPasswordRequired')
          : '',
      newPassword: values.newPassword.length < 6 ? t('auth.validation.passwordMin') : '',
      confirmPassword: values.newPassword !== values.confirmPassword ? t('auth.validation.passwordMismatch') : '',
    };
    setErrors(nextErrors);
    if (!Object.values(nextErrors).some(Boolean)) {
      onPasswordChange(values);
      setValues(emptyPasswords);
      setErrors(emptyPasswords);
    }
  };

  return (
    <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-6">
      <form onSubmit={handleSubmit} className="p-4 border border-neutral-300/30 rounded-md">
        <h3 className="text-neutral-50 font-mono mb-4">{t('user.settings.changePassword')}</h3>
        <div className="space-y-4">
          {['currentPassword', 'newPassword', 'confirmPassword']
            .filter((name) => name !== 'currentPassword' || !user.hasNoPwd)
            .map((name) => (
              <div key={name}>
                <input
                  name={name}
                  type="password"
                  placeholder={t(`user.settings.placeholders.${name}`)}
                  value={values[name]}
                  onChange={(e) => setValues((prev) => ({ ...prev, [name]: e.target.value }))}
                  onBlur={handleBlur}
                  className={`w-full p-3 bg-neutral-900 border rounded-md text-neutral-50 font-mono focus:outline-none transition-all duration-200 ${errors[name] ? 'border-red-400 focus:border-red-400 focus:shadow-error' : 'border-neutral-300/30 focus:border-geek-400 focus:shadow-focus'}`}
                  required
                />
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: errors[name] ? 1 : 0, height: errors[name] ? 'auto' : 0 }}
                  transition={{ duration: 0.2 }}
                >
                  {errors[name] && <p className="mt-1 text-red-400 text-sm">{errors[name]}</p>}
                </motion.div>
              </div>
            ))}
          <div className="flex justify-end">
            <Button type="submit" variant="primary" size="sm" className="w-full sm:w-auto min-h-10 h-auto! py-2">
              {t('user.settings.updatePassword')}
            </Button>
          </div>
        </div>
      </form>
      {onDelete && (
        <div className="p-4 border border-red-400/20 rounded-md">
          <h3 className="text-red-400 font-mono mb-2">{t('user.settings.deleteAccount')}</h3>
          <p className="text-neutral-400 text-sm mb-4">{t('user.settings.deleteAccountHint')}</p>
          <Button
            variant="ghost"
            textColor="text-red-400"
            size="sm"
            className="w-full sm:w-auto min-h-10 h-auto! py-2 hover:bg-red-400/10!"
            onClick={onDelete}
          >
            {t('user.settings.deleteAccountAction')}
          </Button>
        </div>
      )}
    </motion.div>
  );
}
