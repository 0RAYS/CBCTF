import Contact from '../components/features/Support/Contact';
import { useTranslation } from 'react-i18next';

function ContactPage() {
  const { t } = useTranslation();
  const contactConfig = {
    contactMethods: [
      {
        type: 'github',
        title: t('contact.methods.github.title'),
        description: t('contact.methods.github.description'),
        icon: '💬',
        value: '0RAYS',
        action: t('contact.methods.github.action'),
        onClick: () => window.open(`https://github.com/0RAYS`, '_blank'),
      },
      {
        type: 'email',
        title: t('contact.methods.email.title'),
        description: t('contact.methods.email.description'),
        icon: '📧',
        value: 'admin@0rays.club',
        action: t('contact.methods.email.action'),
        onClick: () => (window.location.href = 'mailto:admin@0rays.club'),
      },
    ],
  };

  return <Contact {...contactConfig} />;
}

export default ContactPage;
