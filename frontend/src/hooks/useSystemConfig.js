import { useEffect, useState, useCallback } from 'react';
import { useImmerState } from './useImmerState';
import { normalizeConfig } from '../utils/configNormalizer';
import { buildPayload } from '../utils/configPayloadBuilder';
import { updateSystemConfig, restartSystem } from '../api/admin/system';
import { toast } from '../utils/toast';

/**
 * Hook for managing system configuration
 * @param {Object} initialConfig - Initial config from parent component
 * @param {Function} t - Translation function
 * @returns {Object} - Config state and methods
 */
export function useSystemConfig(initialConfig, t) {
  const [config, updateConfig] = useImmerState(null);
  const [isUpdating, setIsUpdating] = useState(false);
  const [isRestarting, setIsRestarting] = useState(false);
  const [error, setError] = useState(null);

  // Normalize and set config when initialConfig changes
  useEffect(() => {
    if (initialConfig) {
      const normalized = normalizeConfig(initialConfig);
      updateConfig(() => normalized);
    }
  }, [initialConfig, updateConfig]);

  // Update config on server
  const handleUpdateConfig = useCallback(
    async (configToUpdate) => {
      if (!configToUpdate) {
        return { success: false, error: 'No config provided' };
      }

      setIsUpdating(true);
      setError(null);

      try {
        const payload = buildPayload(configToUpdate);
        const response = await updateSystemConfig(payload);

        if (response.code === 200) {
          toast.success({ description: t('admin.system.toast.updateSuccess') });
          return { success: true };
        }
        const errorMsg = response.msg || t('admin.system.toast.updateFailed');
        toast.danger({ description: errorMsg });
        setError(errorMsg);
        return { success: false, error: errorMsg };
      } catch (err) {
        const errorMsg = err.message || t('admin.system.toast.updateFailed');
        toast.danger({ description: errorMsg });
        setError(errorMsg);
        return { success: false, error: errorMsg };
      } finally {
        setIsUpdating(false);
      }
    },
    [t]
  );

  // Restart system
  const handleRestartSystem = useCallback(async () => {
    setIsRestarting(true);
    setError(null);

    try {
      const response = await restartSystem();

      if (response.code === 200) {
        toast.success({ description: t('admin.system.toast.restartSuccess') });
        return { success: true };
      }
      const errorMsg = response.msg || t('admin.system.toast.restartFailed');
      toast.danger({ description: errorMsg });
      setError(errorMsg);
      return { success: false, error: errorMsg };
    } catch (err) {
      const errorMsg = err.message || t('admin.system.toast.restartFailed');
      toast.danger({ description: errorMsg });
      setError(errorMsg);
      return { success: false, error: errorMsg };
    } finally {
      setIsRestarting(false);
    }
  }, [t]);

  return {
    config,
    updateConfig,
    isUpdating,
    isRestarting,
    error,
    handleUpdateConfig,
    handleRestartSystem,
  };
}
