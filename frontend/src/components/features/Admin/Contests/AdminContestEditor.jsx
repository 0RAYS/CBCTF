import { motion } from 'motion/react';
import { useState, useRef } from 'react';
import {
  IconEdit,
  IconTrash,
  IconDeviceFloppy,
  IconPlus,
  IconX,
  IconUpload,
  IconPhoto,
  IconExclamationCircle,
} from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, DateTimeInput } from '../../../../components/common';

/**
 * 比赛编辑组件
 * @param {Object} props
 * @param {Object} props.contest - 比赛详情对象
 * @param {Function} props.onSave - 保存回调函数
 * @param {Function} props.onCancel - 取消回调函数
 * @param {Function} props.onImageUpload - 图片上传回调函数
 */
function AdminContestEditor({ contest: initialContest, onSave, onCancel, onImageUpload }) {
  const { t } = useTranslation();

  const [contest, setContest] = useState(
    initialContest || {
      title: '',
      description: '',
      image: '',
      status: 'upcoming',
      startTime: '',
      endTime: '',
      participants: 0,
      rules: [],
      prizes: [],
      timeline: [],
      prefix: 'CBCTF',
      size: 4,
      hidden: false,
      captcha: '',
      blood: true,
      victims: 1,
    }
  );

  const [newRule, setNewRule] = useState('');
  const [editingRule, setEditingRule] = useState('');
  const [editingIndex, setEditingIndex] = useState({
    rule: null,
    prize: null,
    timeline: null,
  });

  // 添加表单验证状态
  const [validationErrors, setValidationErrors] = useState({});
  // 添加预览图像状态
  const [previewImage, setPreviewImage] = useState(contest.image || '');
  // 添加文件上传引用
  const fileInputRef = useRef(null);

  // 统一所有输入框样式
  const inputBaseClass =
    'w-full bg-black/70 border-2 border-neutral-600 rounded-md p-3 text-neutral-50 font-mono focus:border-geek-400 focus:outline-none transition-colors duration-200';
  const textareaClass = `${inputBaseClass} min-h-[100px] resize-none`;
  const selectClass = 'select-custom select-custom-md';

  const getPrizeRankLabel = (index) => {
    const rank = index + 1;
    if (rank === 1) return t('admin.contests.editor.prizeRank.first');
    if (rank === 2) return t('admin.contests.editor.prizeRank.second');
    if (rank === 3) return t('admin.contests.editor.prizeRank.third');
    return t('admin.contests.editor.prizeRank.other', { rank });
  };

  // 改进日期转换函数
  const formatDateForInput = (dateString) => {
    if (!dateString) return '';

    try {
      // 尝试创建日期对象
      const date = new Date(dateString);

      // 检查日期是否有效
      if (isNaN(date.getTime())) return '';

      // 格式化为YYYY-MM-DDThh:mm
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');

      return `${year}-${month}-${day}T${hours}:${minutes}`;
    } catch (error) {
      return error.message;
    }
  };

  // 处理图片上传
  const handleImageUpload = (e) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // 如果传入了onImageUpload回调，则调用它
    if (onImageUpload) {
      onImageUpload(file).then((imageUrl) => {
        if (imageUrl) {
          setContest((prev) => ({ ...prev, image: imageUrl }));
          setPreviewImage(imageUrl);
        }
      });
    } else {
      // 如果没有回调，则创建本地预览
      const reader = new FileReader();
      reader.onload = (e) => {
        const imageUrl = e.target.result;
        setContest((prev) => ({ ...prev, image: imageUrl }));
        setPreviewImage(imageUrl);
      };
      reader.readAsDataURL(file);
    }

    // 清空文件输入框，以便再次选择同一文件时触发change事件
    e.target.value = '';
  };

  // 触发文件选择
  const triggerFileInput = () => {
    fileInputRef.current?.click();
  };

  // 修改基本信息
  const handleChange = (e) => {
    const { name, value } = e.target;

    // 特殊处理日期时间字段
    if (name === 'startTime' || name === 'endTime') {
      if (value) {
        try {
          // 创建日期对象并验证其有效性
          const date = new Date(value);

          // 检查日期是否有效
          if (!isNaN(date.getTime())) {
            setContest((prev) => ({ ...prev, [name]: date.toISOString() }));

            // 清除该字段的验证错误
            if (validationErrors[name]) {
              setValidationErrors((prev) => ({ ...prev, [name]: '' }));
            }
          } else {
            // 设置验证错误
            setValidationErrors((prev) => ({ ...prev, [name]: t('admin.contests.editor.validation.invalidDate') }));
          }
        } catch (error) {
          // 设置验证错误
          setValidationErrors((prev) => ({ ...prev, [name]: error.message }));
        }
      } else {
        setContest((prev) => ({ ...prev, [name]: '' }));
      }
    } else {
      setContest((prev) => ({ ...prev, [name]: value }));

      // 如果该字段有验证错误，则在用户输入后清除
      if (validationErrors[name]) {
        setValidationErrors((prev) => ({ ...prev, [name]: '' }));
      }
    }
  };

  // 添加规则
  const addRule = () => {
    if (!newRule.trim()) return;
    setContest((prev) => ({
      ...prev,
      rules: [...(prev.rules || []), newRule],
    }));
    setNewRule('');
  };

  // 更新规则
  const updateRule = (index, value) => {
    const updatedRules = [...contest.rules];
    updatedRules[index] = value;
    setContest((prev) => ({ ...prev, rules: updatedRules }));
    setEditingIndex((prev) => ({ ...prev, rule: null }));
    setEditingRule('');
  };

  // 删除规则
  const deleteRule = (index) => {
    const updatedRules = contest.rules.filter((_, i) => i !== index);
    setContest((prev) => ({ ...prev, rules: updatedRules }));
  };

  // 添加奖励
  const addPrize = () => {
    setContest((prev) => ({
      ...prev,
      prizes: [...(prev.prizes || []), { amount: '$0', description: '' }],
    }));
  };

  // 更新奖励
  const updatePrize = (index, field, value) => {
    const updatedPrizes = [...contest.prizes];
    updatedPrizes[index] = { ...updatedPrizes[index], [field]: value };
    setContest((prev) => ({ ...prev, prizes: updatedPrizes }));
  };

  // 删除奖励
  const deletePrize = (index) => {
    const updatedPrizes = contest.prizes.filter((_, i) => i !== index);
    setContest((prev) => ({ ...prev, prizes: updatedPrizes }));
  };

  // 添加时间线事件
  const addTimelineEvent = () => {
    setContest((prev) => ({
      ...prev,
      timeline: [...(prev.timeline || []), { date: '', title: '', description: '' }],
    }));
  };

  // 更新时间线
  const updateTimeline = (index, field, value) => {
    const updatedTimeline = [...contest.timeline];
    updatedTimeline[index] = { ...updatedTimeline[index], [field]: value };
    setContest((prev) => ({ ...prev, timeline: updatedTimeline }));
  };

  // 删除时间线事件
  const deleteTimelineEvent = (index) => {
    const updatedTimeline = contest.timeline.filter((_, i) => i !== index);
    setContest((prev) => ({ ...prev, timeline: updatedTimeline }));
  };

  // 验证表单
  const validateForm = () => {
    const errors = {};

    if (!contest.title.trim()) {
      errors.title = t('admin.contests.editor.validation.titleRequired');
    }

    if (contest.startTime && contest.endTime && new Date(contest.startTime) >= new Date(contest.endTime)) {
      errors.endTime = t('admin.contests.editor.validation.endTimeAfterStart');
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // 保存所有修改
  const handleSubmit = (e) => {
    e.preventDefault();

    if (validateForm()) {
      onSave(contest);
    } else {
      // 滚动到第一个错误字段
      const firstErrorField = Object.keys(validationErrors)[0];
      if (firstErrorField) {
        document.getElementsByName(firstErrorField)[0]?.scrollIntoView({
          behavior: 'smooth',
          block: 'center',
        });
      }
    }
  };

  return (
    <form onSubmit={handleSubmit} className="w-full mx-auto space-y-6">
      {/* 头部信息编辑区域 */}
      <motion.div
        className="relative w-full border-2 border-neutral-300 rounded-md overflow-hidden bg-neutral-900 p-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <div className="mb-6">
          <label className="block text-neutral-300 font-medium mb-2">{t('admin.contests.editor.labels.title')}</label>
          <input
            type="text"
            name="title"
            value={contest.title}
            onChange={handleChange}
            className={`${inputBaseClass} ${validationErrors.title ? 'border-red-400' : ''}`}
            required
          />
          {validationErrors.title && (
            <p className="mt-2 text-red-400 text-sm flex items-center gap-1">
              <IconExclamationCircle size={16} />
              {validationErrors.title}
            </p>
          )}
        </div>

        <div className="mb-6">
          <label className="block text-neutral-300 font-medium mb-2">
            {t('admin.contests.editor.labels.description')}
          </label>
          <textarea
            name="description"
            value={contest.description}
            onChange={handleChange}
            className={`${textareaClass} ${validationErrors.description ? 'border-red-400' : ''}`}
          />
          {validationErrors.description && (
            <p className="mt-2 text-red-400 text-sm flex items-center gap-1">
              <IconExclamationCircle size={16} />
              {validationErrors.description}
            </p>
          )}
        </div>

        <div className="mb-6">
          <label className="block text-neutral-300 font-medium mb-2">
            {t('admin.contests.editor.labels.backgroundImage')}
          </label>
          <div className="flex gap-4">
            <div className="flex-1" onClick={triggerFileInput}>
              {/* 图片预览区域 */}
              <div className="w-full h-[200px] border-2 border-dashed border-neutral-600 rounded-md bg-black/30 flex items-center justify-center overflow-hidden relative group">
                {previewImage ? (
                  <>
                    <img
                      src={previewImage}
                      alt={t('admin.contests.editor.labels.backgroundAlt')}
                      className="w-full h-full object-cover"
                    />
                    <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <Button variant="primary" size="sm" align="icon-left" icon={<IconUpload size={16} />}>
                        {t('admin.contests.editor.actions.replaceImage')}
                      </Button>
                    </div>
                  </>
                ) : (
                  <Button
                    variant="ghost"
                    size="sm"
                    align="icon-left"
                    icon={<IconPhoto size={48} />}
                    className="flex-col !gap-3 !text-neutral-400"
                  >
                    {t('admin.contests.editor.actions.uploadImage')}
                  </Button>
                )}
              </div>

              {/* 隐藏的文件输入 */}
              <input ref={fileInputRef} type="file" accept="image/*" className="hidden" onChange={handleImageUpload} />
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          <div>
            <label className="block text-neutral-300 font-medium mb-2">
              {t('admin.contests.editor.labels.status')}
            </label>
            <select name="status" value={contest.status} className={selectClass} disabled>
              <option value="upcoming">{t('admin.contests.editor.status.upcoming')}</option>
              <option value="active">{t('admin.contests.editor.status.active')}</option>
              <option value="ended">{t('admin.contests.editor.status.ended')}</option>
            </select>
          </div>

          <div>
            <label className="block text-neutral-300 font-medium mb-2">
              {t('admin.contests.editor.labels.startTime')}
            </label>
            <DateTimeInput
              name="startTime"
              value={formatDateForInput(contest.startTime)}
              onChange={handleChange}
              error={validationErrors.startTime}
            />
          </div>

          <div>
            <label className="block text-neutral-300 font-medium mb-2">
              {t('admin.contests.editor.labels.endTime')}
            </label>
            <DateTimeInput
              name="endTime"
              value={formatDateForInput(contest.endTime)}
              onChange={handleChange}
              error={validationErrors.endTime}
            />
          </div>
        </div>

        {/* 高级设置 */}
        <div className="mb-4">
          <h3 className="text-xl font-mono text-neutral-50 tracking-wider mb-4">
            {t('admin.contests.editor.sections.advanced')}
          </h3>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
            <div>
              <label className="block text-neutral-300 font-medium mb-2">
                {t('admin.contests.editor.labels.flagPrefix')}
              </label>
              <input
                type="text"
                name="prefix"
                value={contest.prefix}
                onChange={(e) =>
                  setContest((prev) => ({
                    ...prev,
                    prefix: e.target.value,
                  }))
                }
                className={inputBaseClass}
                placeholder={t('admin.contests.editor.placeholders.flagPrefix')}
              />
              <p className="mt-1 text-neutral-500 text-sm">{t('admin.contests.editor.help.flagPrefix')}</p>
            </div>

            <div>
              <label className="block text-neutral-300 font-medium mb-2">
                {t('admin.contests.editor.labels.teamSize')}
              </label>
              <input
                type="number"
                name="size"
                value={contest.size}
                onChange={(e) =>
                  setContest((prev) => ({
                    ...prev,
                    size: parseInt(e.target.value) || 1,
                  }))
                }
                min="1"
                max="10"
                className={inputBaseClass}
                placeholder={t('admin.contests.editor.placeholders.teamSize')}
              />
              <p className="mt-1 text-neutral-500 text-sm">{t('admin.contests.editor.help.teamSize')}</p>
            </div>

            <div>
              <label className="block text-neutral-300 font-medium mb-2">
                {t('admin.contests.editor.labels.teamVictims')}
              </label>
              <input
                type="number"
                name="size"
                value={contest.victims}
                onChange={(e) =>
                  setContest((prev) => ({
                    ...prev,
                    victims: parseInt(e.target.value) || 1,
                  }))
                }
                min="1"
                max="10"
                className={inputBaseClass}
                placeholder={t('admin.contests.editor.placeholders.teamVictims')}
              />
              <p className="mt-1 text-neutral-500 text-sm">{t('admin.contests.editor.help.teamVictims')}</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-neutral-300 font-medium mb-2">
                {t('admin.contests.editor.labels.captcha')}
              </label>
              <input
                type="text"
                name="captcha"
                value={contest.captcha}
                onChange={(e) =>
                  setContest((prev) => ({
                    ...prev,
                    captcha: e.target.value,
                  }))
                }
                className={inputBaseClass}
                placeholder={t('admin.contests.editor.placeholders.captcha')}
              />
              <p className="mt-1 text-neutral-500 text-sm">{t('admin.contests.editor.help.captcha')}</p>
            </div>

            <div className="flex flex-col gap-4">
              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  id="hidden-toggle"
                  className="w-5 h-5 bg-black border-2 border-neutral-600 rounded text-geek-400 focus:ring-geek-400"
                  checked={contest.hidden}
                  onChange={(e) =>
                    setContest((prev) => ({
                      ...prev,
                      hidden: e.target.checked,
                    }))
                  }
                />
                <label htmlFor="hidden-toggle" className="text-neutral-300 font-medium">
                  {t('admin.contests.editor.labels.hiddenContest')}
                </label>
              </div>
              <p className="text-neutral-500 text-sm">{t('admin.contests.editor.help.hiddenContest')}</p>

              <div className="flex items-center gap-3 mt-2">
                <input
                  type="checkbox"
                  id="blood-toggle"
                  className="w-5 h-5 bg-black border-2 border-neutral-600 rounded text-geek-400 focus:ring-geek-400"
                  checked={contest.blood}
                  onChange={(e) =>
                    setContest((prev) => ({
                      ...prev,
                      blood: e.target.checked,
                    }))
                  }
                />
                <label htmlFor="blood-toggle" className="text-neutral-300 font-medium">
                  {t('admin.contests.editor.labels.bloodBonus')}
                </label>
              </div>
              <p className="text-neutral-500 text-sm">{t('admin.contests.editor.help.bloodBonus')}</p>
            </div>
          </div>
        </div>
      </motion.div>

      {/* 详细信息编辑区域 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* 比赛规则编辑 */}
        <motion.div className="col-span-1 md:col-span-2 border-2 border-neutral-300 rounded-md bg-neutral-900 p-8">
          <h2 className="text-2xl font-mono text-neutral-50 tracking-wider mb-6">
            {t('admin.contests.editor.sections.rules')}
          </h2>

          <div className="flex gap-2 mb-6">
            <input
              type="text"
              value={newRule}
              onChange={(e) => setNewRule(e.target.value)}
              className={inputBaseClass}
              placeholder={t('admin.contests.editor.placeholders.newRule')}
            />
            <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={18} />} onClick={addRule}>
              {t('admin.contests.editor.actions.addRule')}
            </Button>
          </div>

          <div className="space-y-4 text-neutral-300">
            {contest.rules?.map((rule, index) => (
              <motion.div
                key={index}
                className="flex items-center gap-4 p-2 rounded-md hover:bg-neutral-300/5 transition-colors duration-200"
              >
                <span className="text-geek-400 font-mono">{(index + 1).toString().padStart(2, '0')}</span>

                {editingIndex.rule === index ? (
                  <input
                    type="text"
                    value={editingRule}
                    onChange={(e) => setEditingRule(e.target.value)}
                    className={inputBaseClass}
                    autoFocus
                    onFocus={() => setEditingRule(rule)}
                    onBlur={() => updateRule(index, editingRule)}
                  />
                ) : (
                  <p className="flex-1">{rule}</p>
                )}

                <div className="flex gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="!bg-transparent !text-geek-400 hover:!text-geek-300"
                    onClick={() => setEditingIndex((prev) => ({ ...prev, rule: index }))}
                  >
                    <IconEdit size={18} />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="!bg-transparent !text-red-400 hover:!text-red-300"
                    onClick={() => deleteRule(index)}
                  >
                    <IconTrash size={18} />
                  </Button>
                </div>
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* 奖励信息编辑 */}
        <motion.div className="border-2 border-neutral-300 rounded-md bg-neutral-900 p-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-mono text-neutral-50 tracking-wider">奖励设置</h2>
            <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={addPrize}>
              {t('admin.contests.editor.actions.addPrize')}
            </Button>
          </div>

          <div className="space-y-6">
            {contest.prizes?.map((prize, index) => (
              <motion.div key={index} className="p-4 border border-neutral-700 rounded-md bg-black/30">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <span
                      className={`text-lg font-mono
                      ${
                        index === 0
                          ? 'text-yellow-400'
                          : index === 1
                            ? 'text-neutral-300'
                            : index === 2
                              ? 'text-orange-400'
                              : 'text-neutral-400'
                      }
                    `}
                    >
                      {getPrizeRankLabel(index)}
                    </span>
                  </div>

                  <Button
                    variant="ghost"
                    size="icon"
                    className="!bg-transparent !text-red-400 hover:!text-red-300"
                    onClick={() => deletePrize(index)}
                  >
                    <IconTrash size={18} />
                  </Button>
                </div>

                <div className="space-y-4">
                  <div>
                    <label className="block text-neutral-300 font-medium mb-2 text-sm">
                      {t('admin.contests.editor.labels.prizeAmount')}
                    </label>
                    <input
                      type="text"
                      value={prize.amount}
                      onChange={(e) => updatePrize(index, 'amount', e.target.value)}
                      className="w-full bg-black/70 border-2 border-neutral-600 rounded-md p-2 text-neutral-50 font-mono focus:border-geek-400 focus:outline-none transition-colors duration-200"
                    />
                  </div>

                  <div>
                    <label className="block text-neutral-300 font-medium mb-2 text-sm">
                      {t('admin.contests.editor.labels.prizeDescription')}
                    </label>
                    <textarea
                      value={prize.description || ''}
                      onChange={(e) => updatePrize(index, 'description', e.target.value)}
                      className="w-full bg-black/70 border-2 border-neutral-600 rounded-md p-2 text-neutral-50 text-sm min-h-[60px] resize-none focus:border-geek-400 focus:outline-none transition-colors duration-200"
                    />
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </div>

      {/* 时间线编辑 */}
      <motion.div className="border-2 border-neutral-300 rounded-md bg-neutral-900 p-8">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-mono text-neutral-50 tracking-wider">
            {t('admin.contests.editor.sections.timeline')}
          </h2>
          <Button
            variant="primary"
            size="sm"
            align="icon-left"
            icon={<IconPlus size={16} />}
            onClick={addTimelineEvent}
          >
            {t('admin.contests.editor.actions.addTimeline')}
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {contest.timeline?.map((item, index) => (
            <motion.div
              key={index}
              className="p-4 border border-neutral-700 rounded-md bg-black/30 relative"
              whileHover={{ y: -5 }}
            >
              <Button
                variant="ghost"
                size="icon"
                className="!bg-transparent !text-red-400 hover:!text-red-300 absolute top-2 right-2"
                onClick={() => deleteTimelineEvent(index)}
              >
                <IconX size={18} />
              </Button>

              <div className="space-y-4 mt-4">
                <div>
                  <label className="block text-neutral-300 font-medium mb-2 text-sm">
                    {t('admin.contests.editor.labels.timelineDate')}
                  </label>
                  <DateTimeInput
                    value={formatDateForInput(item.date)}
                    onChange={(e) => updateTimeline(index, 'date', e.target.value)}
                  />
                </div>

                <div>
                  <label className="block text-neutral-300 font-medium mb-2 text-sm">
                    {t('admin.contests.editor.labels.timelineTitle')}
                  </label>
                  <input
                    type="text"
                    value={item.title}
                    onChange={(e) => updateTimeline(index, 'title', e.target.value)}
                    className="w-full bg-black/70 border-2 border-neutral-600 rounded-md p-2 text-neutral-50 focus:border-geek-400 focus:outline-none transition-colors duration-200"
                  />
                </div>

                <div>
                  <label className="block text-neutral-300 font-medium mb-2 text-sm">
                    {t('admin.contests.editor.labels.timelineDescription')}
                  </label>
                  <textarea
                    value={item.description}
                    onChange={(e) => updateTimeline(index, 'description', e.target.value)}
                    className="w-full bg-black/70 border-2 border-neutral-600 rounded-md p-2 text-neutral-50 text-sm min-h-[60px] resize-none focus:border-geek-400 focus:outline-none transition-colors duration-200"
                  />
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </motion.div>

      {/* 操作按钮 */}
      <div className="flex justify-end gap-4 mt-8">
        <Button variant="ghost" size="sm" onClick={onCancel}>
          {t('common.cancel')}
        </Button>
        <Button variant="primary" size="sm" align="icon-left" icon={<IconDeviceFloppy size={18} />} type="submit">
          {t('common.saveChanges')}
        </Button>
      </div>
    </form>
  );
}

export default AdminContestEditor;
