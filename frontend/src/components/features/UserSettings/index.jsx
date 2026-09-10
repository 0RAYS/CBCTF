import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Card } from '../../common';
import ProfileSection from './ProfileSection';
import PasswordSection from './PasswordSection';
import AccountDialogs from './AccountDialogs';

export default function UserSettings({
  user,
  onUpdate,
  onPasswordChange,
  onEmailVerify,
  onDeleteAccount,
  onPictureChange,
  onLogout,
}) {
  const { t } = useTranslation();
  const [activeSection, setActiveSection] = useState('profile');
  const [dialog, setDialog] = useState(null);
  const sections = [
    { id: 'profile', label: t('user.settings.sections.profile'), icon: '👤' },
    { id: 'security', label: t('user.settings.sections.security'), icon: '🔒' },
    { id: 'divider', type: 'divider' },
    { id: 'logout', label: t('user.settings.sections.logout'), icon: '🖖' },
  ];

  return (
    <div className="w-full max-w-[1200px] mx-auto">
      <Card variant="default" padding="none" animate className="overflow-hidden">
        <div className="p-4 sm:p-6 border-b border-neutral-300/30">
          <h2 className="text-2xl font-mono text-neutral-50">{t('user.settings.title')}</h2>
        </div>
        <div className="flex flex-col md:flex-row">
          <div className="w-full md:w-[220px] shrink-0 border-b md:border-b-0 md:border-r border-neutral-300/30 p-4">
            <div className="space-y-2">
              {sections.map((section) =>
                section.type === 'divider' ? (
                  <div key={section.id} className="my-4 border-t border-neutral-300/30" />
                ) : (
                  <Button
                    key={section.id}
                    variant="ghost"
                    align="icon-left-text-center"
                    icon={<span className="text-xl">{section.icon}</span>}
                    textColor={
                      section.id === 'logout'
                        ? 'text-neutral-400'
                        : activeSection === section.id
                          ? 'text-geek-400'
                          : 'text-neutral-300'
                    }
                    className={`w-full px-4 py-3 rounded-md min-w-0 h-auto ${section.id === 'logout' ? 'hover:bg-white/5' : activeSection === section.id ? 'bg-geek-400/10' : 'hover:bg-neutral-700/10 hover:text-neutral-200'}`}
                    onClick={() => (section.id === 'logout' ? onLogout() : setActiveSection(section.id))}
                  >
                    <span className="font-mono text-sm">{section.label}</span>
                  </Button>
                )
              )}
            </div>
          </div>
          <div className="min-w-0 flex-1 p-4 sm:p-6">
            {/* Keep both forms mounted so switching tabs does not discard drafts or errors. */}
            <div hidden={activeSection !== 'profile'}>
              <ProfileSection
                user={user}
                onUpdate={onUpdate}
                onPictureChange={onPictureChange}
                onVerify={onEmailVerify ? () => setDialog('email') : undefined}
              />
            </div>
            <div hidden={activeSection !== 'security'}>
              <PasswordSection
                user={user}
                onPasswordChange={onPasswordChange}
                onDelete={onDeleteAccount ? () => setDialog('delete') : undefined}
              />
            </div>
          </div>
        </div>
      </Card>
      <AccountDialogs
        dialog={dialog}
        onClose={() => setDialog(null)}
        user={user}
        onDeleteAccount={onDeleteAccount}
        onEmailVerify={onEmailVerify}
      />
    </div>
  );
}
