import { useState } from 'react';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { Modal } from '../../../common';
import { useSystemConfig } from '../../../../hooks/useSystemConfig';
import { SystemConfigForm } from './SystemConfigForm';

/**
 * SystemConfig - Main container component for system configuration
 * @param {Object} config - Initial config from parent
 */
function SystemConfig({ config }) {
  const { t } = useTranslation();
  const {
    config: editableConfig,
    updateConfig,
    isUpdating,
    isRestarting,
    handleUpdateConfig,
    handleRestartSystem,
  } = useSystemConfig(config, t);

  const [isUpdateConfirmOpen, setIsUpdateConfirmOpen] = useState(false);
  const [isRestartConfirmOpen, setIsRestartConfirmOpen] = useState(false);

  const handleUpdateClick = () => {
    setIsUpdateConfirmOpen(true);
  };

  const handleRestartClick = () => {
    setIsRestartConfirmOpen(true);
  };

  const handleConfirmUpdate = async () => {
    setIsUpdateConfirmOpen(false);
    await handleUpdateConfig(editableConfig);
  };

  const handleConfirmRestart = async () => {
    const result = await handleRestartSystem();
    if (result.success) {
      setIsRestartConfirmOpen(false);
    }
  };

  return (
    <div className="w-full mx-auto space-y-4">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <SystemConfigForm
          config={editableConfig}
          updateConfig={updateConfig}
          onUpdate={handleUpdateClick}
          onRestart={handleRestartClick}
          isUpdating={isUpdating}
          isRestarting={isRestarting}
        />
      </motion.div>

      <Modal
        isOpen={isUpdateConfirmOpen}
        onClose={() => setIsUpdateConfirmOpen(false)}
        title={t('admin.system.update.title')}
        variant="confirm"
        confirmText={t('admin.system.update.confirm')}
        cancelText={t('admin.system.update.cancel')}
        onConfirm={handleConfirmUpdate}
        confirmType="danger"
      >
        <div className="space-y-2 text-sm text-neutral-300">
          <p>{t('admin.system.update.prompt')}</p>
          <p className="text-amber-300">{t('admin.system.update.warning')}</p>
        </div>
      </Modal>

      <Modal
        isOpen={isRestartConfirmOpen}
        onClose={() => setIsRestartConfirmOpen(false)}
        title={t('admin.system.restart.manualTitle')}
        variant="confirm"
        confirmText={t('admin.system.restart.confirm')}
        cancelText={t('admin.system.restart.cancel')}
        onConfirm={handleConfirmRestart}
        confirmType="danger"
      >
        <div className="space-y-2 text-sm text-neutral-300">
          <p>{t('admin.system.restart.manualPrompt')}</p>
          <p className="text-amber-300">{t('admin.system.restart.warning')}</p>
          {isRestarting && <p className="text-neutral-400">{t('admin.system.restart.inProgress')}</p>}
        </div>
      </Modal>
    </div>
  );
}

export default SystemConfig;
