import { motion, AnimatePresence } from 'motion/react';
import { useState, useEffect, useRef, useCallback } from 'react';
import { Button } from '../../../components/common';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { toast } from '../../../utils/toast';
import { downloadBlobResponse } from '../../../utils/fileDownload';
import {
  getTestChallengeStatus,
  downloadTestAttachment,
  startTestVictim,
  stopTestVictim,
} from '../../../api/admin/challenge';
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

// 格式化剩余时间
const formatTimeLeft = (seconds) => {
  seconds = Math.floor(seconds);
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
};

/**
 * 管理员题目测试弹窗组件
 * @param {Object} props
 * @param {Object} props.challenge - 题目信息对象
 * @param {boolean} props.isOpen - 控制弹窗显示/隐藏
 * @param {Function} props.onClose - 关闭弹窗的回调函数
 */
function AdminChallengeTestModal({ challenge, isOpen, onClose }) {
  const { t } = useTranslation();

  // 状态管理
  const [loading, setLoading] = useState({
    status: false,
    downloading: false,
    starting: false,
    stopping: false,
  });
  const [testStatus, setTestStatus] = useState(null);
  const [timeLeft, setTimeLeft] = useState(0);
  const [isCopied, setIsCopied] = useState({});

  // 定时器 & 轮询 refs
  const timerRef = useRef(null);
  const pollingIntervalRef = useRef(null);
  const pollingTimeoutRef = useRef(null);
  const isOpenRef = useRef(isOpen);

  // 保持 isOpenRef 与 isOpen 同步，供轮询回调使用
  useEffect(() => {
    isOpenRef.current = isOpen;
  }, [isOpen]);

  const stopPolling = useCallback(() => {
    if (pollingIntervalRef.current) {
      clearInterval(pollingIntervalRef.current);
      pollingIntervalRef.current = null;
    }
    if (pollingTimeoutRef.current) {
      clearTimeout(pollingTimeoutRef.current);
      pollingTimeoutRef.current = null;
    }
  }, []);

  const currentInstanceStatus = normalizeInstanceStatus(testStatus?.remote?.status);
  const isRunning = currentInstanceStatus === 'running';
  const isWaiting = currentInstanceStatus === 'waiting';
  const isPending = currentInstanceStatus === 'pending';
  const isTerminating = currentInstanceStatus === 'terminating';
  const instanceDuration = Number(testStatus?.remote?.duration) || 0;
  const progressWidth = instanceDuration > 0 ? Math.max(0, Math.min(100, (timeLeft / instanceDuration) * 100)) : 0;

  // 仅在弹窗打开期间轮询；同步 waiting/pending/running 状态
  const startPolling = useCallback(
    (challengeId) => {
      stopPolling();
      pollingIntervalRef.current = setInterval(async () => {
        if (!isOpenRef.current) {
          stopPolling();
          return;
        }
        try {
          const response = await getTestChallengeStatus(challengeId);
          if (response.code === 200) {
            const nextStatus = normalizeInstanceStatus(response.data.remote?.status);
            setTestStatus(response.data);
            setTimeLeft(Number(response.data.remote?.remaining) || 0);
            if (nextStatus === 'running') {
              stopPolling();
              setLoading((prev) => ({ ...prev, starting: false }));
            }
          }
        } catch {
          // 忽略轮询中的网络错误
        }
      }, 5000);
      // 3 分钟超时兜底
      pollingTimeoutRef.current = setTimeout(
        () => {
          stopPolling();
          setLoading((prev) => ({ ...prev, starting: false }));
        },
        3 * 60 * 1000
      );
    },
    [stopPolling]
  );

  // 获取测试状态
  const fetchTestStatus = useCallback(
    async (showSuccessToast = false) => {
      if (!challenge?.id) return;

      setLoading((prev) => ({ ...prev, status: true }));
      try {
        const response = await getTestChallengeStatus(challenge.id);
        if (response.code === 200) {
          setTestStatus(response.data);
          setTimeLeft(Number(response.data.remote?.remaining) || 0);
          if (showSuccessToast) {
            toast.success({ description: t('admin.challenge.testModal.toast.statusRefreshSuccess') });
          }
          // Pod 仍在排队或启动中（页面刷新后恢复）→ 自动开始轮询
          if (['waiting', 'pending', 'terminating'].includes(normalizeInstanceStatus(response.data.remote?.status))) {
            startPolling(challenge.id);
          }
        }
      } catch (error) {
        toast.danger({ description: error.message || t('admin.challenge.testModal.toast.fetchStatusFailed') });
      } finally {
        setLoading((prev) => ({ ...prev, status: false }));
      }
    },
    [challenge?.id, startPolling, t]
  );

  // 当弹窗打开时获取初始状态
  useEffect(() => {
    if (isOpen && challenge?.id) {
      fetchTestStatus(false);
    }
  }, [isOpen, challenge?.id, fetchTestStatus]);

  // 当弹窗关闭时停止轮询并重置所有状态
  useEffect(() => {
    if (!isOpen) {
      stopPolling();
      setLoading({ status: false, downloading: false, starting: false, stopping: false });
      setTestStatus(null);
      setTimeLeft(0);
      setIsCopied({});
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    }
  }, [isOpen, stopPolling]);

  // 初始化时间
  useEffect(() => {
    setTimeLeft(Number(testStatus?.remote?.remaining) || 0);
  }, [testStatus?.remote?.remaining]);

  // 倒计时效果
  useEffect(() => {
    if (!testStatus || !isOpen) return;
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
  }, [isOpen, isRunning, testStatus, timeLeft]);

  // 启动测试靶机 — HTTP 成功后保持 loading，通过轮询等待 Running
  const handleStartVictim = useCallback(async () => {
    setLoading((prev) => ({ ...prev, starting: true }));
    try {
      const response = await startTestVictim(challenge.id);
      if (response.code === 200) {
        setTestStatus((prev) => ({
          ...prev,
          remote: {
            ...(prev?.remote || {}),
            status: 'waiting',
            target: prev?.remote?.target || [],
            duration: prev?.remote?.duration || 0,
            remaining: 0,
          },
        }));
        toast.success({ description: t('admin.challenge.testModal.toast.actionSuccess') });
        startPolling(challenge.id);
      } else {
        setLoading((prev) => ({ ...prev, starting: false }));
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.challenge.testModal.toast.actionFailed') });
      setLoading((prev) => ({ ...prev, starting: false }));
    }
  }, [challenge?.id, startPolling, t]);

  // 停止测试靶机 — 后端同步执行，HTTP 返回后直接刷新状态
  const handleStopVictim = useCallback(async () => {
    setLoading((prev) => ({ ...prev, stopping: true }));
    try {
      const response = await stopTestVictim(challenge.id);
      if (response.code === 200) {
        await fetchTestStatus(false);
        toast.success({ description: t('admin.challenge.testModal.toast.actionSuccess') });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.challenge.testModal.toast.actionFailed') });
    } finally {
      setLoading((prev) => ({ ...prev, stopping: false }));
    }
  }, [challenge?.id, fetchTestStatus, t]);

  // 下载测试附件
  const handleDownloadAttachment = useCallback(async () => {
    if (!challenge?.id) return;

    setLoading((prev) => ({ ...prev, downloading: true }));
    try {
      const response = await downloadTestAttachment(challenge.id);
      if (response.headers?.['file'] === 'true') {
        downloadBlobResponse(response, 'attachment.zip', 'application/octet-stream');
        toast.success({ description: t('admin.challenge.testModal.toast.downloadSuccess') });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.challenge.testModal.toast.downloadFailed') });
    } finally {
      setLoading((prev) => ({ ...prev, downloading: false }));
    }
  }, [challenge?.id, t]);

  // 复制IP地址
  const handleCopyIP = useCallback(
    (ip) => {
      navigator.clipboard.writeText(ip);
      setIsCopied((prev) => ({ ...prev, [ip]: true }));
      toast.success({ description: t('admin.challenge.testModal.toast.copyIpSuccess') });
      setTimeout(() => {
        setIsCopied((prev) => ({ ...prev, [ip]: false }));
      }, 2000);
    },
    [t]
  );

  // 靶机部分的渲染
  const renderInstanceContent = useCallback(() => {
    const targets = testStatus?.remote?.target || [];
    const launchButtonLabel = isWaiting
      ? t('admin.challenge.testModal.instance.waiting')
      : isTerminating
        ? t('admin.challenge.testModal.instance.terminating')
        : isPending
          ? t('admin.challenge.testModal.actions.launching')
          : t('admin.challenge.testModal.actions.launch');

    if (!testStatus && !loading.starting) return null;

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
                          : 'bg-neutral-400'
                }`}
              />
              <span className="text-neutral-50 font-mono text-sm">
                {isRunning
                  ? t('admin.challenge.testModal.instance.running')
                  : isTerminating
                    ? t('admin.challenge.testModal.instance.terminating')
                    : isWaiting
                      ? t('admin.challenge.testModal.instance.waiting')
                      : isPending
                        ? t('admin.challenge.testModal.instance.pending')
                        : t('admin.challenge.testModal.instance.notRunning')}
              </span>
            </div>

            {/* 运行中时显示剩余时间 */}
            {isRunning && (
              <div className="flex items-center gap-2">
                <span className="text-neutral-400 text-sm">{t('admin.challenge.testModal.instance.time')}</span>
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
                onClick={handleStartVictim}
                disabled={loading.starting || isWaiting || isPending || isTerminating}
                className={isWaiting || isPending || isTerminating ? 'border-yellow-400 text-yellow-400' : ''}
              >
                {launchButtonLabel}
              </Button>
            ) : (
              <Button
                variant="danger"
                size="sm"
                onClick={handleStopVictim}
                disabled={loading.stopping}
                loading={loading.stopping}
              >
                {loading.stopping
                  ? t('admin.challenge.testModal.actions.stopping')
                  : t('admin.challenge.testModal.actions.stop')}
              </Button>
            )}
          </div>
        </div>

        {/* 进度条：waiting 显示静态条，pending 显示闪动条，running 显示倒计时 */}
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
        {isRunning && targets.length > 0 && (
          <div>
            <div className="flex items-center justify-between p-2 bg-neutral-900 rounded-md">
              <span className="text-neutral-400 text-xs">{t('admin.challenge.testModal.instance.address')}</span>
            </div>
            {targets.map((target, index) => (
              <div key={index} className="flex justify-between p-1.5">
                <span className="font-mono text-neutral-50 text-sm cursor-pointer" onClick={() => handleCopyIP(target)}>
                  {target}
                </span>
                <Button
                  variant="ghost"
                  size="icon"
                  className="!text-neutral-400 hover:!text-geek-400"
                  onClick={() => handleCopyIP(target)}
                >
                  {isCopied[target] ? '✓' : '📋'}
                </Button>
              </div>
            ))}
          </div>
        )}
      </div>
    );
  }, [
    testStatus,
    timeLeft,
    loading.starting,
    loading.stopping,
    handleStartVictim,
    handleStopVictim,
    handleCopyIP,
    isCopied,
    isPending,
    isRunning,
    isTerminating,
    isWaiting,
    progressWidth,
    t,
  ]);

  if (!isOpen) return null;

  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[900] flex items-center justify-center">
          {/* 背景遮罩 - 确保完全覆盖 */}
          <motion.div
            className="fixed inset-0 bg-black/60 backdrop-blur-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
          />

          {/* 内容容器 - 添加内边距 */}
          <div className="relative z-10 w-full max-w-[800px] p-4">
            {/* 模态框内容 */}
            <motion.div
              className="relative w-full bg-black/80 border border-neutral-300 rounded-md"
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
            >
              {/* 头部 - 减小内边距 */}
              <div className="p-5 border-b border-neutral-300/30">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <span className="text-geek-400 font-mono">{t('admin.challenge.testModal.title')}</span>
                    <h2 className="text-2xl text-neutral-50 font-mono">{challenge?.name}</h2>
                    <span className="text-yellow-400 font-mono">{challenge?.type}</span>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="!text-neutral-400 hover:!text-neutral-50"
                    onClick={onClose}
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
                      {t('admin.challenge.testModal.sections.description')}
                    </h3>
                    <div className="text-neutral-50 prose prose-invert prose-sm max-w-none">
                      <ReactMarkdown remarkPlugins={[remarkGfm]}>{challenge?.description || ''}</ReactMarkdown>
                    </div>
                  </div>
                </div>

                {/* Dynamic类型 - 下载测试附件 */}
                {(challenge?.type === 'dynamic' || testStatus?.file !== '') && (
                  <div className="space-y-1.5">
                    <h3 className="text-neutral-400 font-mono text-sm">
                      {t('admin.challenge.testModal.sections.attachments')}
                    </h3>
                    <div className="space-y-1.5">
                      <motion.div
                        className="flex items-center gap-2 p-2 border border-neutral-300/30 rounded-md
                                  text-neutral-300 hover:text-geek-400 hover:border-geek-400
                                  transition-colors duration-200 cursor-pointer"
                        whileHover={{ x: 5 }}
                        onClick={handleDownloadAttachment}
                      >
                        <span className="text-sm">📎</span>
                        <span className="font-mono text-sm">{testStatus?.file}</span>
                        <span className="text-neutral-400 text-sm ml-auto">
                          {loading.downloading
                            ? t('admin.challenge.testModal.attachments.generating')
                            : challenge?.type === 'dynamic'
                              ? t('admin.challenge.testModal.attachments.generateAndDownload')
                              : t('admin.challenge.testModal.attachments.download')}
                        </span>
                      </motion.div>
                    </div>
                  </div>
                )}

                {/* Pods类型 - 靶机控制 */}
                {challenge?.type === 'pods' && (
                  <div className="space-y-1.5">
                    <h3 className="text-neutral-400 font-mono text-sm">
                      {t('admin.challenge.testModal.sections.instance')}
                    </h3>
                    {renderInstanceContent()}
                  </div>
                )}
              </div>
            </motion.div>
          </div>
        </div>
      )}
    </AnimatePresence>
  );
}

export default AdminChallengeTestModal;
