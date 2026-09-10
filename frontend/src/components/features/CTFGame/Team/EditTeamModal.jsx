/**
 * 编辑队伍模态框
 * @param {Object} props
 * @param {boolean} props.isOpen - 是否显示
 * @param {Function} props.onClose - 关闭回调
 * @param {Object} props.team - 队伍信息
 * @param {Function} props.onSave - 保存回调
 */

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../common';
import Modal from '../../../common/Modal';

function EditTeamModal({ isOpen, onClose, team, onSave }) {
  const { t } = useTranslation();
  const [formData, setFormData] = useState({
    name: team.name,
    description: team.description,
    newLeader: '',
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave(formData);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={t('game.team.editModal.title')} size="md" showCloseButton={false}>
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* 队伍名称 */}
        <div>
          <label className="block text-neutral-400 text-sm mb-2">{t('game.team.editModal.teamName')}</label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) => setFormData((prev) => ({ ...prev, name: e.target.value }))}
            className="w-full p-3 bg-neutral-800/60 border border-neutral-600/60 rounded-md
                                        text-neutral-50 font-mono
                                        focus:outline-none focus:border-geek-400"
          />
        </div>

        {/* 队伍描述 */}
        <div>
          <label className="block text-neutral-400 text-sm mb-2">{t('game.team.editModal.description')}</label>
          <textarea
            value={formData.description.slice(0, 100)}
            onChange={(e) => setFormData((prev) => ({ ...prev, description: e.target.value.slice(0, 100) }))}
            rows={4}
            className="w-full p-3 bg-neutral-800/60 border border-neutral-600/60 rounded-md
                                        text-neutral-50
                                        focus:outline-none focus:border-geek-400"
          />
        </div>

        {/* 转让队长 */}
        <div>
          <label className="block text-neutral-400 text-sm mb-2">{t('game.team.editModal.transferLeadership')}</label>
          <select
            value={formData.newLeader}
            onChange={(e) => setFormData((prev) => ({ ...prev, newLeader: e.target.value }))}
            className="select-custom select-custom-md"
          >
            <option value="">{t('game.team.editModal.selectNewLeader')}</option>
            {team.members.map((member) => (
              <option key={member.name} value={member.name}>
                {member.name}
              </option>
            ))}
          </select>
        </div>

        <div className="flex justify-end gap-4">
          <Button variant="ghost" size="sm" onClick={onClose}>
            {t('game.team.editModal.cancel')}
          </Button>
          <Button variant="primary" size="sm" type="submit">
            {t('game.team.editModal.saveChanges')}
          </Button>
        </div>
      </form>
    </Modal>
  );
}

export default EditTeamModal;
