import { IconPlus, IconTrash } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import Button from '@/components/common/Button';
import Checkbox from '@/components/common/Checkbox';
import Modal from '@/components/common/Modal';
import ComposeEditor from '@/components/features/Admin/challenges/editor/ComposeEditor.jsx';
import useComposeGuide from '@/components/features/Admin/challenges/editor/useComposeGuide.js';
import ChallengeFlagPanel from '@/components/features/Admin/challenges/editor/ChallengeFlagPanel.jsx';
import {
  inputBaseClass,
  selectClass,
  textareaClass,
} from '@/components/features/Admin/challenges/editor/editorStyles.js';

export default function ChallengeEditorDialog({
  isOpen = false,
  mode = 'add',
  challenge = {
    name: '',
    description: '',
    category: '',
    type: 'static',
    contest_id: 0,
    hidden: false,
    attempt: 0,
    hints: [],
    tags: [],
    generator_image: '',
    flags: [{ id: 0, value: '' }],
    docker_compose: '',
    network_policies: [],
  },
  categories = [],
  contests = [],
  contestLocked = false,
  flagVersion = 0,
  onClose,
  onSubmit,
  onChange,
  onAddFlag,
  onRemoveFlag,
  onFlagChange,
  onFlagsChanged,
}) {
  const { t } = useTranslation();
  const compose = useComposeGuide({ isOpen, challenge, onChange, t });
  const isEditMode = mode === 'edit';
  const invalidCompose =
    challenge.type === 'pods' && (compose.guideValidation.list.length > 0 || compose.rawValidation.list.length > 0);
  const podNoticeLines = [
    t('admin.challengeModal.podsNotice.flagFormat', { format: '`static{}`, `leet{}`, `uuid{}`' }),
    t('admin.challengeModal.podsNotice.flagPrefix'),
    t('admin.challengeModal.podsNotice.flagVolume'),
    '',
    t('admin.challengeModal.podsNotice.services'),
    '',
    t('admin.challengeModal.podsNotice.networks'),
    t('admin.challengeModal.podsNotice.networksWithIp'),
    t('admin.challengeModal.podsNotice.networksShared'),
    t('admin.challengeModal.podsNotice.portUnique'),
    '',
    t('admin.challengeModal.podsNotice.exampleIndent'),
  ];
  // 归属比赛在创建后不可更改，跨比赛复用请使用复制。
  const changeContest = (value) => {
    const contestId = Number(value);
    onChange({
      ...challenge,
      contest_id: contestId,
      hidden: isEditMode ? challenge.hidden : contestId !== 0,
    });
  };
  const contestName = contests.find((contest) => contest.id === Number(challenge.contest_id || 0))?.name;
  const contestLabel =
    contestName || (challenge.contest_id ? `#${challenge.contest_id}` : t('admin.challengeModal.contest.unassigned'));

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size={isEditMode ? 'full' : '2xl'}
      title={t(`admin.challengeModal.title.${mode === 'add' ? 'add' : mode === 'edit' ? 'edit' : 'delete'}`)}
      bodyClassName="p-3 sm:p-4 lg:p-5 [&_*]:min-w-0"
      footer={
        <>
          <Button variant="ghost" size="sm" onClick={onClose}>
            {t('common.cancel')}
          </Button>
          <Button variant="primary" size="sm" disabled={invalidCompose} onClick={() => onSubmit(challenge)}>
            {mode === 'add'
              ? t('admin.challengeModal.actions.add')
              : mode === 'edit'
                ? t('common.saveChanges')
                : t('admin.challengeModal.actions.confirmDelete')}
          </Button>
        </>
      }
    >
      {mode === 'delete' ? (
        <div className="text-center py-12">
          <div className="mb-8">
            <div className="mx-auto w-20 h-20 bg-red-400/20 rounded-full flex items-center justify-center mb-6">
              <IconTrash size={40} className="text-red-400" />
            </div>
            <h3 className="text-lg font-mono text-neutral-50 mb-4">{t('admin.challengeModal.delete.title')}</h3>
            <p className="text-neutral-400 font-mono text-lg mb-2">{t('admin.challengeModal.delete.prompt')}</p>
            <p className="text-red-400 font-mono text-sm">{t('admin.challengeModal.delete.warning')}</p>
          </div>
        </div>
      ) : (
        <div className={isEditMode ? 'grid grid-cols-1 xl:grid-cols-2 gap-4' : 'space-y-3'}>
          <div className="space-y-3">
            <div className="rounded-md border border-neutral-700/70 bg-black/10 p-3">
              <h3 className="text-lg font-mono text-neutral-50 mb-3">{t('admin.challengeModal.sections.basic')}</h3>
              <div className="space-y-3">
                <div className="grid grid-cols-1 md:grid-cols-4 gap-3 lg:gap-4">
                  <div className="md:col-span-2">
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.name')}
                    </label>
                    <input
                      type="text"
                      value={challenge.name}
                      onChange={(e) => onChange({ ...challenge, name: e.target.value })}
                      className={inputBaseClass}
                      placeholder={t('admin.challengeModal.placeholders.name')}
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.category')}
                    </label>
                    <input
                      type="text"
                      value={challenge.category}
                      onChange={(e) => onChange({ ...challenge, category: e.target.value })}
                      className={selectClass}
                      placeholder={t('admin.challengeModal.placeholders.category')}
                      list="category-options"
                    />
                    <datalist id="category-options">
                      {categories.map((category) => (
                        <option key={category} value={category} />
                      ))}
                    </datalist>
                  </div>
                  <div>
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.type')}
                    </label>
                    <select
                      value={challenge.type}
                      onChange={(e) => onChange({ ...challenge, type: e.target.value })}
                      className={selectClass}
                      required
                    >
                      <option value="static">{t('common.challengeTypes.static')}</option>
                      <option value="dynamic">{t('common.challengeTypes.dynamic')}</option>
                      <option value="pods">{t('common.challengeTypes.pods')}</option>
                    </select>
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-mono text-neutral-400 mb-1">
                    {t('admin.challengeModal.labels.description')}
                  </label>
                  <textarea
                    value={challenge.description}
                    onChange={(e) => onChange({ ...challenge, description: e.target.value })}
                    className={textareaClass}
                    placeholder={t('admin.challengeModal.placeholders.description')}
                  />
                </div>
              </div>
            </div>
            <div className="rounded-md border border-neutral-700/70 bg-black/10 p-3">
              <h3 className="text-lg font-mono text-neutral-50 mb-3">{t('admin.challengeModal.sections.ownership')}</h3>
              <div className="space-y-3">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3 lg:gap-4">
                  <div>
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.contest')}
                    </label>
                    {contestLocked ? (
                      <div className={`${selectClass} flex items-center`}>{contestLabel}</div>
                    ) : (
                      <select
                        value={String(challenge.contest_id || 0)}
                        onChange={(e) => changeContest(e.target.value)}
                        className={selectClass}
                        disabled={isEditMode}
                      >
                        <option value="0">{t('admin.challengeModal.contest.unassigned')}</option>
                        {contests.map((contest) => (
                          <option key={contest.id} value={contest.id}>
                            {contest.name}
                          </option>
                        ))}
                      </select>
                    )}
                  </div>
                  <div>
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.attempts')}
                    </label>
                    <input
                      type="number"
                      min="0"
                      value={challenge.attempt || 0}
                      onChange={(e) => onChange({ ...challenge, attempt: Number(e.target.value) || 0 })}
                      className={inputBaseClass}
                    />
                  </div>
                </div>
                <Checkbox
                  checked={!!challenge.hidden}
                  onChange={(e) => onChange({ ...challenge, hidden: e.target.checked })}
                  label={t('admin.challengeModal.labels.hidden')}
                />
                {contestLocked ? (
                  <p className="text-xs font-mono text-neutral-500">{t('admin.challengeModal.contest.fixedHint')}</p>
                ) : (
                  isEditMode && (
                    <p className="text-xs font-mono text-neutral-500">{t('admin.challengeModal.contest.lockedHint')}</p>
                  )
                )}
              </div>
            </div>
            {['tags', 'hints'].map((field) => {
              const singular = field === 'tags' ? 'Tag' : 'Hint';
              return (
                <div key={field} className="rounded-md border border-neutral-700/70 bg-black/10 p-3">
                  <div className="flex flex-wrap justify-between gap-2 mb-3">
                    <h3 className="text-lg font-mono text-neutral-50">{t(`admin.challengeModal.sections.${field}`)}</h3>
                    <Button
                      size="sm"
                      variant="primary"
                      icon={<IconPlus size={16} />}
                      onClick={() => onChange({ ...challenge, [field]: [...(challenge[field] || []), ''] })}
                    >
                      {t(`admin.challengeModal.actions.add${singular}`)}
                    </Button>
                  </div>
                  <div className="space-y-2">
                    {(challenge[field] || []).map((value, index) => (
                      <div key={index} className="flex items-center gap-2">
                        <input
                          type="text"
                          aria-label={t(`admin.challengeModal.placeholders.${singular.toLowerCase()}`, {
                            index: index + 1,
                          })}
                          value={value}
                          placeholder={t(`admin.challengeModal.placeholders.${singular.toLowerCase()}`, {
                            index: index + 1,
                          })}
                          onChange={(e) =>
                            onChange({
                              ...challenge,
                              [field]: (challenge[field] || []).map((item, i) => (i === index ? e.target.value : item)),
                            })
                          }
                          className={inputBaseClass}
                        />
                        <Button
                          size="icon"
                          variant="ghost"
                          onClick={() =>
                            onChange({
                              ...challenge,
                              [field]: (challenge[field] || []).filter((_, i) => i !== index),
                            })
                          }
                        >
                          <IconTrash size={18} />
                        </Button>
                      </div>
                    ))}
                  </div>
                </div>
              );
            })}
            {challenge.type === 'dynamic' && (
              <div className="rounded-md border border-neutral-700/70 bg-black/10 p-3">
                <h3 className="text-lg font-mono text-neutral-50 mb-3">
                  {t('admin.challengeModal.sections.generator')}
                </h3>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.challengeModal.labels.generatorImage')}
                </label>
                <input
                  type="text"
                  value={challenge.generator_image}
                  onChange={(e) => onChange({ ...challenge, generator_image: e.target.value })}
                  className={inputBaseClass}
                  placeholder={t('admin.challengeModal.placeholders.generatorImage')}
                  required
                />
              </div>
            )}
            {challenge.type !== 'pods' && (
              <div className="rounded-md border border-neutral-700/70 bg-black/10 p-3">
                <div className="flex justify-between items-center mb-3">
                  <h3 className="text-lg font-mono text-neutral-50">{t('admin.challengeModal.sections.flags')}</h3>
                  <Button
                    variant="primary"
                    size="sm"
                    align="icon-left"
                    icon={<IconPlus size={16} />}
                    onClick={onAddFlag}
                  >
                    {t('admin.challengeModal.actions.addFlag')}
                  </Button>
                </div>
                <div className="space-y-3">
                  <label className="block text-sm font-mono text-neutral-400 mb-1">
                    {challenge.type === 'static' ? 'static{}' : 'leet{} / uuid{}'}
                  </label>
                  {(challenge.flags || []).map((flag, index) => (
                    <div key={flag.id || `new-${index}`} className="flex gap-2 items-center">
                      <input
                        type="text"
                        value={
                          typeof flag === 'string'
                            ? flag
                            : flag.value || (challenge.type === 'static' ? 'static{}' : 'uuid{}')
                        }
                        onChange={(e) => onFlagChange(index, e.target.value)}
                        className={inputBaseClass}
                        placeholder={t('admin.challengeModal.placeholders.flag', { index: index + 1 })}
                      />
                      {(challenge.flags || []).length > 1 && (
                        <Button
                          variant="ghost"
                          size="icon"
                          className="!bg-transparent !text-red-400 hover:!text-red-300"
                          onClick={() => onRemoveFlag(index)}
                        >
                          <IconTrash size={18} />
                        </Button>
                      )}
                    </div>
                  ))}
                  <p className="text-xs font-mono text-neutral-500">
                    {t('admin.challengeModal.labels.flagChangeNote')}
                  </p>
                </div>
              </div>
            )}
            {challenge.type === 'pods' && (
              <>
                <div className="rounded-md border border-neutral-700/70 bg-black/10 p-3">
                  <h3 className="text-lg font-mono text-neutral-50 mb-3">
                    {t('admin.challengeModal.sections.podsNotice')}
                  </h3>
                  <div className="p-3 bg-geek-400/10 border border-geek-400/20 rounded-md">
                    <p className="text-sm text-geek-400 font-mono">
                      {podNoticeLines.map((line, index) =>
                        line ? (
                          <span key={index}>
                            {line}
                            <br />
                          </span>
                        ) : (
                          <br key={index} />
                        )
                      )}
                    </p>
                  </div>
                </div>
                <ComposeEditor challenge={challenge} onChange={onChange} compose={compose} />
              </>
            )}
          </div>
          {isEditMode && (
            <div className="space-y-4">
              <ChallengeFlagPanel
                challengeId={challenge.id}
                contestId={challenge.contest_id || 0}
                version={flagVersion}
                onChanged={onFlagsChanged}
              />
            </div>
          )}
        </div>
      )}
    </Modal>
  );
}
