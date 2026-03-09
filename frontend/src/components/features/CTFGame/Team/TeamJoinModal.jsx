import { motion, AnimatePresence } from 'motion/react';
import { useState } from 'react';
import { Button } from '../../../../components/common';
import { useTranslation } from 'react-i18next';

function TeamJoinModal({ isOpen, onClose, onCreateTeam, onJoinTeam }) {
  const [mode, setMode] = useState('select'); // 'select' | 'create' | 'join'
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const { t } = useTranslation();

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
    setError(null);
    setLoading(true);

    try {
      if (mode === 'create') {
        await onCreateTeam(createForm);
      } else if (mode === 'join') {
        await onJoinTeam(joinForm);
      }
      onClose();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[900] flex items-center justify-center">
      <motion.div
        className="fixed inset-0 bg-black/60 backdrop-blur-sm"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        onClick={onClose}
      />

      <div className="relative z-10 w-full max-w-[500px] p-4">
        <motion.div
          className="relative w-full bg-black/80 border border-neutral-300 rounded-md overflow-hidden"
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
        >
          {/* 头部 */}
          <div className="p-5 border-b border-neutral-300/30">
            <div className="flex items-center justify-between">
              <h2 className="text-2xl text-neutral-50 font-mono">
                {t(`game.team.joinModal.title.${mode}`)}
              </h2>
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

          {/* 内容区域 */}
          <div className="p-6">
            <AnimatePresence mode="wait">
              {mode === 'select' ? (
                <motion.div
                  key="select"
                  className="space-y-4"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                >
                  <Button
                    variant="primary"
                    size="lg"
                    fullWidth
                    className="shadow-[0_0_20px_rgba(89,126,247,0.4)]"
                    onClick={() => setMode('create')}
                  >
                    {t('game.team.joinModal.createTeam')}
                  </Button>

                  <div className="flex items-center gap-4">
                    <div className="h-[1px] flex-1 bg-neutral-300/30" />
                    <span className="text-neutral-400 text-sm">{t('game.team.joinModal.or')}</span>
                    <div className="h-[1px] flex-1 bg-neutral-300/30" />
                  </div>

                  <Button variant="outline" size="lg" fullWidth onClick={() => setMode('join')}>
                    {t('game.team.joinModal.joinTeam')}
                  </Button>
                </motion.div>
              ) : mode === 'create' ? (
                <motion.form
                  key="create"
                  className="space-y-4"
                  onSubmit={handleSubmit}
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                >
                  <div className="space-y-2">
                    <label className="text-neutral-400 text-sm">
                      {t('game.team.joinModal.form.teamName')}
                    </label>
                    <input
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
                    <label className="text-neutral-400 text-sm">
                      {t('game.team.joinModal.form.description')}
                    </label>
                    <textarea
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
                    <label className="text-neutral-400 text-sm">
                      {t('game.team.joinModal.form.contestCode')}
                    </label>
                    <input
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
                    <motion.div
                      className="text-red-400 text-sm"
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                    >
                      {error}
                    </motion.div>
                  )}

                  <div className="flex items-center gap-3 pt-2">
                    <Button type="button" variant="outline" fullWidth onClick={() => setMode('select')}>
                      {t('game.team.joinModal.actions.back')}
                    </Button>
                    <Button
                      type="submit"
                      variant="primary"
                      fullWidth
                      className="shadow-[0_0_20px_rgba(89,126,247,0.4)]"
                      disabled={!createForm.teamName || loading}
                      loading={loading}
                    >
                      {loading
                        ? t('game.team.joinModal.actions.creating')
                        : t('game.team.joinModal.actions.create')}
                    </Button>
                  </div>
                </motion.form>
              ) : (
                <motion.form
                  key="join"
                  className="space-y-4"
                  onSubmit={handleSubmit}
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                >
                  <div className="space-y-2">
                    <label className="text-neutral-400 text-sm">
                      {t('game.team.joinModal.form.teamName')}
                    </label>
                    <input
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
                    <label className="text-neutral-400 text-sm">
                      {t('game.team.joinModal.form.inviteCode')}
                    </label>
                    <input
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
                    <motion.div
                      className="text-red-400 text-sm"
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                    >
                      {error}
                    </motion.div>
                  )}

                  <div className="flex items-center gap-3 pt-2">
                    <Button type="button" variant="outline" fullWidth onClick={() => setMode('select')}>
                      {t('game.team.joinModal.actions.back')}
                    </Button>
                    <Button
                      type="submit"
                      variant="primary"
                      fullWidth
                      className="shadow-[0_0_20px_rgba(89,126,247,0.4)]"
                      disabled={!joinForm.teamName || !joinForm.teamCode || loading}
                      loading={loading}
                    >
                      {loading
                        ? t('game.team.joinModal.actions.joining')
                        : t('game.team.joinModal.actions.join')}
                    </Button>
                  </div>
                </motion.form>
              )}
            </AnimatePresence>
          </div>
        </motion.div>
      </div>
    </div>
  );
}

export default TeamJoinModal;
