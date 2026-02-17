import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getIpInfo } from '../../api/admin/contest';
import { toast } from '../../utils/toast';
import Modal from './Modal';
import Card from './Card';

/**
 * IP 地址显示组件，点击后弹出 IP 地理位置反查对话框
 * @param {Object} props
 * @param {string} props.ip - IP 地址
 * @param {string} [props.className] - 额外的 CSS 类名
 */
function IpAddress({ ip, className = '' }) {
  const { t } = useTranslation();
  const [show, setShow] = useState(false);
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);

  if (!ip) return <span className={`text-neutral-300 font-mono ${className}`}>-</span>;

  const handleClick = async (e) => {
    e.stopPropagation();
    setLoading(true);
    setShow(true);
    try {
      const response = await getIpInfo(ip);
      if (response.code === 200) {
        setData({ ip, ...response.data });
      } else {
        setData({ ip });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.cheats.ipDetail.fetchFailed') });
      setData({ ip });
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setShow(false);
    setData(null);
  };

  return (
    <>
      <span
        className={`text-geek-400 hover:text-geek-300 cursor-pointer transition-colors font-mono ${className}`}
        onClick={handleClick}
      >
        {ip}
      </span>

      <Modal isOpen={show} onClose={handleClose} title={t('admin.contests.cheats.ipDetail.title')} size="sm">
        {loading ? (
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-geek-400"></div>
          </div>
        ) : data ? (
          <Card variant="default" padding="md">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.ipDetail.ip')}
                </label>
                <p className="text-neutral-300 font-mono">{data.ip || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.ipDetail.iso')}
                </label>
                <p className="text-neutral-300">{data.iso || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.ipDetail.country')}
                </label>
                <p className="text-neutral-300">{data.country || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.ipDetail.subdivision')}
                </label>
                <p className="text-neutral-300">{data.subdivision || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.ipDetail.city')}
                </label>
                <p className="text-neutral-300">{data.city || '-'}</p>
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.ipDetail.timezone')}
                </label>
                <p className="text-neutral-300">{data.timezone || '-'}</p>
              </div>
              <div className="col-span-2">
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.ipDetail.coordinates')}
                </label>
                <p className="text-neutral-300 font-mono">
                  {data.latitude != null && data.longitude != null ? `${data.latitude}, ${data.longitude}` : '-'}
                </p>
              </div>
            </div>
          </Card>
        ) : null}
      </Modal>
    </>
  );
}

export default IpAddress;
