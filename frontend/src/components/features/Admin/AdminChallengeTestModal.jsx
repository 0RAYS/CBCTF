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

  // 倒计时定时器
  const timerRef = useRef(null);

  // 获取测试状态
  const fetchTestStatus = useCallback(
    async (showSuccessToast = false) => {
      if (!challenge?.id) return;

      setLoading((prev) => ({ ...prev, status: true }));
      try {
        const response = await getTestChallengeStatus(challenge.id);
        if (response.code === 200) {
          setTestStatus(response.data);
          // 根据实际API数据结构设置剩余时间
          if (response.data.remote?.remaining) {
            setTimeLeft(response.data.remote.remaining);
          }
          // 如果是手动刷新，显示成功提示
          if (showSuccessToast) {
            toast.success({ description: t('admin.challenge.testModal.toast.statusRefreshSuccess') });
          }
        }
      } catch (error) {
        toast.danger({ description: error.message || t('admin.challenge.testModal.toast.fetchStatusFailed') });
      } finally {
        setLoading((prev) => ({ ...prev, status: false }));
      }
    },
    [challenge?.id]
  );

  // 当弹窗打开时获取状态
  useEffect(() => {
    if (isOpen && challenge?.id) {
      fetchTestStatus(false);
    }
  }, [isOpen, challenge?.id, fetchTestStatus]);

  // 当弹窗关闭时清理状态
  useEffect(() => {
    if (!isOpen) {
      setTestStatus(null);
      setTimeLeft(0);
      setIsCopied({});
      // 清理定时器
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    }
  }, [isOpen]);

  // 初始化时间
  useEffect(() => {
    if (testStatus?.remote?.remaining) {
      setTimeLeft(testStatus.remote.remaining);
    }
  }, [testStatus?.remote?.remaining]);

  // 倒计时效果
  useEffect(() => {
    if (!testStatus || !isOpen) return;
    // 根据实际API数据结构判断靶机是否运行
    const isRunning = testStatus.remote?.status === 'Running';
    if (!isRunning || !timeLeft) return;

    if (timerRef.current) {
      clearInterval(timerRef.current);
    }

    timerRef.current = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 0) {
          clearInterval(timerRef.current);
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
  }, [testStatus?.remote?.status, isOpen, timeLeft]);

  // 处理异步操作的通用函数
  const handleAsyncAction = useCallback(
    async (actionType, action, ...args) => {
      setLoading((prev) => ({ ...prev, [actionType]: true }));
      try {
        const response = await action(...args);
        if (response.code === 200) {
          // 操作成功后刷新状态
          await fetchTestStatus(false);
          // 显示成功提示
          toast.success({ description: t('admin.challenge.testModal.toast.actionSuccess') });
        }
      } catch (error) {
        toast.danger({ description: error.message || t('admin.challenge.testModal.toast.actionFailed') });
      } finally {
        setLoading((prev) => ({ ...prev, [actionType]: false }));
      }
    },
    [fetchTestStatus]
  );

  // 下载测试附件
  const handleDownloadAttachment = useCallback(async () => {
    if (!challenge?.id) return;

    setLoading((prev) => ({ ...prev, downloading: true }));
    try {
      const response = await downloadTestAttachment(challenge.id);
      downloadBlobResponse(response, 'attachment.zip', 'application/octet-stream');
      toast.success({ description: t('admin.challenge.testModal.toast.downloadSuccess') });
    } catch (error) {
      toast.danger({ description: error.message || t('admin.challenge.testModal.toast.downloadFailed') });
    } finally {
      setLoading((prev) => ({ ...prev, downloading: false }));
    }
  }, [challenge?.id]);

  // 启动测试靶机
  const handleStartVictim = () => {
    handleAsyncAction('starting', startTestVictim, challenge.id);
  };

  // 停止测试靶机
  const handleStopVictim = () => {
    handleAsyncAction('stopping', stopTestVictim, challenge.id);
  };

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

  // 靶机部分的渲染 - 完全仿照用户比赛时的样式
  const renderInstanceContent = useCallback(() => {
    if (!testStatus) return null;

    // 根据实际API数据结构判断靶机状态
    const isRunning = testStatus.remote?.status === 'Running';
    const targets = testStatus.remote?.target || [];

    return (
      <div className="space-y-3">
        {/* 状态行 - 包含状态和操作按钮 */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            {/* 状态指示器 */}
            <div className="flex items-center gap-2">
              <span className={`w-2 h-2 rounded-full ${isRunning ? 'bg-green-400' : 'bg-neutral-400'}`} />
              <span className="text-neutral-50 font-mono text-sm">
                {isRunning
                  ? t('admin.challenge.testModal.instance.running')
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
              /* 启动按钮 */
              <Button
                variant="primary"
                size="sm"
                onClick={handleStartVictim}
                disabled={loading.starting}
                loading={loading.starting}
                className={loading.starting ? 'border-yellow-400 text-yellow-400' : ''}
              >
                {loading.starting
                  ? t('admin.challenge.testModal.actions.launching')
                  : t('admin.challenge.testModal.actions.launch')}
              </Button>
            ) : (
              /* 运行中显示停止按钮 */
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

        {/* 运行中状态 - 进度条 */}
        {isRunning && (
          <div className="h-1.5 bg-neutral-700 rounded-full overflow-hidden">
            <motion.div
              className="h-full bg-yellow-400"
              initial={{ width: 0 }}
              animate={{
                // 假设总时长为1小时（3600秒），根据剩余时间计算进度
                width: `${Math.max(0, Math.min(100, (timeLeft / 3600) * 100))}%`,
              }}
              transition={{ duration: 0.5 }}
            />
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
