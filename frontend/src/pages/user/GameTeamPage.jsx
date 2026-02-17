import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { toast } from '../../utils/toast';
import TeamSettings from '../../components/features/CTFGame/Team/TeamSettings';
import {
  getTeamInfo,
  getTeamMembers,
  updateTeamInfo,
  uploadTeamPicture,
  deleteTeam,
  kickTeamMember,
  getTeamCaptcha,
  updateTeamCaptcha,
} from '../../api/game/team';
import { useSelector } from 'react-redux';
import Loading from '../../components/common/Loading';
import { useTranslation } from 'react-i18next';

function GameTeamPage() {
  const { contestId } = useParams();
  const navigate = useNavigate();
  const [team, setTeam] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isLeader, setIsLeader] = useState(false);
  const currentUser = useSelector((state) => state.user.user);
  const { t } = useTranslation();

  // 转换API数据为组件所需格式
  const transformTeamData = async (teamInfo, members) => {
    const leader = members.find((member) => member.id === teamInfo.captain_id);
    const otherMembers = members.filter((member) => member.id !== teamInfo.captain_id);

    return {
      name: teamInfo.name,
      picture: teamInfo.picture,
      description: teamInfo.description,
      inviteCode: await getTeamCaptcha(contestId).then((res) => res.data),
      leader: {
        name: leader.name,
        picture: leader.picture,
        email: leader.email,
      },
      members: otherMembers.map((member) => ({
        name: member.name,
        picture: member.picture,
        email: member.email,
      })),
    };
  };

  // 获取队伍信息
  useEffect(() => {
    const fetchTeamInfo = async () => {
      try {
        const [teamInfoRes, membersRes] = await Promise.all([getTeamInfo(contestId), getTeamMembers(contestId)]);
        if (teamInfoRes.code === 200 && membersRes.code === 200) {
          const transformedData = await transformTeamData(teamInfoRes.data, membersRes.data);
          setTeam(transformedData);
          setIsLeader(currentUser?.id === teamInfoRes.data.captain_id);
        }
      } catch (error) {
        toast.danger({ title: t('game.team.toast.fetchFailed'), description: error.message });
      } finally {
        setLoading(false);
      }
    };

    fetchTeamInfo();
  }, [contestId, currentUser]);

  // 复制邀请码
  const handleCopyCode = () => {
    navigator.clipboard.writeText(team.inviteCode);
    toast.success({ title: t('game.team.toast.inviteCopied') });
  };

  // 刷新邀请码
  const handleRefreshCode = async () => {
    try {
      const res = await updateTeamCaptcha(contestId);
      if (res.code === 200) {
        setTeam((prev) => ({
          ...prev,
          inviteCode: res.data,
        }));
        toast.success({ title: t('game.team.toast.inviteRefreshed') });
      }
    } catch (error) {
      toast.danger({ title: t('game.team.toast.inviteRefreshFailed'), description: error.message });
    }
  };

  // 编辑队伍信息
  const handleEditTeam = async (updatedTeam) => {
    try {
      const res = await updateTeamInfo(contestId, {
        name: updatedTeam.name,
        description: updatedTeam.description,
      });

      if (res.code === 200) {
        setTeam((prev) => ({
          ...prev,
          name: updatedTeam.name,
          description: updatedTeam.description,
        }));
        toast.success({ title: t('game.team.toast.updateSuccess') });
      }
    } catch (error) {
      toast.danger({ title: t('game.team.toast.updateFailed'), description: error.message });
    }
  };

  // 踢出队员
  const handleKickMember = async (memberName) => {
    try {
      const member = team.members.find((m) => m.name === memberName);
      if (!member) return;

      const res = await kickTeamMember(contestId, member.id);
      if (res.code === 200) {
        setTeam((prev) => ({
          ...prev,
          members: prev.members.filter((m) => m.name !== memberName),
        }));
        toast.success({
          title: t('game.team.toast.memberRemoved'),
          description: t('game.team.toast.memberRemovedDescription', { name: memberName }),
        });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.team.toast.removeFailed') });
    }
  };

  // 解散队伍
  const handleDisbandTeam = async () => {
    try {
      const res = await deleteTeam(contestId);
      if (res.code === 200) {
        toast.success({ title: t('game.team.toast.disbanded') });
        navigate(`/contests/${contestId}`);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.team.toast.disbandFailed') });
    }
  };

  // 上传头像
  const handlePictureUpload = async (file) => {
    try {
      const res = await uploadTeamPicture(contestId, file);
      if (res.code === 200) {
        // 重新获取队伍信息以获取新的头像URL
        const teamInfoRes = await getTeamInfo(contestId);
        if (teamInfoRes.code === 200) {
          setTeam((prev) => ({
            ...prev,
            picture: teamInfoRes.data.picture,
          }));
          toast.success({ title: t('game.team.toast.avatarUpdated') });
        }
      }
    } catch (error) {
      toast.danger({ title: t('game.team.toast.avatarUploadFailed'), description: error.message });
    }
  };

  if (loading) {
    return <Loading />;
  }

  if (!team) {
    return (
      <div className="flex items-center justify-center min-h-[500px]">
        <span className="text-neutral-300">{t('game.team.empty')}</span>
      </div>
    );
  }

  return (
    <TeamSettings
      team={team}
      isLeader={isLeader}
      onCopyCode={handleCopyCode}
      onRefreshCode={handleRefreshCode}
      onEditTeam={handleEditTeam}
      onKickMember={handleKickMember}
      onDisbandTeam={handleDisbandTeam}
      onPictureUpload={handlePictureUpload}
    />
  );
}

export default GameTeamPage;
