import { motion } from 'motion/react';
import { useState } from 'react';
import { Button } from '../../common';
import { useTranslation } from 'react-i18next';

function AdminAuthPanel({ onSubmit }) {
  const { t } = useTranslation();
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  });
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const validateForm = () => {
    const newErrors = {};

    // Username validation
    if (!formData.username) {
      newErrors.username = t('auth.validation.usernameRequired');
    } else if (formData.username.length < 3) {
      newErrors.username = t('auth.validation.usernameMin');
    }

    // Password validation
    if (!formData.password) {
      newErrors.password = t('auth.validation.passwordRequired');
    } else if (formData.password.length < 6) {
      newErrors.password = t('auth.validation.passwordMin');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!validateForm()) return;

    setIsSubmitting(true);
    try {
      // Call external handler function
      await onSubmit({
        type: 'login',
        data: {
          username: formData.username,
          password: formData.password,
        },
      });
    } catch (error) {
      // Handle error, possibly display error message
      setErrors((prev) => ({
        ...prev,
        submit: error.message,
      }));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    // Clear corresponding field error
    if (errors[name]) {
      setErrors((prev) => ({
        ...prev,
        [name]: '',
      }));
    }
  };

  return (
    <motion.form
      onSubmit={handleSubmit}
      className="w-[400px] bg-black/40 backdrop-blur-[2px] border border-neutral-300 rounded-md p-8"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2 }}
    >
      {/* Title area */}
      <div className="relative flex justify-center mb-8">
        <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
        <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
        <motion.span
          className="text-neutral-300 text-2xl font-mono tracking-wider"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.15 }}
        >
          {t('auth.adminLogin')}
        </motion.span>
      </div>

      {/* Form area */}
      <motion.div className="space-y-4">
        <div>
          <input
            type="text"
            name="username"
            required
            value={formData.username}
            onChange={handleChange}
            placeholder={t('auth.placeholders.username')}
            className={`
                w-full h-[40px] bg-black/20 border rounded-md px-4 
                text-neutral-50 placeholder-neutral-400
                transition-all duration-200
                ${
                  errors.username
                    ? 'border-red-500 focus:border-red-500 focus:shadow-[0_0_15px_rgba(239,68,68,0.3)]'
                    : 'border-neutral-300 focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]'
                }
            `}
          />
          {errors.username && (
            <motion.span
              className="text-red-500 text-xs mt-1 block"
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.15 }}
            >
              {errors.username}
            </motion.span>
          )}
        </div>

        <div>
          <input
            type="password"
            name="password"
            required
            value={formData.password}
            onChange={handleChange}
            placeholder={t('auth.placeholders.password')}
            className={`
                w-full h-[40px] bg-black/20 border rounded-md px-4 
                text-neutral-50 placeholder-neutral-400
                transition-all duration-200
                ${
                  errors.password
                    ? 'border-red-500 focus:border-red-500 focus:shadow-[0_0_15px_rgba(239,68,68,0.3)]'
                    : 'border-neutral-300 focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]'
                }
            `}
          />
          {errors.password && (
            <motion.span
              className="text-red-500 text-xs mt-1 block"
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.15 }}
            >
              {errors.password}
            </motion.span>
          )}
        </div>

        {errors.submit && (
          <motion.div className="text-red-500 text-sm text-center" initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
            {errors.submit}
          </motion.div>
        )}

        <Button
          type="submit"
          variant="primary"
          fullWidth
          className="shadow-[0_0_20px_rgba(89,126,247,0.4)]"
          disabled={isSubmitting}
        >
          {isSubmitting ? t('common.processing') : t('auth.login')}
        </Button>
      </motion.div>
    </motion.form>
  );
}

export default AdminAuthPanel;
