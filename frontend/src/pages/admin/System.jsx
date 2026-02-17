import { useState, useEffect } from 'react';
import { toast } from '../../utils/toast';
import SystemConfig from '../../components/features/Admin/SystemConfig';
import { getSystemConfig } from '../../api/admin/system';
import { useTranslation } from 'react-i18next';

function SystemSettings() {
  const [config, setConfig] = useState(null);
  const { t } = useTranslation();

  const fetchConfig = async () => {
    try {
      const response = await getSystemConfig();
      if (response.code === 200) {
        setConfig(response.data);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.system.toast.fetchFailed') });
    }
  };

  useEffect(() => {
    fetchConfig();
  }, []);

  if (!config) {
    return <div className="p-4 text-red-500">{t('admin.system.noData')}</div>;
  }

  return <SystemConfig config={config} />;
}

export default SystemSettings;
