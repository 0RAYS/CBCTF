/**
 * 题目详情模态框组件
 * @param {Object} props
 * @param {Object} props.challenge - 题目信息对象
 * @param {string} props.challenge.category - 题目类别 (WEB/CRYPTO/PWN等)
 * @param {string} props.challenge.title - 题目标题
 * @param {number} props.challenge.score - 题目分值
 * @param {boolean} props.challenge.isInitialized - 题目是否已初始化
 * @param {string} props.challenge.description - 题目描述（初始化后可见）
 * @param {Object} props.challenge.attachment - 附件（初始化后可见）
 * @param {string} props.challenge.attachment.name - 附件名称
 * @param {string} props.challenge.attachment.url - 附件下载链接
 * @param {string} props.challenge.attachment.size - 附件大小
 * @param {boolean} props.challenge.hasInstance - 是否有靶机
 * @param {boolean} props.challenge.instanceRunning - 靶机是否运行中
 * @param {string} [props.challenge.instanceIP] - 靶机地址（运行状态时可见）
 * @param {number} [props.challenge.instanceDuration] - 靶机运行时长（秒）
 * @param {number} [props.challenge.instanceTimeLeft] - 靶机剩余时间（秒）
 * @param {Object} props.contest - 比赛信息对象
 * @param {string} props.contest.prefix - flag前缀
 * @param {boolean} props.isOpen - 控制模态框显示/隐藏
 * @param {Function} props.onClose - 关闭模态框的回调函数
 * @param {Function} props.onInitialize - 初始化题目的回调函数, 返回Promise
 * @param {Function} props.onLaunchInstance - 启动靶机的回调函数, 返回Promise
 * @param {Function} props.onExtendInstance - 延长靶机时间的回调函数, 返回Promise
 * @param {Function} props.onDestroyInstance - 销毁靶机的回调函数, 返回Promise
 * @param {Function} props.onSubmitFlag - 提交flag的回调函数, 返回Promise
 * @param {Function} props.onDownloadAttachment - 下载附件的回调函数, 参数为附件对象
 */

import { useState, useEffect, useRef } from 'react';
import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Button, Modal } from '../../../../components/common';
import { IconCheck, IconCopy, IconPaperclip } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';

const normalizeInstanceStatus = (status) => {
  const normalizedStatus = typeof status === 'string' ? status.toLowerCase() : '';
  if (
    normalizedStatus === 'waiting' ||
    normalizedStatus === 'pending' ||
    normalizedStatus === 'terminating' ||
    normalizedStatus === 'running'
  ) {
    return normalizedStatus;
  }
  return '';
};

// 将 HintItem 组件移到外部, 使用 React.memo 包装以避免不必要的重新渲染
const HintItem = React.memo(({ hint, index }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const { t } = useTranslation();

  return (
    <div
      className={`
                border border-neutral-300/30 rounded-md overflow-hidden
                transition-colors duration-200
                ${isExpanded ? 'bg-neutral-900/50' : 'bg-black/20'}
            `}
    >
      <button
        type="button"
        aria-expanded={isExpanded}
        className="flex w-full items-center justify-between p-2 text-left cursor-pointer hover:bg-neutral-800/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-geek-400"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <div className="flex items-center gap-2">
          <span className="text-geek-400 font-mono text-sm">#{index + 1}</span>
          <span className="text-neutral-300 font-mono text-sm">
            {isExpanded ? t('game.challengeModal.hint.hide') : t('game.challengeModal.hint.show')}
          </span>
        </div>
        <span className="text-neutral-400 text-[10px]" aria-hidden="true">
          ▼
        </span>
      </button>

      <div hidden={!isExpanded}>
        <div className="p-3 border-t border-neutral-300/10">
          <span className="break-words text-neutral-400 font-mono text-sm">{hint}</span>
        </div>
      </div>
    </div>
  );
});

// 添加 displayName
HintItem.displayName = 'HintItem';

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

  // 状态管理
  const [loading, setLoading] = useState({
    initializing: false,
    resetting: false,
    launching: false,
    extending: false,
    destroying: false,
    submitting: false,
  });
  const [error, setError] = useState(null);
  const [flagError, setFlagError] = useState(null);
  const [flag, setFlag] = useState('');
  const [isCopied, setIsCopied] = useState({});
  const [timeLeft, setTimeLeft] = useState(0);

  // 修改倒计时效果, 使用 useRef 来避免不必要的重新渲染
  const timerRef = useRef(null);
  const sessionRef = useRef(null);
  const pendingRef = useRef(new Set());
  const copyTimeoutsRef = useRef(new Set());
  const instanceStatus = normalizeInstanceStatus(challenge?.instanceStatus);
  const isRunning = instanceStatus === 'running';
  const isWaiting = instanceStatus === 'waiting';
  const isPending = instanceStatus === 'pending';
  const isTerminating = instanceStatus === 'terminating';
  const instanceDuration = Number(challenge?.instanceDuration) || 0;
  const progressWidth = instanceDuration > 0 ? Math.max(0, Math.min(100, (timeLeft / instanceDuration) * 100)) : 0;
  const attemptsExhausted =
    Number(challenge?.maxAttempts) > 0 && Number(challenge?.attempts) >= Number(challenge.maxAttempts);
  const flagDisabled =
    !challenge?.isInitialized || challenge.isSolved || attemptsExhausted || loading.submitting || loading.resetting;

  useEffect(() => {
    sessionRef.current = {};
    const timeouts = copyTimeoutsRef.current;
    return () => {
      sessionRef.current = null;
      timeouts.forEach(clearTimeout);
      timeouts.clear();
    };
  }, [challenge?.id, isOpen]);

  // 初始化时间
  useEffect(() => {
    setTimeLeft(Number(challenge?.instanceTimeLeft) || 0);
  }, [challenge?.instanceTimeLeft]);

  // 倒计时效果 - 优化以减少重新渲染
  useEffect(() => {
    if (!challenge || !isOpen) return;
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
    if (!isRunning || timeLeft <= 0) return;

    timerRef.current = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(timerRef.current);
          timerRef.current = null;
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
    };
  }, [challenge, isOpen, isRunning, timeLeft]);

  // 格式化剩余时间
  const formatTimeLeft = (seconds) => {
    seconds = Math.floor(seconds); // 确保是整数
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60); // 确保秒数也是整数
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  // 处理异步操作的通用函数
  const handleAsyncAction = async (actionType, action, ...args) => {
    if (pendingRef.current.has(actionType)) return;
    if (actionType === 'resetting' && pendingRef.current.has('submitting')) return;
    const session = sessionRef.current;
    pendingRef.current.add(actionType);
    setError(null);
    setLoading((prev) => ({ ...prev, [actionType]: true }));
    try {
      if (!(await action(...args)) && sessionRef.current === session) {
        setError(t('errors.requestFailed'));
      }
    } catch (err) {
      if (sessionRef.current === session) setError(err.message || t('errors.requestFailed'));
    } finally {
      pendingRef.current.delete(actionType);
      if (sessionRef.current === session) {
        setLoading((prev) => ({ ...prev, [actionType]: false }));
      }
    }
  };

  // 初始化题目
  const handleInitialize = () => {
    handleAsyncAction('initializing', onInitialize, challenge.id);
  };

  // 重置题目
  const handleReset = () => {
    handleAsyncAction('resetting', onReset, challenge.id);
  };

  const handleLaunchInstance = () => handleAsyncAction('launching', onLaunchInstance, challenge.id);

  // 延长靶机时间
  const handleExtendTime = () => handleAsyncAction('extending', onExtendInstance, challenge.id);

  // 销毁靶机
  const handleDestroy = () => handleAsyncAction('destroying', onDestroyInstance, challenge.id);

  // 处理flag提交
  const handleSubmitFlag = async (e) => {
    e.preventDefault();
    if (!flag.trim() || flagDisabled || pendingRef.current.has('submitting') || pendingRef.current.has('resetting'))
      return;
    const session = sessionRef.current;
    pendingRef.current.add('submitting');

    setLoading((prev) => ({ ...prev, submitting: true }));
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
    } catch (err) {
      if (sessionRef.current === session) setFlagError(err.message || t('errors.requestFailed'));
    } finally {
      pendingRef.current.delete('submitting');
      if (sessionRef.current === session) setLoading((prev) => ({ ...prev, submitting: false }));
    }
  };

  // 复制IP地址
  const handleCopyIP = async (ip) => {
    const session = sessionRef.current;
    try {
      await navigator.clipboard.writeText(ip);
      if (sessionRef.current !== session) return;
      setIsCopied((prev) => ({ ...prev, [ip]: true }));
      const timeout = setTimeout(() => {
        setIsCopied((prev) => ({ ...prev, [ip]: false }));
        copyTimeoutsRef.current.delete(timeout);
      }, 2000);
      copyTimeoutsRef.current.add(timeout);
    } catch (err) {
      if (sessionRef.current === session) setError(err.message || t('errors.requestFailed'));
    }
  };

  // 靶机部分的渲染
  const renderInstanceContent = () => {
    const launchButtonLabel = isWaiting
      ? t('game.challengeModal.instance.waiting')
      : isTerminating
        ? t('game.challengeModal.instance.terminating')
        : isPending
          ? t('game.challengeModal.actions.launching')
          : t('game.challengeModal.actions.launch');

    return (
      <div className="space-y-3">
        {/* 状态行 */}
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div className="flex flex-wrap items-center gap-4">
            {/* 状态指示器 */}
            <div className="flex items-center gap-2">
              <span
                className={`w-2 h-2 rounded-full transition-colors duration-300 ${
                  isRunning
                    ? 'bg-green-400'
                    : isTerminating
                      ? 'bg-orange-400'
                      : isPending
                        ? 'bg-yellow-400'
                        : isWaiting
                          ? 'bg-yellow-400'
                          : 'bg-neutral-500'
                }`}
              />
              <span className="text-neutral-50 font-mono text-sm">
                {isRunning
                  ? t('game.challengeModal.instance.running')
                  : isTerminating
                    ? t('game.challengeModal.instance.terminating')
                    : isWaiting
                      ? t('game.challengeModal.instance.waiting')
                      : isPending
                        ? t('game.challengeModal.instance.pending')
                        : t('game.challengeModal.instance.notRunning')}
              </span>
            </div>

            {/* 运行中时显示剩余时间 */}
            {isRunning && (
              <div className="flex items-center gap-1.5">
                <span className="text-neutral-400 text-xs">{t('game.challengeModal.instance.time')}</span>
                <span className="text-yellow-400 font-mono text-sm">{formatTimeLeft(timeLeft)}</span>
              </div>
            )}
          </div>

          {/* 操作按钮 */}
          <div className="flex flex-wrap items-center gap-2">
            {!isRunning ? (
              <Button
                variant="primary"
                size="sm"
                onClick={handleLaunchInstance}
                disabled={loading.launching || isWaiting || isPending || isTerminating}
                loading={loading.launching}
                className={isWaiting || isPending || isTerminating ? 'border-yellow-400 text-yellow-400' : ''}
              >
                {launchButtonLabel}
              </Button>
            ) : (
              <>
                <Button
                  variant="primary"
                  size="sm"
                  onClick={handleExtendTime}
                  disabled={loading.extending}
                  loading={loading.extending}
                  className="border-yellow-400 text-yellow-400 hover:bg-yellow-400/10"
                >
                  {loading.extending
                    ? t('game.challengeModal.actions.extending')
                    : t('game.challengeModal.actions.extend')}
                </Button>
                <Button
                  variant="danger"
                  size="sm"
                  onClick={handleDestroy}
                  disabled={loading.destroying}
                  loading={loading.destroying}
                >
                  {loading.destroying
                    ? t('game.challengeModal.actions.destroying')
                    : t('game.challengeModal.actions.destroy')}
                </Button>
              </>
            )}
          </div>
        </div>

        {/* Transition states use a static bar; running shows remaining time. */}
        {(isRunning || isWaiting || isPending || isTerminating) && (
          <div className="h-1.5 bg-neutral-700 rounded-full overflow-hidden">
            {isRunning ? (
              <div className="h-full bg-yellow-400" style={{ width: `${progressWidth}%` }} />
            ) : isWaiting ? (
              <div className="h-full w-full bg-yellow-400/35" />
            ) : isTerminating ? (
              <div className="h-full w-full bg-orange-400/35" />
            ) : (
              <div className="h-full w-full bg-yellow-400/35" />
            )}
          </div>
        )}

        {/* 靶机地址 - 只在运行中显示 */}
        {isRunning && challenge.instanceIP && (
          <div>
            <div className="flex items-center justify-between p-1 rounded-md">
              <span className="text-neutral-400 text-xs">{t('game.challengeModal.instance.address')}</span>
            </div>
            {challenge.instanceIP.map((ip, index) => (
              <div key={index} className="flex min-w-0 items-center justify-between gap-2 bg-neutral-900">
                <span className="min-w-0 break-all font-mono text-neutral-50 text-sm">{ip}</span>
                <Button
                  variant="ghost"
                  size="icon"
                  className="shrink-0 !text-neutral-400 hover:!text-geek-400"
                  aria-label={t('common.copyToClipboard', {
                    defaultValue: 'Copy to clipboard',
                  })}
                  onClick={() => handleCopyIP(ip)}
                >
                  {isCopied[ip] ? <IconCheck size={16} /> : <IconCopy size={16} />}
                </Button>
              </div>
            ))}
          </div>
        )}
      </div>
    );
  };

  // 在渲染 hint 部分时, 为每个 hint 提供一个稳定的 key
  const renderHints = () => {
    if (!challenge.hints || challenge.hints.length === 0) return null;

    return (
      <div className="space-y-1.5">
        <h3 className="text-neutral-400 font-mono text-sm">{t('game.challengeModal.sections.hints')}</h3>
        <div className="space-y-2">
          {challenge.hints.map((hint, index) => (
            <HintItem key={`hint-${challenge.id}-${index}`} hint={hint} index={index} />
          ))}
        </div>
      </div>
    );
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
          <div className="min-w-0 w-full space-y-1.5">
            <label htmlFor="challenge-flag" className="text-neutral-400 font-mono text-sm">
              {t('game.challengeModal.sections.submitFlag')}
            </label>
            <form
              onSubmit={handleSubmitFlag}
              className="flex flex-col gap-2 sm:flex-row"
              aria-busy={loading.submitting}
            >
              <input
                id="challenge-flag"
                type="text"
                value={flag}
                autoComplete="off"
                autoCapitalize="none"
                spellCheck={false}
                readOnly={flagDisabled}
                aria-invalid={!!flagError}
                aria-describedby={flagError ? 'challenge-flag-error' : undefined}
                onChange={(e) => {
                  setFlag(e.target.value);
                  setFlagError(null);
                }}
                placeholder={`${contest.prefix}{...}`}
                className="min-w-0 w-full sm:flex-1 h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4
                text-neutral-50 placeholder-neutral-400 focus:border-geek-400 focus:shadow-focus transition-colors duration-200"
              />
              <Button
                type="submit"
                variant="primary"
                size="action"
                className="!h-auto min-h-10 whitespace-normal break-words sm:shrink-0"
                loading={loading.submitting}
                disabled={flagDisabled || !flag.trim()}
              >
                {loading.submitting
                  ? t('game.challengeModal.submit.submitting')
                  : t('game.challengeModal.submit.button', {
                      status: challenge.isSolved ? t('common.solved') : t('common.submit'),
                      attempts: challenge.attempts,
                      max: challenge.maxAttempts || '∞',
                    })}
              </Button>
            </form>
            {(flagError || error) && (
              <div className="max-h-20 space-y-1 overflow-y-auto text-red-400 text-sm [overflow-wrap:anywhere]">
                {flagError && (
                  <div id="challenge-flag-error" role="alert">
                    {flagError}
                  </div>
                )}
                {error && <div role="alert">{error}</div>}
              </div>
            )}
          </div>
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
            onClick={handleInitialize}
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
          {/* 描述 */}
          <div className="flex flex-col items-start justify-between gap-4 sm:flex-row">
            <div className="w-full space-y-1.5 flex-1 min-w-0">
              <h3 className="text-neutral-400 font-mono text-sm">{t('game.challengeModal.sections.description')}</h3>
              <div className="break-words text-neutral-50 prose prose-invert prose-sm max-w-none [&_pre]:max-w-full [&_pre]:overflow-x-auto [&_table]:block [&_table]:overflow-x-auto">
                <ReactMarkdown remarkPlugins={[remarkGfm]}>{challenge.description || ''}</ReactMarkdown>
              </div>
            </div>
            <div className="flex-shrink-0">
              <Button
                variant="primary"
                size="action"
                onClick={handleReset}
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

          {/* 附件 - 仅在有附件时显示 */}
          {challenge.attachment && (
            <div className="space-y-1.5">
              <h3 className="text-neutral-400 font-mono text-sm">{t('game.challengeModal.sections.attachments')}</h3>
              <div className="space-y-1.5">
                <button
                  type="button"
                  className="flex w-full items-center gap-2 p-2 text-left border border-neutral-300/30 rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400
                                                        text-neutral-300 hover:text-geek-400 hover:border-geek-400
                                                        transition-colors duration-200 cursor-pointer"
                  onClick={(e) => {
                    e.preventDefault();
                    onDownloadAttachment(challenge.attachment);
                  }}
                >
                  <IconPaperclip size={16} className="shrink-0" aria-hidden="true" />
                  <span className="min-w-0 break-all font-mono text-sm">{challenge.attachment}</span>
                </button>
              </div>
            </div>
          )}

          {/* 提示 - 调整间距 */}
          {renderHints()}

          {/* 靶机信息 - 调整间距 */}
          {challenge.hasInstance && renderInstanceContent()}
        </div>
      )}
    </Modal>
  );
}

export default ChallengeModal;
