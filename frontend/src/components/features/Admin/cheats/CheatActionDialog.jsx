import { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { deleteContestCheat, deleteAllContestCheats, checkContestCheats } from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';
import { Button, Modal } from '../../../common';

export default function CheatActionDialog({ contestId, action, cheat, onClose, onCompleted }) {
  const { t } = useTranslation();
  const [pending, setPending] = useState(false);
  const busy = useRef(false);
  const close = () => {
    if (!busy.current) onClose();
  };
  const confirmKey = { delete: 'confirmDelete', deleteAll: 'confirmDeleteAll', check: 'confirmCheck' }[action];

  async function submit() {
    if (busy.current || (action === 'delete' && !cheat)) return;
    busy.current = true;
    setPending(true);
    try {
      const response =
        action === 'delete'
          ? await deleteContestCheat(contestId, cheat.id)
          : action === 'deleteAll'
            ? await deleteAllContestCheats(contestId)
            : await checkContestCheats(contestId);
      if (response.code !== 200) throw new Error(t(`admin.contests.cheats.toast.${action}Failed`));
      toast.success({ description: t(`admin.contests.cheats.toast.${action}Success`) });
      onCompleted();
      onClose();
    } catch (error) {
      toast.danger({ description: error.message || t(`admin.contests.cheats.toast.${action}Failed`) });
    } finally {
      busy.current = false;
      setPending(false);
    }
  }

  return (
    <Modal
      isOpen
      onClose={close}
      title={t(`admin.contests.cheats.actions.${confirmKey}Title`)}
      size="sm"
      footer={
        <>
          <Button size="sm" variant="ghost" onClick={close} disabled={pending}>
            {t('common.cancel')}
          </Button>
          <Button size="sm" variant={action === 'check' ? 'primary' : 'danger'} onClick={submit} disabled={pending}>
            {t(`admin.contests.cheats.actions.${action}`)}
          </Button>
        </>
      }
    >
      <p className="text-neutral-300">{t(`admin.contests.cheats.actions.${confirmKey}Prompt`)}</p>
    </Modal>
  );
}
