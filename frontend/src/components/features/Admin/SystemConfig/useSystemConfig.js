import { useCallback, useEffect, useState } from 'react';
import { useImmerState } from '@/hooks/useImmerState';
import { normalizeConfig } from '@/components/features/Admin/SystemConfig/configNormalizer';
import { buildPayload } from '@/components/features/Admin/SystemConfig/configPayloadBuilder';
import { restartSystem, updateSystemConfig, uploadGeoCityDB } from '@/api/admin/system';
import { toast } from '@/utils/toast';

/** Hook for managing system configuration. */
export function useSystemConfig(initialConfig, t) {
  const [config, updateConfig] = useImmerState(null);
  const [isUpdating, setIsUpdating] = useState(false);
  const [isUploadingGeoCityDB, setIsUploadingGeoCityDB] = useState(false);
  const [isRestarting, setIsRestarting] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (initialConfig) {
      const normalized = normalizeConfig(initialConfig);
      updateConfig(() => normalized);
    }
  }, [initialConfig, updateConfig]);

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

  const handleUploadGeoCityDB = useCallback(
    async (file) => {
      if (!file) {
        return { success: false, error: 'No file provided' };
      }

      setIsUploadingGeoCityDB(true);
      setError(null);

      try {
        const response = await uploadGeoCityDB(file);
        if (response.code === 200) {
          toast.success({ description: t('admin.system.toast.uploadGeoCityDBSuccess') });
          return { success: true };
        }
        const errorMsg = response.msg || t('admin.system.toast.uploadGeoCityDBFailed');
        toast.danger({ description: errorMsg });
        setError(errorMsg);
        return { success: false, error: errorMsg };
      } catch (err) {
        const errorMsg = err.message || t('admin.system.toast.uploadGeoCityDBFailed');
        toast.danger({ description: errorMsg });
        setError(errorMsg);
        return { success: false, error: errorMsg };
      } finally {
        setIsUploadingGeoCityDB(false);
      }
    },
    [t]
  );

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
    isUploadingGeoCityDB,
    isRestarting,
    error,
    handleUpdateConfig,
    handleUploadGeoCityDB,
    handleRestartSystem,
  };
}
