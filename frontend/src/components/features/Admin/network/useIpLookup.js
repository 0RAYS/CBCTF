import { useLayoutEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getIpInfo } from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';

const CLOSED = { isOpen: false, data: null, loading: false };

export default function useIpLookup(scope) {
  const { t } = useTranslation();
  const [state, setState] = useState(CLOSED);
  const request = useRef(0);

  useLayoutEffect(() => {
    setState(CLOSED);
    return () => {
      request.current += 1;
    };
  }, [scope]);

  const close = () => {
    request.current += 1;
    setState(CLOSED);
  };

  const lookup = async (ip) => {
    if (!ip) return;
    const current = ++request.current;
    setState({ scope, isOpen: true, data: { ip }, loading: true });
    try {
      const response = await getIpInfo(ip);
      if (current !== request.current) return;
      setState({
        scope,
        isOpen: true,
        data: { ...(response.code === 200 ? response.data : {}), ip },
        loading: false,
      });
    } catch (error) {
      // Closing, changing scope, unmounting, or a newer lookup invalidates both data and errors.
      if (current !== request.current) return;
      toast.danger({ description: error.message || t('admin.contests.cheats.ipDetail.fetchFailed') });
      setState({ scope, isOpen: true, data: { ip }, loading: false });
    }
  };

  return { lookup, close, dialogProps: { ...state, isOpen: state.isOpen && state.scope === scope, onClose: close } };
}
