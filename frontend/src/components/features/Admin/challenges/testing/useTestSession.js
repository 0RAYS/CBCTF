import { useEffect, useEffectEvent, useLayoutEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getTestChallengeStatus, startTestVictim, stopTestVictim } from '../../../../../api/admin/challenge';
import { toast } from '../../../../../utils/toast';
import { createTestSession, initialTestState, normalizeInstanceStatus } from './testSession';

const api = { getTestChallengeStatus, startTestVictim, stopTestVictim };

export default function useTestSession(challengeId) {
  const { t } = useTranslation();
  const [state, setState] = useState(initialTestState);
  const sessionRef = useRef(null);
  const notify = useEffectEvent((type, key, error) => {
    toast[type]({ description: error?.message || t(`admin.challenge.testModal.toast.${key}`) });
  });

  useLayoutEffect(() => {
    const session = createTestSession({ challengeId, api, onChange: setState, notify });
    sessionRef.current = session;
    setState(initialTestState());
    void session.load();
    return () => session.dispose();
  }, [challengeId]);

  const remaining = Math.max(0, Number(state.testStatus?.remote?.remaining) || 0);
  const running = normalizeInstanceStatus(state.testStatus?.remote?.status) === 'running';
  const [timeLeft, setTimeLeft] = useState(0);
  useEffect(() => {
    setTimeLeft(remaining);
    if (!running || remaining <= 0) return;
    const expiresAt = state.receivedAt + remaining * 1000;
    const timer = setInterval(() => {
      const seconds = Math.max(0, Math.ceil((expiresAt - Date.now()) / 1000));
      setTimeLeft(seconds);
      if (seconds === 0) clearInterval(timer);
    }, 1000);
    return () => clearInterval(timer);
  }, [remaining, running, state.receivedAt]);

  return {
    ...state,
    timeLeft,
    start: () => sessionRef.current?.start(),
    stop: () => sessionRef.current?.stop(),
  };
}
