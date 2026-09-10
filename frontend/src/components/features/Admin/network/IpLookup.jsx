import { useTranslation } from 'react-i18next';
import IpLookupDialog from './IpLookupDialog';
import useIpLookup from './useIpLookup';

export default function IpLookup({ ip, className = '' }) {
  const { t } = useTranslation();
  const { lookup, dialogProps } = useIpLookup(ip);

  if (!ip) return <span className={`text-neutral-300 font-mono ${className}`}>-</span>;

  return (
    <>
      <button
        type="button"
        className={`text-geek-400 hover:text-geek-300 cursor-pointer transition-colors font-mono rounded-sm focus-visible:outline-2 focus-visible:outline-geek-400 ${className}`}
        aria-label={`${t('admin.contests.cheats.ipDetail.title')}: ${ip}`}
        onClick={(event) => {
          event.stopPropagation();
          lookup(ip);
        }}
      >
        {ip}
      </button>
      <IpLookupDialog {...dialogProps} />
    </>
  );
}
