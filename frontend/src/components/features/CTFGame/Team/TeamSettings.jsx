/**
 * 队伍设置页面组件
 * @param {Object} props
 * @param {Object} props.team - 队伍信息
 * @param {string} props.team.name - 队伍名称
 * @param {string} props.team.picture - 队伍头像URL
 * @param {string} props.team.inviteCode - 邀请码
 * @param {string} props.team.description - 队伍描述
 * @param {Object} props.team.leader - 队长信息
 * @param {string} props.team.leader.name - 队长用户名
 * @param {string} props.team.leader.picture - 队长头像URL
 * @param {string} props.team.leader.email - 队长邮箱
 * @param {Array<Object>} props.team.members - 队员列表
 * @param {string} props.team.members[].name - 队员用户名
 * @param {string} props.team.members[].picture - 队员头像URL
 * @param {string} props.team.members[].email - 队员邮箱
 * @param {boolean} props.isLeader - 当前用户是否为队长
 * @param {Function} props.onCopyCode - 复制邀请码处理函数
 * @param {Function} props.onRefreshCode - 刷新邀请码处理函数
 * @param {Function} props.onEditTeam - 编辑队伍信息处理函数, 参数为更新后的队伍数据
 * @param {Function} props.onKickMember - 踢出队员处理函数, 参数为队员用户名
 * @param {Function} props.onDisbandTeam - 解散队伍处理函数
 * @param {Function} props.onPictureUpload - 头像上传处理函数, 参数为File对象
 * @example
 * const team = {
 *   name: "Binary Bandits",
 *   picture: "https://example.com/team-picture.jpg",
 *   inviteCode: "BAND-1234-5678",
 *   description: "We hack for fun!",
 *   leader: {
 *     name: "hacker1",
 *     picture: "https://avatars.githubusercontent.com/u/1",
 *     email: "hacker1@example.com"
 *   },
 *   members: [
 *     { name: "hacker2", picture: "https://...", email: "hacker2@example.com" },
 *     { name: "hacker3", picture: "https://...", email: "hacker3@example.com" }
 *   ]
 * }
 *
 *
 *   onCopyCode: () => { ... },
 *   onRefreshCode: () => { ... },
 *   onEditTeam: (updatedTeam) => { ... },
 *   onKickMember: (memberName) => { ... },
 *   onDisbandTeam: () => { ... },
 *   onPictureUpload: (file) => { ... }
 *
 */

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import ConfirmModal from '../../../common/ConfirmModal';
import Avatar from '../../../common/Avatar';
import EditTeamModal from '../EditTeamModal';
import { Button, Card } from '../../../../components/common';
import remarkGfm from 'remark-gfm';
import ReactMarkdown from 'react-markdown';

function TeamSettings({
  team,
  isLeader,
  onCopyCode = () => {},
  onRefreshCode = () => {},
  onEditTeam = () => {},
  onKickMember = () => {},
  onDisbandTeam = () => {},
  onPictureUpload = () => {},
}) {
  const [isCodeCopied, setIsCodeCopied] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showKickModal, setShowKickModal] = useState(null);
  const [showDisbandModal, setShowDisbandModal] = useState(false);
  const { t } = useTranslation();

  const handleCopyCode = () => {
    onCopyCode();
    setIsCodeCopied(true);
    setTimeout(() => setIsCodeCopied(false), 2000);
  };

  const handlePictureUpload = (event) => {
    const file = event.target.files[0];
    event.target.value = '';
    if (file) {
      onPictureUpload(file);
    }
  };

  return (
    <div className="contest-container mx-auto space-y-6">
      {/* 队伍信息卡片 */}
      <Card variant="default" padding="md" animate>
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-4">
            <div className="relative group">
              <Avatar src={team.picture} name={team.name} size="lg" className="border-2 border-neutral-300" />
              {isLeader && (
                <label
                  className="absolute inset-0 flex items-center justify-center 
                                    bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity 
                                    cursor-pointer rounded-lg"
                >
                  <input type="file" className="hidden" accept="image/*" onChange={handlePictureUpload} />
                  <span className="text-neutral-50 text-sm font-mono">{t('game.team.settings.changeAvatar')}</span>
                </label>
              )}
            </div>
            <div>
              <div className="text-neutral-50 font-mono text-lg">{team.name}</div>
              <div className="text-neutral-400 text-sm">
                {t('game.team.settings.members', { count: team.members.length + 1 })}
              </div>
            </div>
          </div>
          {/* 队伍信息卡片底部添加解散按钮 */}
          <div className="flex items-center gap-2">
            {isLeader && (
              <Button variant="danger" size="sm" onClick={() => setShowDisbandModal(true)}>
                {t('game.team.settings.disbandTeam')}
              </Button>
            )}
            {isLeader && (
              <Button variant="primary" size="sm" onClick={() => setShowEditModal(true)}>
                {t('game.team.settings.editTeam')}
              </Button>
            )}
          </div>
        </div>

        {/* 队伍描述 */}
        <div className="mb-6">
          <div className="text-neutral-400 text-sm mb-2">{t('game.team.settings.description')}</div>
          <div className="p-3 bg-neutral-900 rounded-md">
            <div className="text-neutral-300 text-sm prose prose-invert prose-sm line-clamp-3">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {team.description || t('game.team.settings.noDescription')}
              </ReactMarkdown>
            </div>
          </div>
        </div>

        {/* 邀请码部分 */}
        <div className="space-y-2">
          <div className="text-neutral-400 text-sm">{t('game.team.settings.invitationCode')}</div>
          <div className="flex items-center gap-2">
            <div className="flex-1 flex items-center justify-between p-3 bg-neutral-900 rounded-md">
              <span className="font-mono text-neutral-50">{team.inviteCode}</span>
              <Button variant="ghost" className="p-0 min-w-0 h-auto" onClick={handleCopyCode}>
                {isCodeCopied ? '✓' : '📋'}
              </Button>
            </div>
            {isLeader && (
              <Button
                variant="primary"
                size="icon"
                className="border-yellow-400 text-yellow-400 hover:bg-yellow-400/10"
                onClick={onRefreshCode}
              >
                ↻
              </Button>
            )}
          </div>
        </div>
      </Card>

      {/* 成员列表 */}
      <Card variant="default" padding="none" animate className="overflow-hidden">
        <div className="p-4 border-b border-neutral-300/30">
          <h2 className="text-neutral-50 font-mono">{t('game.team.settings.teamMembers')}</h2>
        </div>

        <div className="divide-y divide-neutral-300/10">
          {/* 队长 */}
          <div className="p-4 flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Avatar
                src={team.leader.picture}
                name={team.leader.name}
                size="sm"
                className="border-2 border-yellow-400"
              />
              <div>
                <div className="text-neutral-50 font-mono">{team.leader.name}</div>
                <div className="text-yellow-400 text-sm font-mono">{t('game.team.settings.leader')}</div>
              </div>
            </div>
            <div className="text-neutral-400 text-sm">{team.leader.email}</div>
          </div>

          {/* 队员 */}
          {team.members.map((member, index) => (
            <div key={index} className="p-4 flex items-center justify-between">
              <div className="flex items-center gap-4">
                <Avatar src={member.picture} name={member.name} size="sm" className="border-2 border-neutral-300" />
                <div>
                  <div className="text-neutral-50 font-mono">{member.name}</div>
                  <div className="text-neutral-400 text-sm font-mono">{t('game.team.settings.member')}</div>
                </div>
              </div>
              <div className="flex items-center gap-4">
                <div className="text-neutral-400 text-sm">{member.email}</div>
                {isLeader && (
                  <Button
                    variant="danger"
                    size="icon"
                    className="w-6 h-6"
                    onClick={() => setShowKickModal(member.name)}
                  >
                    ✕
                  </Button>
                )}
              </div>
            </div>
          ))}
        </div>
      </Card>

      {/* 编辑模态框 */}
      <EditTeamModal
        isOpen={showEditModal}
        onClose={() => setShowEditModal(false)}
        team={team}
        onSave={(data) => {
          onEditTeam(data);
          setShowEditModal(false);
        }}
      />

      {/* 踢出确认模态框 */}
      <ConfirmModal
        isOpen={!!showKickModal}
        onClose={() => setShowKickModal(null)}
        onConfirm={() => {
          onKickMember(showKickModal);
          setShowKickModal(null);
        }}
        title={t('game.team.settings.kickModal.title')}
        message={t('game.team.settings.kickModal.message', { name: showKickModal })}
        confirmText={t('game.team.settings.kickModal.confirm')}
        type="danger"
      />

      {/* 解散确认模态框 */}
      <ConfirmModal
        isOpen={showDisbandModal}
        onClose={() => setShowDisbandModal(false)}
        onConfirm={() => {
          onDisbandTeam();
          setShowDisbandModal(false);
        }}
        title={t('game.team.settings.disbandModal.title')}
        message={t('game.team.settings.disbandModal.message')}
        confirmText={t('game.team.settings.disbandModal.confirm')}
        type="danger"
      />
    </div>
  );
}

export default TeamSettings;
