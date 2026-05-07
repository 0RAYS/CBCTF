import { motion } from 'motion/react';
import { IconX, IconPlus, IconTrash, IconUsers } from '@tabler/icons-react';
import { useState, useEffect } from 'react';
import { Button } from '../../../../components/common';
import { useTranslation } from 'react-i18next';
import ScoreCurveChart from './ScoreCurveChart';
import FlagSolversModal from './FlagSolversModal';

/**
 * 赛事内部赛题编辑弹窗组件
 * @param {Object} props
 * @param {boolean} props.isOpen - 是否显示弹窗
 * @param {string} props.mode - 模式, 'add'或'edit'
 * @param {Object} props.challenge - 当前编辑的赛题对象
 * @param {Array} props.flags - 赛题的flag列表
 * @param {Function} props.onClose - 关闭弹窗回调
 * @param {Function} props.onSubmit - 提交表单回调
 * @param {Function} props.onChange - 表单字段变更回调
 * @param {Function} props.onAddHint - 添加提示回调
 * @param {Function} props.onRemoveHint - 移除提示回调
 * @param {Function} props.onHintChange - 提示变更回调
 * @param {Function} props.onAddTag - 添加标签回调
 * @param {Function} props.onRemoveTag - 移除标签回调
 * @param {Function} props.onTagChange - 标签变更回调
 * @param {Function} props.onFlagChange - Flag变更回调
 */
function AdminContestChallengeModal({
  isOpen = false,
  mode = 'add',
  contestId,
  challengeId,
  challenge = {
    name: '',
    description: '',
    attempt: 0,
    hidden: false,
    hints: [],
    tags: [],
  },
  flags = [
    {
      blood: [0, 0, 0],
      current_score: 1000,
      decay: 100,
      id: 1,
      last: '0001-01-01T00:00:00Z',
      min_score: 100,
      score: 1000,
      score_type: 0,
      solvers: 0,
      value: '',
    },
  ],
  onClose,
  onSubmit,
  onChange,
  onAddHint,
  onRemoveHint,
  onHintChange,
  onAddTag,
  onRemoveTag,
  onTagChange,
  onFlagChange,
}) {
  const { t } = useTranslation();
  // 用于本地编辑的flags状态
  const [editingFlags, setEditingFlags] = useState([...flags]);
  const [solversModal, setSolversModal] = useState({ open: false, flagId: null, flagIndex: 0 });

  // 本地更新flag, 不触发onFlagChange
  const handleFlagChange = (index, updatedFlag) => {
    const newFlags = [...editingFlags];
    newFlags[index] = updatedFlag;
    setEditingFlags(newFlags);
  };

  useEffect(() => {
    setEditingFlags(flags);
  }, [flags]);

  // 提交单个flag的更新
  const handleUpdateFlag = (index) => {
    onFlagChange(index, editingFlags[index]);
  };

  // 统一输入框样式
  const inputBaseClass =
    'w-full bg-black/20 border border-neutral-300/30 rounded-md p-3 text-neutral-50 font-mono focus:border-geek-400 focus:outline-none transition-colors duration-200';
  const textareaClass = `${inputBaseClass} min-h-[100px] resize-none`;

  // 禁用输入框样式
  const disabledInputClass = `${inputBaseClass} opacity-70 cursor-not-allowed`;

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <motion.div
        className="w-full max-w-4xl bg-neutral-900 border border-neutral-300 rounded-md overflow-hidden"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.2 }}
      >
        {/* 标题栏 */}
        <div className="flex justify-between items-center p-4 border-b border-neutral-700">
          <h2 className="text-xl font-mono text-neutral-50">
            {mode === 'add'
              ? t('admin.contests.challengeModal.titleAdd')
              : t('admin.contests.challengeModal.titleEdit')}
          </h2>
          <Button
            variant="ghost"
            size="icon"
            className="!text-neutral-400 hover:!text-neutral-300"
            aria-label={t('common.close')}
            onClick={onClose}
          >
            <IconX size={18} />
          </Button>
        </div>

        {/* 表单内容 */}
        <div className="p-6 max-h-[70vh] overflow-y-auto">
          <div className="space-y-4">
            {/* 基本信息 */}
            <div className="mb-4">
              <h3 className="text-lg font-mono text-neutral-50 mb-3">
                {t('admin.contests.challengeModal.sections.basic')}
              </h3>

              <div className="space-y-3">
                <div>
                  <label className="block text-sm font-mono text-neutral-400 mb-1">
                    {t('admin.contests.challengeModal.labels.name')}
                  </label>
                  <input
                    type="text"
                    value={challenge.name}
                    onChange={(e) => onChange({ ...challenge, name: e.target.value })}
                    className={inputBaseClass}
                    placeholder={t('admin.contests.challengeModal.placeholders.name')}
                  />
                </div>

                <div>
                  <label className="block text-sm font-mono text-neutral-400 mb-1">
                    {t('admin.contests.challengeModal.labels.description')}
                  </label>
                  <textarea
                    value={challenge.description}
                    onChange={(e) => onChange({ ...challenge, description: e.target.value })}
                    className={textareaClass}
                    placeholder={t('admin.contests.challengeModal.placeholders.description')}
                  />
                </div>

                <div className="flex gap-6">
                  <div className="w-1/3">
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.contests.challengeModal.labels.attempts')}
                    </label>
                    <input
                      type="number"
                      value={challenge.attempt}
                      onChange={(e) => onChange({ ...challenge, attempt: parseInt(e.target.value) || 0 })}
                      className={inputBaseClass}
                      placeholder="0"
                    />
                  </div>

                  <div className="w-1/3 flex items-center pt-6">
                    <input
                      type="checkbox"
                      id="hidden-checkbox"
                      checked={challenge.hidden}
                      onChange={(e) => onChange({ ...challenge, hidden: e.target.checked })}
                      className="mr-2 h-4 w-4 accent-geek-400"
                    />
                    <label htmlFor="hidden-checkbox" className="text-sm font-mono text-neutral-400">
                      {t('admin.contests.challengeModal.labels.hidden')}
                    </label>
                  </div>
                </div>
              </div>
            </div>

            {/* Flag设置 */}
            <div className="border-t border-neutral-700 pt-4">
              <div className="flex justify-between items-center mb-3">
                <h3 className="text-lg font-mono text-neutral-50">
                  {t('admin.contests.challengeModal.sections.flags')}
                </h3>
              </div>

              {editingFlags.map((flag, index) => (
                <div key={index} className="mb-6 p-4 border border-neutral-700 rounded-md bg-black/20">
                  <div className="flex justify-between mb-2">
                    <h4 className="text-md font-mono text-geek-400">Flag #{index + 1}</h4>
                    <div className="flex gap-1">
                      {mode === 'edit' && flag.id && (
                        <Button
                          variant="ghost"
                          size="sm"
                          align="icon-left"
                          icon={<IconUsers size={15} />}
                          className="!text-neutral-300 hover:!text-neutral-100 !bg-neutral-700/40 hover:!bg-neutral-700/70 !text-xs"
                          onClick={() => setSolversModal({ open: true, flagId: flag.id, flagIndex: index })}
                        >
                          {t('admin.contests.challengeModal.actions.viewSolvers')}
                        </Button>
                      )}
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.contests.challengeModal.labels.flagValueHint')}
                      </label>
                      <input
                        type="text"
                        value={flag.value}
                        onChange={(e) => handleFlagChange(index, { ...flag, value: e.target.value })}
                        className={inputBaseClass}
                        placeholder={t('admin.contests.challengeModal.placeholders.flagValue')}
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.contests.challengeModal.labels.scoreCurve')}
                      </label>
                      <select
                        value={flag.score_type}
                        onChange={(e) => handleFlagChange(index, { ...flag, score_type: parseInt(e.target.value) })}
                        className={inputBaseClass}
                      >
                        <option value={0}>{t('admin.contests.challengeModal.scoreCurve.static')}</option>
                        <option value={1}>{t('admin.contests.challengeModal.scoreCurve.linear')}</option>
                        <option value={2}>{t('admin.contests.challengeModal.scoreCurve.log')}</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.contests.challengeModal.labels.initialScore')}
                      </label>
                      <input
                        type="number"
                        value={flag.score}
                        onChange={(e) => handleFlagChange(index, { ...flag, score: parseFloat(e.target.value) || 0 })}
                        className={inputBaseClass}
                        placeholder="1000"
                        step="0.1"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.contests.challengeModal.labels.currentScore')}
                      </label>
                      <input
                        type="number"
                        value={flag.current_score}
                        className={disabledInputClass}
                        placeholder="1000"
                        disabled
                      />
                    </div>

                    {flag.score_type !== 0 && (
                      <>
                        <div>
                          <label className="block text-sm font-mono text-neutral-400 mb-1">
                            {t('admin.contests.challengeModal.labels.decay')}
                          </label>
                          <input
                            type="number"
                            value={flag.decay}
                            onChange={(e) =>
                              handleFlagChange(index, { ...flag, decay: parseFloat(e.target.value) || 0 })
                            }
                            className={inputBaseClass}
                            placeholder="100"
                            step="0.1"
                          />
                        </div>

                        <div>
                          <label className="block text-sm font-mono text-neutral-400 mb-1">
                            {t('admin.contests.challengeModal.labels.minScore')}
                          </label>
                          <input
                            type="number"
                            value={flag.min_score}
                            onChange={(e) =>
                              handleFlagChange(index, { ...flag, min_score: parseFloat(e.target.value) || 0 })
                            }
                            className={inputBaseClass}
                            placeholder="100"
                            step="0.1"
                          />
                        </div>
                      </>
                    )}

                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.contests.challengeModal.labels.solvers')}
                      </label>
                      <input
                        type="number"
                        value={flag.solvers}
                        className={disabledInputClass}
                        placeholder="0"
                        disabled
                      />
                    </div>
                  </div>

                  {/* 分数曲线预览（拖拽控制点可调整参数） */}
                  <ScoreCurveChart
                    scoreType={flag.score_type}
                    score={flag.score}
                    decay={flag.decay}
                    minScore={flag.min_score}
                    onChange={(patch) => handleFlagChange(index, { ...flag, ...patch })}
                  />
                </div>
              ))}
            </div>

            {/* 标签设置 */}
            <div className="border-t border-neutral-700 pt-4">
              <div className="flex justify-between items-center mb-3">
                <h3 className="text-lg font-mono text-neutral-50">
                  {t('admin.contests.challengeModal.sections.tags')}
                </h3>
                <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onAddTag}>
                  {t('admin.contests.challengeModal.actions.addTag')}
                </Button>
              </div>

              <div className="space-y-2">
                {challenge.tags.map((tag, index) => (
                  <div key={index} className="flex gap-2 items-center">
                    <input
                      type="text"
                      value={tag}
                      onChange={(e) => onTagChange(index, e.target.value)}
                      className={inputBaseClass}
                      placeholder={t('admin.contests.challengeModal.placeholders.tag', { index: index + 1 })}
                    />
                    {
                      <Button
                        variant="ghost"
                        size="icon"
                        className="!text-red-400 hover:!text-red-300"
                        onClick={() => onRemoveTag(index)}
                      >
                        <IconTrash size={18} />
                      </Button>
                    }
                  </div>
                ))}
              </div>
            </div>

            {/* 提示设置 */}
            <div className="border-t border-neutral-700 pt-4">
              <div className="flex justify-between items-center mb-3">
                <h3 className="text-lg font-mono text-neutral-50">
                  {t('admin.contests.challengeModal.sections.hints')}
                </h3>
                <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onAddHint}>
                  {t('admin.contests.challengeModal.actions.addHint')}
                </Button>
              </div>

              <div className="space-y-2">
                {challenge.hints.map((hint, index) => (
                  <div key={index} className="flex gap-2 items-center">
                    <input
                      type="text"
                      value={hint}
                      onChange={(e) => onHintChange(index, e.target.value)}
                      className={inputBaseClass}
                      placeholder={t('admin.contests.challengeModal.placeholders.hint', { index: index + 1 })}
                    />
                    {
                      <Button
                        variant="ghost"
                        size="icon"
                        className="!text-red-400 hover:!text-red-300"
                        onClick={() => onRemoveHint(index)}
                      >
                        <IconTrash size={18} />
                      </Button>
                    }
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* 底部按钮 */}
        <div className="flex justify-end gap-4 p-4 border-t border-neutral-700">
          <Button variant="ghost" onClick={onClose}>
            {t('common.cancel')}
          </Button>
          <Button
            variant="primary"
            onClick={() => {
              editingFlags.map((_, index) => handleUpdateFlag(index));
              onSubmit(challenge, flags);
            }}
          >
            {mode === 'add' ? t('admin.contests.challengeModal.actions.add') : t('common.saveChanges')}
          </Button>
        </div>
      </motion.div>
      <FlagSolversModal
        isOpen={solversModal.open}
        onClose={() => setSolversModal({ ...solversModal, open: false })}
        flagIndex={solversModal.flagIndex}
        contestId={contestId}
        challengeId={challengeId}
        flagId={solversModal.flagId}
      />
    </div>
  );
}

export default AdminContestChallengeModal;
