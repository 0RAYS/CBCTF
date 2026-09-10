import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  getChallengeStatus,
  initChallenge,
  resetChallenge,
  startRemoteTarget,
  extendContainerTime,
  stopContainer,
  submitFlag,
  downloadChallengeAttachment,
} from '../../../../../api/challenge';
import { toast } from '../../../../../utils/toast';
import { downloadBlobResponse } from '../../../../../utils/fileDownload';
import { isInstanceTransitioning, mapChallengeStatusToViewModel } from '../models/challengeViewModel';

export default function useChallengeSession(contestId, { updateChallenge, onSolved }) {
  const { t } = useTranslation();
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const selectedChallengeRef = useRef(null);
  const pollingIntervalRef = useRef(null);
  const pollingTimeoutRef = useRef(null);
  const pollingGenerationRef = useRef(0);
  const selectionRef = useRef(0);
  const statusRequestRef = useRef(0);
  const scopeRef = useRef(null);

  const applyStatus = (challenge) => {
    selectedChallengeRef.current = challenge;
    setSelectedChallenge(challenge);
    updateChallenge(challenge);
  };

  const stopPolling = () => {
    pollingGenerationRef.current += 1;
    clearTimeout(pollingIntervalRef.current);
    clearTimeout(pollingTimeoutRef.current);
    pollingIntervalRef.current = null;
    pollingTimeoutRef.current = null;
  };

  const startPolling = (challengeId, targetStatus, selection) => {
    if (selectionRef.current !== selection || selectedChallengeRef.current?.id !== challengeId) return;
    stopPolling();
    const scope = scopeRef.current;
    const generation = pollingGenerationRef.current;
    const isCurrent = () =>
      scopeRef.current === scope && generation === pollingGenerationRef.current && selection === selectionRef.current;
    // Serial completion-based scheduling prevents overlapping slow requests.
    const poll = async () => {
      if (!isCurrent()) return;
      const request = ++statusRequestRef.current;
      try {
        const response = await getChallengeStatus(contestId, challengeId, { noToast: true, noLoading: true });
        if (!isCurrent()) return;
        if (response.code === 200 && request === statusRequestRef.current) {
          const challenge = mapChallengeStatusToViewModel(selectedChallengeRef.current, response.data);
          applyStatus(challenge);
          if (
            (targetStatus === 'running' && challenge.instanceStatus === 'running') ||
            (targetStatus === 'stopped' &&
              challenge.instanceStatus !== 'running' &&
              !isInstanceTransitioning(challenge.instanceStatus))
          ) {
            stopPolling();
          }
        }
      } catch {
        // Polling failures are silent; retry until this generation times out.
      }
      if (isCurrent()) pollingIntervalRef.current = setTimeout(poll, 5000);
    };
    pollingIntervalRef.current = setTimeout(poll, 5000);
    pollingTimeoutRef.current = setTimeout(stopPolling, 3 * 60 * 1000);
  };

  const closeChallenge = () => {
    stopPolling();
    selectionRef.current += 1;
    statusRequestRef.current += 1;
    selectedChallengeRef.current = null;
    setSelectedChallenge(null);
  };

  useEffect(() => {
    scopeRef.current = { contestId };
    return () => {
      scopeRef.current = null;
      selectionRef.current += 1;
      statusRequestRef.current += 1;
      selectedChallengeRef.current = null;
      stopPolling();
    };
  }, [contestId]);

  const refreshChallengeStatus = async (challengeId, selection) => {
    if (selectionRef.current !== selection || selectedChallengeRef.current?.id !== challengeId) return;
    const scope = scopeRef.current;
    const request = ++statusRequestRef.current;
    const isCurrent = () =>
      scopeRef.current === scope && selectionRef.current === selection && request === statusRequestRef.current;
    try {
      const response = await getChallengeStatus(contestId, challengeId);
      if (!isCurrent()) return;
      if (response.code === 200)
        applyStatus(mapChallengeStatusToViewModel(selectedChallengeRef.current, response.data));
    } catch (error) {
      if (isCurrent()) toast.danger({ description: error.message || t('game.challenges.toast.refreshStatusFailed') });
    }
  };

  // All mutations capture the same selection owner, including their follow-up requests.
  const runMutation = async (challengeId, action, label, afterSuccess) => {
    const scope = scopeRef.current;
    const selection = selectionRef.current;
    const isCurrent = () =>
      scopeRef.current === scope &&
      selectionRef.current === selection &&
      selectedChallengeRef.current?.id === challengeId;
    if (!scope || !isCurrent()) return false;
    try {
      const response = await action(contestId, challengeId);
      if (!isCurrent()) return false;
      if (response.code !== 200) throw new Error(response.msg || t(`game.challenges.toast.${label}Failed`));
      toast.success({
        title: (label !== 'reset' && label !== 'launch' && response.msg) || t(`game.challenges.toast.${label}Success`),
      });
      if (afterSuccess) await afterSuccess(challengeId, selection);
      else await refreshChallengeStatus(challengeId, selection);
      return isCurrent();
    } catch (error) {
      if (isCurrent()) toast.danger({ title: t(`game.challenges.toast.${label}Failed`), description: error.message });
      return false;
    }
  };

  const handleLaunchInstance = (challengeId) =>
    runMutation(challengeId, startRemoteTarget, 'launch', (id, selection) => {
      statusRequestRef.current += 1;
      applyStatus({
        ...selectedChallengeRef.current,
        instanceStatus: 'waiting',
        instanceRunning: false,
        instancePending: false,
        instanceWaiting: true,
        instanceTerminating: false,
      });
      startPolling(id, 'running', selection);
    });

  const handleDestroyInstance = (challengeId) =>
    runMutation(challengeId, stopContainer, 'destroy', async (id, selection) => {
      await refreshChallengeStatus(id, selection);
      startPolling(id, 'stopped', selection);
    });

  const handleSubmitFlag = async (challengeId, value) => {
    const scope = scopeRef.current;
    const selection = selectionRef.current;
    const isCurrent = () => scopeRef.current === scope && selectionRef.current === selection;
    try {
      const response = await submitFlag(contestId, challengeId, { flag: value });
      if (!isCurrent()) return { success: false };
      // Incorrect flags still refresh attempts without remounting the input.
      void refreshChallengeStatus(challengeId, selection);
      if (response.code === 200) {
        onSolved();
        toast.success({ title: response.msg || t('game.challenges.toast.submitSuccess') });
        return { success: true, message: t('game.challenges.toast.submitSuccessMessage') };
      }
      return { success: false, message: response.msg || t('errors.requestFailed') };
    } catch (error) {
      return { success: false, message: isCurrent() ? error.message || t('errors.requestFailed') : undefined };
    }
  };

  const openChallenge = async (challenge) => {
    closeChallenge();
    const scope = scopeRef.current;
    const selection = selectionRef.current;
    const request = ++statusRequestRef.current;
    const isCurrent = () =>
      scopeRef.current === scope && selectionRef.current === selection && request === statusRequestRef.current;
    try {
      const response = await getChallengeStatus(contestId, challenge.id);
      if (!isCurrent()) return;
      if (response.code !== 200) throw new Error(response.msg || t('game.challenges.toast.fetchStatusFailed'));
      const updated = mapChallengeStatusToViewModel(challenge, response.data);
      selectedChallengeRef.current = updated;
      setSelectedChallenge(updated);
      if (isInstanceTransitioning(updated.instanceStatus)) {
        startPolling(challenge.id, updated.instanceStatus === 'terminating' ? 'stopped' : 'running', selection);
      }
    } catch (error) {
      if (!isCurrent()) return;
      selectedChallengeRef.current = challenge;
      setSelectedChallenge(challenge);
      toast.danger({ description: error.message || t('game.challenges.toast.fetchStatusFailed') });
    }
  };

  const handleDownloadAttachment = async (attachment) => {
    const scope = scopeRef.current;
    const selection = selectionRef.current;
    const challengeId = selectedChallengeRef.current?.id;
    if (challengeId == null) return;
    const isCurrent = () => scopeRef.current === scope && selectionRef.current === selection;
    try {
      const response = await downloadChallengeAttachment(contestId, challengeId);
      if (!isCurrent()) return;
      if (response.headers?.['file'] !== 'true') {
        throw new Error(response.msg || t('game.challenges.toast.downloadFailed'));
      }
      downloadBlobResponse(response, attachment, 'application/octet-stream');
      toast.success({ title: t('game.challenges.toast.downloadSuccess') });
    } catch (error) {
      if (isCurrent()) toast.danger({ description: error.message || t('game.challenges.toast.downloadFailed') });
    }
  };

  return {
    selectedChallenge,
    openChallenge,
    closeChallenge,
    modalProps: {
      challenge: selectedChallenge,
      isOpen: !!selectedChallenge,
      onClose: closeChallenge,
      onInitialize: (id) => runMutation(id, initChallenge, 'init'),
      onReset: (id) => runMutation(id, resetChallenge, 'reset'),
      onLaunchInstance: handleLaunchInstance,
      onExtendInstance: (id) => runMutation(id, extendContainerTime, 'extend'),
      onDestroyInstance: handleDestroyInstance,
      onSubmitFlag: handleSubmitFlag,
      onDownloadAttachment: handleDownloadAttachment,
    },
  };
}
