import { useState, useRef, useEffect } from 'react';
import { Button, Modal } from '../../../../components/common';
import { useTranslation } from 'react-i18next';

function TeamJoinModal({ isOpen, onClose, onCreateTeam, onJoinTeam }) {
  const [mode, setMode] = useState('select'); // 'select' | 'create' | 'join'
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const { t } = useTranslation();
  const activeRef = useRef(false);
  const submittingRef = useRef(false);
  useEffect(() => {
    activeRef.current = isOpen;
    return () => {
      activeRef.current = false;
    };
  }, [isOpen]);

  // 创建队伍表单
  const [createForm, setCreateForm] = useState({
    teamName: '',
    description: '',
    contestCode: '',
  });

  // 加入队伍表单
  const [joinForm, setJoinForm] = useState({
    teamName: '',
    teamCode: '',
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (submittingRef.current) return;
    submittingRef.current = true;
    setError(null);
    setLoading(true);

    try {
      const success = mode === 'create' ? await onCreateTeam(createForm) : await onJoinTeam(joinForm);
      if (!activeRef.current) return;
      if (success === true) onClose();
      else setError(t(mode === 'create' ? 'toast.team.createFailed' : 'toast.team.joinFailed'));
    } catch (err) {
      if (activeRef.current) setError(err.message || t('errors.requestFailed'));
    } finally {
      submittingRef.current = false;
      if (activeRef.current) setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={() => {
        if (!submittingRef.current) onClose();
      }}
      title={<span className="block break-words text-lg sm:text-xl">{t(`game.team.joinModal.title.${mode}`)}</span>}
      size="md"
      className="max-w-[500px]"
    >
      {mode === 'select' ? (
        <div className="space-y-4">
          <Button
            variant="primary"
            size="lg"
            fullWidth
            onClick={() => {
              setError(null);
              setMode('create');
            }}
          >
            {t('game.team.joinModal.createTeam')}
          </Button>

          <div className="flex items-center gap-4">
            <div className="h-[1px] flex-1 bg-neutral-300/30" />
            <span className="text-neutral-400 text-sm">{t('game.team.joinModal.or')}</span>
            <div className="h-[1px] flex-1 bg-neutral-300/30" />
          </div>

          <Button
            variant="outline"
            size="lg"
            fullWidth
            onClick={() => {
              setError(null);
              setMode('join');
            }}
          >
            {t('game.team.joinModal.joinTeam')}
          </Button>
        </div>
      ) : mode === 'create' ? (
        <form key="create" className="space-y-4" onSubmit={handleSubmit} aria-busy={loading}>
          <div className="space-y-2">
            <label htmlFor="create-team-name" className="text-neutral-400 text-sm">
              {t('game.team.joinModal.form.teamName')}
            </label>
            <input
              id="create-team-name"
              type="text"
              required
              value={createForm.teamName}
              onChange={(e) =>
                setCreateForm((prev) => ({
                  ...prev,
                  teamName: e.target.value,
                }))
              }
              className="w-full h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4
                                                text-neutral-50 placeholder-neutral-400
                                                focus:border-geek-400 focus:shadow-focus
                                                transition-all duration-200"
              placeholder={t('game.team.joinModal.form.teamNamePlaceholder')}
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="create-team-description" className="text-neutral-400 text-sm">
              {t('game.team.joinModal.form.description')}
            </label>
            <textarea
              id="create-team-description"
              value={createForm.description}
              onChange={(e) =>
                setCreateForm((prev) => ({
                  ...prev,
                  description: e.target.value,
                }))
              }
              className="w-full h-[80px] bg-black/20 border border-neutral-300 rounded-md p-4
                                                text-neutral-50 placeholder-neutral-400
                                                focus:border-geek-400 focus:shadow-focus
                                                transition-all duration-200 resize-none"
              placeholder={t('game.team.joinModal.form.descriptionPlaceholder')}
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="create-contest-code" className="text-neutral-400 text-sm">
              {t('game.team.joinModal.form.contestCode')}
            </label>
            <input
              id="create-contest-code"
              type="text"
              value={createForm.contestCode}
              onChange={(e) =>
                setCreateForm((prev) => ({
                  ...prev,
                  contestCode: e.target.value,
                }))
              }
              className="w-full h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4
                                                text-neutral-50 placeholder-neutral-400
                                                focus:border-geek-400 focus:shadow-focus
                                                transition-all duration-200"
              placeholder={t('game.team.joinModal.form.contestCodePlaceholder')}
            />
          </div>

          {error && (
            <div role="alert" className="text-red-400 text-sm break-words">
              {error}
            </div>
          )}

          <div className="flex items-center gap-3 pt-2">
            <Button type="button" variant="outline" fullWidth disabled={loading} onClick={() => setMode('select')}>
              {t('game.team.joinModal.actions.back')}
            </Button>
            <Button
              type="submit"
              variant="primary"
              fullWidth
              disabled={!createForm.teamName.trim() || loading}
              loading={loading}
            >
              {loading ? t('game.team.joinModal.actions.creating') : t('game.team.joinModal.actions.create')}
            </Button>
          </div>
        </form>
      ) : (
        <form key="join" className="space-y-4" onSubmit={handleSubmit} aria-busy={loading}>
          <div className="space-y-2">
            <label htmlFor="join-team-name" className="text-neutral-400 text-sm">
              {t('game.team.joinModal.form.teamName')}
            </label>
            <input
              id="join-team-name"
              type="text"
              required
              value={joinForm.teamName}
              onChange={(e) =>
                setJoinForm((prev) => ({
                  ...prev,
                  teamName: e.target.value,
                }))
              }
              className="w-full h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4
                                                text-neutral-50 placeholder-neutral-400
                                                focus:border-geek-400 focus:shadow-focus
                                                transition-all duration-200"
              placeholder={t('game.team.joinModal.form.teamNamePlaceholder')}
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="join-team-code" className="text-neutral-400 text-sm">
              {t('game.team.joinModal.form.inviteCode')}
            </label>
            <input
              id="join-team-code"
              type="text"
              required
              value={joinForm.teamCode}
              onChange={(e) =>
                setJoinForm((prev) => ({
                  ...prev,
                  teamCode: e.target.value,
                }))
              }
              className="w-full h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4
                                                text-neutral-50 placeholder-neutral-400
                                                focus:border-geek-400 focus:shadow-focus
                                                transition-all duration-200"
              placeholder={t('game.team.joinModal.form.inviteCodePlaceholder')}
            />
          </div>

          {error && (
            <div role="alert" className="text-red-400 text-sm break-words">
              {error}
            </div>
          )}

          <div className="flex items-center gap-3 pt-2">
            <Button type="button" variant="outline" fullWidth disabled={loading} onClick={() => setMode('select')}>
              {t('game.team.joinModal.actions.back')}
            </Button>
            <Button
              type="submit"
              variant="primary"
              fullWidth
              disabled={!joinForm.teamName.trim() || !joinForm.teamCode.trim() || loading}
              loading={loading}
            >
              {loading ? t('game.team.joinModal.actions.joining') : t('game.team.joinModal.actions.join')}
            </Button>
          </div>
        </form>
      )}
    </Modal>
  );
}

export default TeamJoinModal;
