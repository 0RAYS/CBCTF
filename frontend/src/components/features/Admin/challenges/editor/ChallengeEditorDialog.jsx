import { IconPlus, IconTrash } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import Button from '../../../../common/Button';
import Modal from '../../../../common/Modal';
import ComposeEditor from './ComposeEditor.jsx';
import useComposeGuide from './useComposeGuide.js';
import { inputBaseClass, selectClass, textareaClass } from './editorStyles.js';

export default function ChallengeEditorDialog({
  isOpen = false,
  mode = 'add',
  challenge = {
    name: '',
    description: '',
    category: '',
    type: 'static',
    generator_image: '',
    flags: [''],
    docker_compose: '',
    network_policies: [{ from: [{ cidr: '', except: [''] }], to: [{ cidr: '', except: [''] }] }],
  },
  categories = [],
  onClose,
  onSubmit,
  onChange,
  onAddFlag,
  onRemoveFlag,
  onFlagChange,
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
            <h3 className="text-2xl font-mono text-neutral-50 mb-4">{t('admin.challengeModal.delete.title')}</h3>
            <p className="text-neutral-400 font-mono text-lg mb-2">{t('admin.challengeModal.delete.prompt')}</p>
            <p className="text-red-400 font-mono text-sm">{t('admin.challengeModal.delete.warning')}</p>
          </div>
        </div>
      ) : (
        <div className={isEditMode ? 'space-y-3' : 'space-y-4'}>
          <div className={isEditMode ? 'rounded-md border border-neutral-700/70 bg-black/10 p-3' : 'mb-4'}>
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
                    <option value="static">{t('admin.challengeModal.types.static')}</option>
                    <option value="dynamic">{t('admin.challengeModal.types.dynamic')}</option>
                    <option value="pods">{t('admin.challengeModal.types.pods')}</option>
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
          {challenge.type === 'dynamic' && (
            <div className="border-t border-neutral-700 pt-3 lg:pt-4">
              <h3 className="text-lg font-mono text-neutral-50 mb-3">{t('admin.challengeModal.sections.generator')}</h3>
              <div>
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
            </div>
          )}
          {challenge.type !== 'pods' && (
            <div className="border-t border-neutral-700 pt-3 lg:pt-4">
              <div className="flex justify-between items-center mb-3">
                <h3 className="text-lg font-mono text-neutral-50">{t('admin.challengeModal.sections.flags')}</h3>
                <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onAddFlag}>
                  {t('admin.challengeModal.actions.addFlag')}
                </Button>
              </div>
              <div className="space-y-3">
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {challenge.type === 'static' ? 'static{}' : 'leet{} / uuid{}'}
                </label>
                {challenge.flags.map((flag, index) => (
                  <div key={index} className="flex gap-2 items-center">
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
                    {challenge.flags.length > 1 && (
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
              </div>
            </div>
          )}
          {challenge.type === 'pods' && (
            <>
              <div className="border-t border-neutral-700 pt-3 lg:pt-4">
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
      )}
    </Modal>
  );
}
