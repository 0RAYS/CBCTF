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

import { motion, AnimatePresence, useAnimationControls } from 'motion/react';
import { useState, useEffect, useRef } from 'react';
import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Button } from '../../../../components/common';
import { useTranslation } from 'react-i18next';
import { EASE_T2, backdropVariants } from '../../../../config/motion';

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
    <motion.div
      className={`
                border border-neutral-300/30 rounded-md overflow-hidden
                transition-colors duration-200
                ${isExpanded ? 'bg-neutral-900/50' : 'bg-black/20'}
            `}
    >
      <div
        className="flex items-center justify-between p-2 cursor-pointer hover:bg-neutral-800/30"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <div className="flex items-center gap-2">
          <span className="text-geek-400 font-mono text-sm">#{index + 1}</span>
          <span className="text-neutral-300 font-mono text-sm">
            {isExpanded ? t('game.challengeModal.hint.hide') : t('game.challengeModal.hint.show')}
          </span>
        </div>
        <motion.span
          className="text-neutral-400 text-[10px]"
          animate={{ rotate: isExpanded ? 180 : 0 }}
          transition={{ duration: 0.2 }}
        >
          ▼
        </motion.span>
      </div>

      <motion.div
        initial={false}
        animate={{
          height: isExpanded ? 'auto' : 0,
          opacity: isExpanded ? 1 : 0,
        }}
        transition={{ duration: 0.2 }}
        className="overflow-hidden"
      >
        <div className="p-3 border-t border-neutral-300/10">
          <span className="text-neutral-400 font-mono text-sm">{hint}</span>
        </div>
      </motion.div>
    </motion.div>
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
  const flagControls = useAnimationControls();

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
  const [flag, setFlag] = useState('');
  const [selectedOptions, setSelectedOptions] = useState([]);
  const [isCopied, setIsCopied] = useState({});
  const [timeLeft, setTimeLeft] = useState(0);

  // 修改倒计时效果, 使用 useRef 来避免不必要的重新渲染
  const timerRef = useRef(null);
  const prevRunningRef = useRef(normalizeInstanceStatus(challenge?.instanceStatus) === 'running');
  const launchingTimeoutRef = useRef(null);
  const instanceStatus = normalizeInstanceStatus(challenge?.instanceStatus);
  const isRunning = instanceStatus === 'running';
  const isWaiting = instanceStatus === 'waiting';
  const isPending = instanceStatus === 'pending';
  const isTerminating = instanceStatus === 'terminating';
  const instanceDuration = Number(challenge?.instanceDuration) || 0;
  const progressWidth = instanceDuration > 0 ? Math.max(0, Math.min(100, (timeLeft / instanceDuration) * 100)) : 0;

  // Clear launching state when instanceRunning transitions false → true (via polling or WS)
  useEffect(() => {
    const prev = prevRunningRef.current;
    prevRunningRef.current = isRunning;
    if (!prev && isRunning) {
      if (launchingTimeoutRef.current) {
        clearTimeout(launchingTimeoutRef.current);
        launchingTimeoutRef.current = null;
      }
      setLoading((p) => ({ ...p, launching: false }));
    }
  }, [isRunning]);

  // Flag 提交错误时抖动动画
  useEffect(() => {
    if (error) {
      flagControls.start({
        x: [-8, 8, -6, 6, -3, 3, 0],
        transition: { duration: 0.35 },
      });
    }
  }, [error]);

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
    setError(null);
    try {
      if (await action(...args)) {
        setLoading((prev) => ({ ...prev, [actionType]: true }));
      }
    } catch (err) {
      setError(err.message);
    } finally {
      if (actionType !== 'launching') {
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

  // 启动靶机 — 从点击瞬间到 Pod Ready 全程 loading, 3 分钟超时兜底
  const handleLaunchInstance = async () => {
    setError(null);
    setLoading((p) => ({ ...p, launching: true }));
    try {
      const ok = await onLaunchInstance(challenge.id);
      if (!ok) {
        setLoading((p) => ({ ...p, launching: false }));
      } else {
        // HTTP 成功, 等待 Pod Ready（由 Fix3 / WS 清除）, 3 分钟后强制清除
        launchingTimeoutRef.current = setTimeout(
          () => {
            setLoading((p) => ({ ...p, launching: false }));
            launchingTimeoutRef.current = null;
          },
          3 * 60 * 1000
        );
      }
    } catch (err) {
      setError(err.message);
      setLoading((p) => ({ ...p, launching: false }));
    }
  };

  // 延长靶机时间
  const handleExtendTime = async () => {
    setError(null);
    setLoading((p) => ({ ...p, extending: true }));
    try {
      await onExtendInstance(challenge.id);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading((p) => ({ ...p, extending: false }));
    }
  };

  // 销毁靶机
  const handleDestroy = async () => {
    setError(null);
    setLoading((p) => ({ ...p, destroying: true }));
    try {
      await onDestroyInstance(challenge.id);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading((p) => ({ ...p, destroying: false }));
    }
  };

  // 处理选项选择
  const handleOptionToggle = (randId) => {
    setSelectedOptions((prev) => {
      if (prev.includes(randId)) {
        return prev.filter((id) => id !== randId);
      } else {
        return [...prev, randId];
      }
    });
  };

  // 处理flag提交
  const handleSubmitFlag = async (e) => {
    e.preventDefault();

    // 如果是question类型, 检查是否选择了选项
    if (challenge.type === 'question') {
      if (selectedOptions.length === 0) {
        setError(t('game.challengeModal.errors.selectOption'));
        return;
      }
      // 将选中的选项rand_id用逗号拼接
      const value = selectedOptions.join(',');
      setLoading((prev) => ({ ...prev, submitting: true }));
      setError(null);

      try {
        const result = await onSubmitFlag(challenge.id, value);
        if (result.success) {
          setSelectedOptions([]);
          onClose();
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading((prev) => ({ ...prev, submitting: false }));
      }
    } else {
      // 原有的flag提交逻辑
      if (!flag.trim()) return;

      setLoading((prev) => ({ ...prev, submitting: true }));
      setError(null);

      try {
        const result = await onSubmitFlag(challenge.id, flag);
        if (result.success) {
          setFlag('');
          onClose();
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading((prev) => ({ ...prev, submitting: false }));
      }
    }
  };

  // 复制IP地址
  const handleCopyIP = (ip) => {
    navigator.clipboard.writeText(ip);
    setIsCopied((prev) => ({ ...prev, [ip]: true }));
    setTimeout(() => {
      setIsCopied((prev) => ({ ...prev, [ip]: false }));
    }, 2000);
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
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            {/* 状态指示器 */}
            <div className="flex items-center gap-2">
              <span
                className={`w-2 h-2 rounded-full transition-colors duration-300 ${
                  isRunning
                    ? 'bg-green-400'
                    : isTerminating
                      ? 'bg-orange-400 animate-pulse'
                      : isPending
                        ? 'bg-yellow-400 animate-pulse'
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
          <div className="flex items-center gap-2">
            {!isRunning ? (
              <Button
                variant="primary"
                size="sm"
                onClick={handleLaunchInstance}
                disabled={loading.launching || isWaiting || isPending || isTerminating}
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

        {/* 进度条: waiting 显示静态条, pending 显示闪动条, running 显示倒计时 */}
        {(isRunning || isWaiting || isPending || isTerminating) && (
          <div className="h-1.5 bg-neutral-700 rounded-full overflow-hidden">
            {isRunning ? (
              <motion.div
                className="h-full bg-yellow-400"
                initial={{ width: 0 }}
                animate={{ width: `${progressWidth}%` }}
                transition={{ duration: 0.5 }}
              />
            ) : isWaiting ? (
              <div className="h-full w-full bg-yellow-400/35" />
            ) : isTerminating ? (
              <motion.div
                className="h-full w-2/5 bg-orange-400/70 rounded-full"
                animate={{ x: ['-100%', '350%'] }}
                transition={{ duration: 1.1, repeat: Infinity, ease: 'easeInOut' }}
              />
            ) : (
              <motion.div
                className="h-full w-2/5 bg-yellow-400/60 rounded-full"
                animate={{ x: ['-100%', '350%'] }}
                transition={{ duration: 1.4, repeat: Infinity, ease: 'easeInOut' }}
              />
            )}
          </div>
        )}

        {/* 靶机地址 - 只在运行中显示 */}
        {isRunning && challenge.instanceIP && (
          <div>
            <div className="p-2 bg-neutral-900 rounded-md">
              <span className="text-neutral-400 text-xs">{t('game.challengeModal.instance.address')}</span>
            </div>
            {challenge.instanceIP.map((ip, index) => (
              <div key={index} className="flex items-center justify-between p-1.5">
                <span className="font-mono text-neutral-50 text-sm">{ip}</span>
                <Button
                  variant="ghost"
                  size="icon"
                  className="!text-neutral-400 hover:!text-geek-400"
                  onClick={() => handleCopyIP(ip)}
                >
                  {isCopied[ip] ? '✓' : '📋'}
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

  // 渲染问题选项
  const renderQuestionOptions = () => {
    if (challenge.type !== 'question' || !challenge.options || challenge.options.length === 0) return null;

    return (
      <div className="space-y-1.5">
        <h3 className="text-neutral-400 font-mono text-sm">{t('game.challengeModal.sections.questionOptions')}</h3>
        <div className="space-y-2">
          {challenge.options.map((option, index) => (
            <div
              key={option.rand_id || index}
              className="flex items-center gap-3 p-3 border border-neutral-300/30 rounded-md hover:border-geek-400/50 transition-colors duration-200 cursor-pointer"
              onClick={() => handleOptionToggle(option.rand_id)}
            >
              <div className="flex items-center justify-center w-5 h-5 border-2 border-neutral-400 rounded transition-colors duration-200">
                {selectedOptions.includes(option.rand_id) && <div className="w-2 h-2 bg-geek-400 rounded-sm" />}
              </div>
              <span className="text-neutral-50 font-mono text-sm flex-1">{option.content}</span>
            </div>
          ))}
        </div>
      </div>
    );
  };

  if (!isOpen) return null;

  // 未初始化状态下的内容
  if (!challenge.isInitialized) {
    return (
      <div className="fixed inset-0 z-[900] flex items-center justify-center">
        <motion.div
          className="fixed inset-0 bg-neutral-900/70 backdrop-blur-sm"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          onClick={() => {
            setFlag('');
            onClose();
          }}
        />

        <div className="relative z-10 w-full max-w-[800px] p-4">
          <motion.div
            className="relative w-full bg-neutral-800/90 border border-neutral-600/60 rounded-md"
            initial={{ scale: 0.97, opacity: 0, y: 8 }}
            animate={{ scale: 1, opacity: 1, y: 0 }}
            exit={{ scale: 0.97, opacity: 0, y: 8 }}
            transition={{ type: 'tween', ease: EASE_T2, duration: 0.22 }}
          >
            {/* 头部 */}
            <div className="p-5 border-b border-neutral-600/50">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <span className="text-geek-400 font-mono">{challenge.category}</span>
                  <h2 className="text-2xl text-neutral-50 font-mono">{challenge.title}</h2>
                  <span className="text-yellow-400 font-mono">{t('common.points', { count: challenge.score })}</span>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="!text-neutral-400 hover:!text-neutral-50"
                  onClick={() => {
                    setFlag('');
                    onClose();
                  }}
                >
                  ✕
                </Button>
              </div>
            </div>

            {/* 初始化提示 */}
            <div className="p-12 flex flex-col items-center justify-center space-y-4">
              {error && (
                <motion.div
                  className="text-red-400 text-sm"
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                >
                  {error}
                </motion.div>
              )}
              <span className="text-neutral-400 text-sm">{t('game.challengeModal.initialize.message')}</span>
              <Button
                variant="primary"
                size="action"
                onClick={handleInitialize}
                disabled={loading.initializing}
                loading={loading.initializing}
                className={loading.initializing ? 'border-yellow-400 text-yellow-400' : ''}
              >
                {loading.initializing
                  ? t('game.challengeModal.initialize.loading')
                  : t('game.challengeModal.initialize.action')}
              </Button>
            </div>
          </motion.div>
        </div>
      </div>
    );
  }

  // 已初始化状态, 显示完整内容
  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[900] flex items-center justify-center">
          {/* 背景遮罩 - 确保完全覆盖 */}
          <motion.div
            className="fixed inset-0 bg-neutral-900/70 backdrop-blur-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => {
              setFlag('');
              onClose();
            }}
          />

          {/* 内容容器 - 添加内边距 */}
          <div className="relative z-10 w-full max-w-[800px] p-4">
            {/* 模态框内容 */}
            <motion.div
              className="relative w-full bg-black/80 border border-neutral-300 rounded-md"
              initial={{ scale: 0.97, opacity: 0, y: 8 }}
              animate={{ scale: 1, opacity: 1, y: 0 }}
              exit={{ scale: 0.97, opacity: 0, y: 8 }}
              transition={{ type: 'tween', ease: [0.25, 1, 0.5, 1], duration: 0.22 }}
            >
              {/* 头部 - 减小内边距 */}
              <div className="p-5 border-b border-neutral-600/50">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <span className="text-geek-400 font-mono">{challenge.category}</span>
                    <h2 className="text-2xl text-neutral-50 font-mono">{challenge.title}</h2>
                    <span className="text-yellow-400 font-mono">{t('common.points', { count: challenge.score })}</span>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="!text-neutral-400 hover:!text-neutral-50"
                    onClick={() => {
                      setFlag('');
                      onClose();
                    }}
                  >
                    ✕
                  </Button>
                </div>
              </div>

              {/* 主体内容 - 调整间距 */}
              <div className="p-5 space-y-5">
                {/* 描述 */}
                <div className="flex items-start justify-between gap-4">
                  <div className="space-y-1.5 flex-1 min-w-0">
                    <h3 className="text-neutral-400 font-mono text-sm">
                      {t('game.challengeModal.sections.description')}
                    </h3>
                    <div className="text-neutral-50 prose prose-invert prose-sm max-w-none">
                      <ReactMarkdown remarkPlugins={[remarkGfm]}>{challenge.description || ''}</ReactMarkdown>
                    </div>
                  </div>
                  <div className="flex-shrink-0">
                    <Button
                      variant="primary"
                      size="action"
                      onClick={handleReset}
                      disabled={loading.resetting}
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
                    <h3 className="text-neutral-400 font-mono text-sm">
                      {t('game.challengeModal.sections.attachments')}
                    </h3>
                    <div className="space-y-1.5">
                      <motion.a
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-2 p-2 border border-neutral-300/30 rounded-md
                                                        text-neutral-300 hover:text-geek-400 hover:border-geek-400
                                                        transition-colors duration-200 cursor-pointer"
                        whileHover={{ x: 5 }}
                        onClick={(e) => {
                          e.preventDefault();
                          onDownloadAttachment(challenge.attachment);
                        }}
                      >
                        <span className="text-sm">📎</span>
                        <span className="font-mono text-sm">{challenge.attachment}</span>
                      </motion.a>
                    </div>
                  </div>
                )}

                {/* 提示 - 调整间距 */}
                {renderHints()}

                {/* 靶机信息 - 调整间距 */}
                {challenge.hasInstance && renderInstanceContent()}

                {/* 问题选项 - question类型才显示 */}
                {renderQuestionOptions()}

                {/* Flag 提交 - 调整间距 */}
                <motion.div className="space-y-1.5" animate={flagControls}>
                  <h3 className="text-neutral-400 font-mono text-sm">
                    {challenge.type === 'question'
                      ? t('game.challengeModal.sections.submitAnswer')
                      : t('game.challengeModal.sections.submitFlag')}
                  </h3>
                  {challenge.type === 'question' ? (
                    <form onSubmit={handleSubmitFlag} className="flex gap-2">
                      <div className="flex-1 h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4 flex items-center">
                        <span className="text-neutral-400 font-mono text-sm">
                          {selectedOptions.length > 0
                            ? t('game.challengeModal.submit.selectedOptions', { count: selectedOptions.length })
                            : t('game.challengeModal.submit.selectOptions')}
                        </span>
                      </div>
                      <Button
                        type="submit"
                        variant="primary"
                        size="action"
                        disabled={challenge.isSolved || loading.submitting || selectedOptions.length === 0}
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
                  ) : (
                    <form onSubmit={handleSubmitFlag} className="flex gap-2">
                      <input
                        type="text"
                        onChange={(e) => setFlag(e.target.value)}
                        placeholder={`${contest.prefix}{...}`}
                        className="flex-1 h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4
                                                  text-neutral-50 placeholder-neutral-400
                                                  focus:border-geek-400 focus:shadow-focus
                                                  transition-all duration-200"
                      />
                      <Button
                        type="submit"
                        variant="primary"
                        size="action"
                        disabled={challenge.isSolved || loading.submitting || !flag.trim()}
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
                  )}
                  {error && (
                    <motion.div
                      className="text-red-400 text-sm mt-2"
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                    >
                      {error}
                    </motion.div>
                  )}
                </motion.div>
              </div>
            </motion.div>
          </div>
        </div>
      )}
    </AnimatePresence>
  );
}

export default ChallengeModal;
