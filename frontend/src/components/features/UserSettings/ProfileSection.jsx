import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import Avatar from '../../common/Avatar';
import { Button } from '../../common';

export default function ProfileSection({ user, onUpdate, onPictureChange, onVerify }) {
  const { t } = useTranslation();
  const handlePictureClick = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.onchange = (e) => {
      const file = e.target.files[0];
      if (file) onPictureChange(file);
    };
    input.click();
  };

  return (
    <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-6">
      <div className="flex flex-col items-start sm:flex-row sm:items-center gap-4 sm:gap-6">
        <button
          type="button"
          aria-label={t('common.change')}
          className="relative group shrink-0 rounded-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400"
          onClick={handlePictureClick}
        >
          <Avatar src={user.picture} name={user.name} size="xl" shape="circle" />
          <div className="absolute inset-0 flex items-center justify-center bg-black/60 rounded-full opacity-0 group-hover:opacity-100 group-focus-visible:opacity-100 transition-opacity cursor-pointer">
            <span className="text-neutral-200 text-sm font-mono">{t('common.change')}</span>
          </div>
        </button>
        <div className="min-w-0 [overflow-wrap:anywhere]">
          <h3 className="text-neutral-50 font-mono mb-1">{user.name}</h3>
          <p className="text-neutral-400 text-sm">
            {user.emailVerified ? t('user.settings.emailVerified') : t('user.settings.emailNotVerified')}
          </p>
        </div>
      </div>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          onUpdate({ name: e.target.name.value, email: e.target.email.value, description: e.target.description.value });
        }}
        className="space-y-4"
      >
        <div>
          <label className="block text-neutral-400 text-sm mb-2">{t('user.settings.displayName')}</label>
          <input
            required
            name="name"
            type="text"
            defaultValue={user.name}
            className="w-full p-3 bg-neutral-900 border border-neutral-300/30 rounded-md text-neutral-50 font-mono focus:outline-none focus:border-geek-400"
          />
        </div>
        <div>
          <label className="block text-neutral-400 text-sm mb-2">{t('auth.placeholders.email')}</label>
          <div className="flex flex-col sm:flex-row gap-2">
            <input
              required
              name="email"
              type="email"
              defaultValue={user.email}
              className="min-w-0 w-full flex-1 p-3 bg-neutral-900 border border-neutral-300/30 rounded-md text-neutral-50 font-mono focus:outline-none focus:border-geek-400"
            />
            {!user.emailVerified && onVerify && (
              <Button variant="outline" size="sm" className="shrink-0 w-full sm:w-auto" onClick={onVerify}>
                {t('user.settings.verify')}
              </Button>
            )}
          </div>
        </div>
        <div>
          <label className="block text-neutral-400 text-sm mb-2">{t('user.settings.bio')}</label>
          <textarea
            name="description"
            defaultValue={user.description}
            rows={4}
            className="w-full p-3 bg-neutral-900 border border-neutral-300/30 rounded-md text-neutral-50 focus:outline-none focus:border-geek-400"
          />
        </div>
        <div className="flex justify-end pt-4 border-t border-neutral-300/30">
          <Button type="submit" variant="primary" size="sm" className="w-full sm:w-auto min-h-10 h-auto! py-2">
            {t('common.saveChanges')}
          </Button>
        </div>
      </form>
    </motion.div>
  );
}
