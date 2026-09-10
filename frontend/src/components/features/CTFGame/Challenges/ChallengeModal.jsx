import { useState, useEffect, useRef } from 'react';
import { Button, Modal } from '../../../common';
import MarkdownContent from '../../../common/MarkdownContent';
import { IconPaperclip } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { normalizeInstanceStatus } from './models/challengeViewModel';
import FlagForm from './display/FlagForm';
import Hints from './display/Hints';
import InstancePanel from './display/InstancePanel';

function ChallengeModal({
  challenge,
  contest,
  isOpen,
  onClose,
  onInitialize,
  onReset,
  onLaunchInstance,
  onExtendInstance,
  onDestroyInstance,
  onSubmitFlag,
  onDownloadAttachment,
}) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState({});
  const [error, setError] = useState(null);
  const [flagError, setFlagError] = useState(null);
  const [flag, setFlag] = useState('');
  const [isCopied, setIsCopied] = useState({});
  const [timeLeft, setTimeLeft] = useState(0);
  const timerRef = useRef(null);
  const sessionRef = useRef(null);
  const pendingRef = useRef(new Set());
  const copyTimeoutsRef = useRef(new Map());
  const isRunning = normalizeInstanceStatus(challenge?.instanceStatus) === 'running';
  const attemptsExhausted =
    Number(challenge?.maxAttempts) > 0 && Number(challenge?.attempts) >= Number(challenge.maxAttempts);
  const flagDisabled =
    !challenge?.isInitialized || challenge.isSolved || attemptsExhausted || loading.submitting || loading.resetting;

  useEffect(() => {
    sessionRef.current = {};
    pendingRef.current = new Set();
    const timeouts = copyTimeoutsRef.current;
    return () => {
      sessionRef.current = null;
      timeouts.forEach(clearTimeout);
      timeouts.clear();
    };
  }, [challenge?.id, isOpen]);

  useEffect(() => {
    setTimeLeft(Number(challenge?.instanceTimeLeft) || 0);
  }, [challenge?.instanceTimeLeft]);

  useEffect(() => {
    if (!challenge || !isOpen || !isRunning || timeLeft <= 0) return;
    timerRef.current = setInterval(() => {
      setTimeLeft((previous) => Math.max(0, previous - 1));
    }, 1000);
    return () => {
      clearInterval(timerRef.current);
      timerRef.current = null;
    };
  }, [challenge, isOpen, isRunning, timeLeft]);

  // Keep synchronous guards and loading in one owner: reset and submit exclude each other.
  const handleAsyncAction = async (actionType, action) => {
    const pending = pendingRef.current;
    if (pending.has(actionType) || (actionType === 'resetting' && pending.has('submitting'))) return;
    const session = sessionRef.current;
    pending.add(actionType);
    setError(null);
    setLoading((previous) => ({ ...previous, [actionType]: true }));
    try {
      if (!(await action(challenge.id)) && sessionRef.current === session) setError(t('errors.requestFailed'));
    } catch (error) {
      if (sessionRef.current === session) setError(error.message || t('errors.requestFailed'));
    } finally {
      pending.delete(actionType);
      if (sessionRef.current === session) setLoading((previous) => ({ ...previous, [actionType]: false }));
    }
  };

  const handleSubmitFlag = async (event) => {
    event.preventDefault();
    const pending = pendingRef.current;
    if (!flag.trim() || flagDisabled || pending.has('submitting') || pending.has('resetting')) return;
    const session = sessionRef.current;
    pending.add('submitting');
    setLoading((previous) => ({ ...previous, submitting: true }));
    setFlagError(null);
    try {
      const result = await onSubmitFlag(challenge.id, flag);
      if (sessionRef.current !== session) return;
      if (result.success) {
        setFlag('');
        onClose();
      } else {
        setFlagError(result.message || t('errors.requestFailed'));
      }
    } catch (error) {
      if (sessionRef.current === session) setFlagError(error.message || t('errors.requestFailed'));
    } finally {
      pending.delete('submitting');
      if (sessionRef.current === session) setLoading((previous) => ({ ...previous, submitting: false }));
    }
  };

  const handleCopyIP = async (ip) => {
    const session = sessionRef.current;
    try {
      await navigator.clipboard.writeText(ip);
      if (sessionRef.current !== session) return;
      setIsCopied((previous) => ({ ...previous, [ip]: true }));
      clearTimeout(copyTimeoutsRef.current.get(ip));
      const timeout = setTimeout(() => {
        if (sessionRef.current !== session) return;
        setIsCopied((previous) => ({ ...previous, [ip]: false }));
        copyTimeoutsRef.current.delete(ip);
      }, 2000);
      copyTimeoutsRef.current.set(ip, timeout);
    } catch (error) {
      if (sessionRef.current === session) setError(error.message || t('errors.requestFailed'));
    }
  };

  if (!isOpen || !challenge) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="lg"
      title={
        <span className="flex min-w-0 flex-wrap items-baseline gap-x-3 gap-y-1">
          <span className="break-all text-sm text-geek-400">{challenge.category}</span>
          <span className="order-last min-w-0 basis-full [overflow-wrap:anywhere] text-lg sm:order-none sm:basis-auto sm:text-xl">
            {challenge.title}
          </span>
          <span className="text-sm text-yellow-400">{t('common.points', { count: challenge.score })}</span>
        </span>
      }
      footer={
        challenge.isInitialized ? (
          <FlagForm
            challenge={challenge}
            prefix={contest.prefix}
            flag={flag}
            onChange={(value) => {
              setFlag(value);
              setFlagError(null);
            }}
            onSubmit={handleSubmitFlag}
            disabled={flagDisabled}
            submitting={loading.submitting}
            flagError={flagError}
            error={error}
          />
        ) : null
      }
    >
      {!challenge.isInitialized ? (
        <div className="flex flex-col items-center justify-center gap-4 py-8 text-center">
          {error && (
            <div role="alert" className="text-red-400 text-sm break-words">
              {error}
            </div>
          )}
          <p className="text-neutral-400 text-sm">{t('game.challengeModal.initialize.message')}</p>
          <Button
            variant="primary"
            size="action"
            onClick={() => handleAsyncAction('initializing', onInitialize)}
            disabled={loading.initializing}
            loading={loading.initializing}
          >
            {loading.initializing
              ? t('game.challengeModal.initialize.loading')
              : t('game.challengeModal.initialize.action')}
          </Button>
        </div>
      ) : (
        <div className="space-y-5">
          <div className="flex flex-col items-start justify-between gap-4 sm:flex-row">
            <div className="w-full space-y-1.5 flex-1 min-w-0">
              <h3 className="text-neutral-400 font-mono text-sm">{t('game.challengeModal.sections.description')}</h3>
              <MarkdownContent className="break-words text-neutral-50 prose prose-invert prose-sm max-w-none [&_pre]:max-w-full [&_pre]:overflow-x-auto [&_table]:block [&_table]:overflow-x-auto">
                {challenge.description || ''}
              </MarkdownContent>
            </div>
            <div className="flex-shrink-0">
              <Button
                variant="primary"
                size="action"
                onClick={() => handleAsyncAction('resetting', onReset)}
                disabled={loading.resetting || loading.submitting}
                loading={loading.resetting}
                className={loading.resetting ? 'border-yellow-400 text-yellow-400' : ''}
              >
                {loading.resetting
                  ? t('game.challengeModal.actions.resetting')
                  : t('game.challengeModal.actions.reset')}
              </Button>
            </div>
          </div>
          {challenge.attachment && (
            <div className="space-y-1.5">
              <h3 className="text-neutral-400 font-mono text-sm">{t('game.challengeModal.sections.attachments')}</h3>
              <div className="space-y-1.5">
                <button
                  type="button"
                  className="flex w-full items-center gap-2 p-2 text-left border border-neutral-300/30 rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400 text-neutral-300 hover:text-geek-400 hover:border-geek-400 transition-colors duration-200 cursor-pointer"
                  onClick={() => onDownloadAttachment(challenge.attachment)}
                >
                  <IconPaperclip size={16} className="shrink-0" aria-hidden="true" />
                  <span className="min-w-0 break-all font-mono text-sm">{challenge.attachment}</span>
                </button>
              </div>
            </div>
          )}
          <Hints challenge={challenge} />
          {challenge.hasInstance && (
            <InstancePanel
              challenge={challenge}
              timeLeft={timeLeft}
              loading={loading}
              onLaunch={() => handleAsyncAction('launching', onLaunchInstance)}
              onExtend={() => handleAsyncAction('extending', onExtendInstance)}
              onDestroy={() => handleAsyncAction('destroying', onDestroyInstance)}
              onCopy={handleCopyIP}
              isCopied={isCopied}
            />
          )}
        </div>
      )}
    </Modal>
  );
}

export default ChallengeModal;
