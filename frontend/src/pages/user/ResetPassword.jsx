import { useState, useEffect } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { Button, Input } from '../../components/common';
import { resetPassword } from '../../api/auth';

function ResetPassword() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();

  const [token, setToken] = useState('');
  const [id, setId] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [succeeded, setSucceeded] = useState(false);
  const [tokenInvalid, setTokenInvalid] = useState(false);

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const t = params.get('token');
    const i = params.get('id');
    if (!t || !i) {
      setTokenInvalid(true);
    } else {
      setToken(t);
      setId(i);
    }
  }, [location.search]);

  const validate = () => {
    const newErrors = {};
    if (password.length < 6) {
      newErrors.password = t('auth.validation.passwordMin');
    }
    if (confirmPassword !== password) {
      newErrors.confirmPassword = t('auth.validation.passwordMismatch');
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!validate()) return;
    setIsSubmitting(true);
    try {
      const response = await resetPassword({ token, id, password });
      if (response.code === 200) {
        setSucceeded(true);
      } else {
        setTokenInvalid(true);
      }
    } catch {
      setTokenInvalid(true);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-190px)] px-4">
      <motion.div
        className="w-full max-w-[400px] bg-neutral-800/80 border border-neutral-600/60 rounded-md p-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3, ease: [0.25, 1, 0.5, 1] }}
      >
        <div className="relative flex justify-center mb-8">
          <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
          <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
          <h1 className="text-neutral-50 text-2xl font-mono tracking-wider">{t('auth.resetPassword.title')}</h1>
        </div>

        {tokenInvalid ? (
          <div className="space-y-4 text-center">
            <p className="text-red-400 text-sm font-mono leading-relaxed">{t('auth.resetPassword.invalidToken')}</p>
            <Link to="/login">
              <Button variant="outline" fullWidth>
                {t('auth.resetPassword.requestNew')}
              </Button>
            </Link>
          </div>
        ) : succeeded ? (
          <div className="space-y-4 text-center">
            <p className="text-neutral-300 text-sm leading-relaxed">{t('auth.resetPassword.successMessage')}</p>
            <Button variant="primary" fullWidth onClick={() => navigate('/login')}>
              {t('auth.resetPassword.goToLogin')}
            </Button>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <p className="text-neutral-400 text-sm leading-relaxed">{t('auth.resetPassword.description')}</p>
            <Input
              type="password"
              name="password"
              required
              value={password}
              onChange={(e) => {
                setPassword(e.target.value);
                if (errors.password) setErrors((prev) => ({ ...prev, password: '' }));
              }}
              label={t('auth.placeholders.newPassword')}
              placeholder={t('auth.placeholders.newPassword')}
              error={errors.password}
            />
            <Input
              type="password"
              name="confirmPassword"
              required
              value={confirmPassword}
              onChange={(e) => {
                setConfirmPassword(e.target.value);
                if (errors.confirmPassword) setErrors((prev) => ({ ...prev, confirmPassword: '' }));
              }}
              label={t('auth.placeholders.confirmNewPassword')}
              placeholder={t('auth.placeholders.confirmNewPassword')}
              error={errors.confirmPassword}
            />
            <Button type="submit" variant="primary" fullWidth disabled={isSubmitting}>
              {isSubmitting ? t('auth.resetPassword.submitting') : t('auth.resetPassword.submit')}
            </Button>
          </form>
        )}
      </motion.div>
    </div>
  );
}

export default ResetPassword;
