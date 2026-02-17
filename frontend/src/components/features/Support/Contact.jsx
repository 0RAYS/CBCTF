import { motion } from 'motion/react';
import { Button, Card } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function Contact({ contactMethods = [] }) {
  const { t } = useTranslation();
  return (
    <div className="w-full max-w-[1000px] mx-auto space-y-8">
      <motion.h2
        className="text-2xl text-neutral-50 font-mono border-b border-neutral-300/30 pb-4"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        {t('contact.title')}
      </motion.h2>

      <motion.div
        className="grid grid-cols-2 gap-6"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        {contactMethods.map((method, index) => (
          <Card
            key={index}
            variant="default"
            padding="md"
            className="hover:border-geek-400 transition-all duration-200 hover:-translate-y-0.5"
          >
            <div className="text-4xl mb-4">{method.icon}</div>
            <h3 className="text-neutral-50 font-mono text-lg mb-2">{method.title}</h3>
            <p className="text-neutral-400 text-sm mb-4">{method.description}</p>

            <Button
              variant="primary"
              fullWidth
              className="shadow-[0_0_20px_rgba(89,126,247,0.4)]"
              onClick={method.onClick}
            >
              {method.action}
            </Button>
          </Card>
        ))}
      </motion.div>
    </div>
  );
}

export default Contact;
