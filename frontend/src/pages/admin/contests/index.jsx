import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import { getContestInfo, updateContestInfo, updateContestPicture } from '../../../api/admin/contest';
import AdminContestEditor from '../../../components/features/Admin/Contests/AdminContestEditor';
import { useTranslation } from 'react-i18next';

function AdminContestDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [contest, setContest] = useState(null);
  const { t } = useTranslation();

  const getContestStatus = (contest) => {
    const now = new Date().getTime() / 1000;
    const startTime = new Date(contest?.start).getTime() / 1000;
    const duration = contest?.duration;
    const endTime = startTime + duration;

    if (now < startTime) {
      return 'upcoming';
    } else if (now > endTime) {
      return 'ended';
    } else {
      return 'active';
    }
  };

  const calculateEndTime = (startTime, duration) => {
    if (!startTime || !duration) return '';
    const start = new Date(startTime);
    const end = new Date(start.getTime() + duration * 1000);
    return end.toISOString();
  };

  const fetchContestInfo = async () => {
    try {
      const response = await getContestInfo(parseInt(id));
      if (response.code === 200) {
        // 转换数据格式以适配编辑器组件
        const contestData = {
          title: response.data.name,
          description: response.data.description,
          image: response.data.picture,
          status: getContestStatus(response.data),
          startTime: response.data.start,
          endTime: calculateEndTime(response.data.start, response.data.duration),
          participants: response.data.users,
          rules: response.data.rules || [],
          prizes:
            response.data.prizes?.map((prize) => ({
              amount: prize.amount,
              description: prize.description,
            })) || [],
          timeline: response.data.timelines
            ? response.data.timelines.map((item) => ({
                date: item.date,
                title: item.title,
                description: item.description,
              }))
            : [],
          // 保存原始数据以便更新时使用
          prefix: response.data.prefix,
          size: response.data.size,
          hidden: response.data.hidden,
          captcha: response.data.captcha,
          blood: response.data.blood,
          victims: response.data.victims,
        };
        setContest(contestData);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.editor.toast.fetchFailed') });
    }
  };

  useEffect(() => {
    fetchContestInfo();
  }, [id]);

  const handleImageUpload = async (file) => {
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

  const handleSave = async (updatedContest) => {
    try {
      // 将编辑器的数据格式转换回API所需格式
      const updateData = {
        name: updatedContest.title,
        description: updatedContest.description,
        // 更新原始数据
        prefix: updatedContest.prefix,
        size: updatedContest.size,
        hidden: updatedContest.hidden,
        captcha: updatedContest.captcha,
        blood: updatedContest.blood,
        // 计算新的开始时间和持续时间
        start: new Date(updatedContest.startTime).toISOString(),
        duration: Math.floor((new Date(updatedContest.endTime) - new Date(updatedContest.startTime)) / 1000),
        victims: updatedContest.victims,
        // 转换timelines格式
        timelines: updatedContest.timeline.map((item) => ({
          date: item.date ? new Date(item.date).toISOString() : '',
          title: item.title,
          description: item.description,
        })),
        rules: updatedContest.rules,
        prizes: updatedContest.prizes.map((prize) => ({
          amount: prize.amount,
          description: prize.description,
        })),
      };

      const response = await updateContestInfo(parseInt(id), updateData);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.editor.toast.updateSuccess') });
        await fetchContestInfo();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.editor.toast.updateFailed') });
    }
  };

  const handleCancel = () => {
    // 返回到比赛列表页面
    navigate('/admin/contests');
  };

  return contest ? (
    <AdminContestEditor
      contest={contest}
      onSave={handleSave}
      onCancel={handleCancel}
      onImageUpload={handleImageUpload}
    />
  ) : null;
}

export default AdminContestDetail;
