import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from '../../../../utils/toast.js';
import { batchTone, getBatchResult } from './batchModel.js';

export default function useBatchAction(scopeKey) {
  const { t } = useTranslation();
  const sessionRef = useRef(null);
  const [state, setState] = useState({});
  useEffect(() => {
    const session = { active: true, busy: false };
    sessionRef.current = session;
    setState({ scopeKey, pending: false, result: null, error: '' });
    return () => {
      session.active = false;
    };
  }, [scopeKey]);

  const run = async (request, { nested = false, onResult, successMessage, failureMessage } = {}) => {
    const session = sessionRef.current;
    if (!session?.active || session.busy) return;
    session.busy = true;
    setState({ scopeKey, pending: true, result: null, error: '' });
    try {
      let response;
      try {
        response = await request();
      } catch {
        if (!session.active) return;
        const message = t('admin.batch.unknownOutcome');
        setState({ scopeKey, pending: false, result: null, error: message });
        toast.danger({ description: message });
        return;
      }
      if (!session.active) return;
      const result = getBatchResult(response, nested);
      if (result) {
        setState({ scopeKey, pending: false, result, error: '' });
        toast[batchTone(result)]({ description: t('admin.batch.summary', result) });
        onResult?.(result, response);
      } else if (response?.code === 200) {
        toast.success({ description: successMessage || response.msg || t('admin.batch.states.success') });
        onResult?.(null, response);
      } else {
        const message = response?.msg || failureMessage || t('errors.requestFailed');
        setState({ scopeKey, pending: false, result: null, error: message });
        toast.danger({ description: message });
      }
    } finally {
      session.busy = false;
      if (session.active) setState((current) => ({ ...current, pending: false }));
    }
  };

  return { ...(state.scopeKey === scopeKey ? state : { pending: false, result: null, error: '' }), run };
}
