import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { toast } from '../../../utils/toast';
import { getContestInfo, updateContestInfo, updateContestPicture } from '../../../api/admin/contest';
import ContestEditor from '../../../components/features/Admin/Contests/editor/ContestEditor';
import {
  createContestDraft,
  contestUpdatePayload,
} from '../../../components/features/Admin/Contests/editor/contestForm.js';

export default function AdminContestDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [loaded, setLoaded] = useState(null);
  useEffect(() => {
    let cancelled = false;
    getContestInfo(Number(id))
      .then((response) => {
        if (!cancelled && response.code === 200) setLoaded({ id, draft: createContestDraft(response.data) });
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.contests.editor.toast.fetchFailed') });
      });
    return () => {
      cancelled = true;
    };
  }, [id, t]);
  const upload = async (file) => {
    try {
      const response = await updateContestPicture(id, file);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.editor.toast.coverUpdated') });
        return response.data.picture;
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.editor.toast.coverUpdateFailed') });
    }
    return null;
  };
  const save = async (draft) => {
    try {
      const response = await updateContestInfo(Number(id), contestUpdatePayload(draft));
      if (response.code === 200) toast.success({ description: t('admin.contests.editor.toast.updateSuccess') });
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.editor.toast.updateFailed') });
    }
  };
  return loaded?.id === id ? (
    <ContestEditor
      key={id}
      contest={loaded.draft}
      onSave={save}
      onImageUpload={upload}
      onCancel={() => navigate('/admin/contests')}
    />
  ) : null;
}
