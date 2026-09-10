import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getWriteups, uploadWriteup } from '../../../../api/challenge';
import { toast } from '../../../../utils/toast';

// Callers own access checks and whether an upload replaces the whole page with loading.
export default function useContestWriteups(contestId) {
  const { t } = useTranslation();
  const [writeups, setWriteups] = useState([]);
  const scopeRef = useRef(null);
  const requestRef = useRef(0);

  useEffect(() => {
    scopeRef.current = { contestId };
    return () => {
      scopeRef.current = null;
      requestRef.current += 1;
    };
  }, [contestId]);

  const refresh = async () => {
    const scope = scopeRef.current;
    const request = ++requestRef.current;
    const isCurrent = () => scopeRef.current === scope && requestRef.current === request;
    try {
      const response = await getWriteups(contestId);
      if (!isCurrent()) return;
      if (response.code !== 200) throw new Error(response.msg || t('game.challenges.toast.fetchWriteupsFailed'));
      setWriteups(response.data?.writeups || []);
    } catch (error) {
      if (isCurrent()) toast.danger({ description: error.message || t('game.challenges.toast.fetchWriteupsFailed') });
    }
  };

  const upload = async (file) => {
    const scope = scopeRef.current;
    try {
      const response = await uploadWriteup(contestId, file);
      if (scopeRef.current !== scope) return;
      if (response.code !== 200) throw new Error(response.msg || t('game.challenges.toast.uploadFailed'));
      toast.success({
        title: t('game.challenges.toast.uploadSuccess'),
        description: t('game.challenges.toast.uploadThanks'),
      });
      void refresh();
    } catch (error) {
      if (scopeRef.current === scope) {
        toast.danger({ description: error.message || t('game.challenges.toast.uploadFailed') });
      }
    }
  };

  return { writeups, refresh, upload };
}
