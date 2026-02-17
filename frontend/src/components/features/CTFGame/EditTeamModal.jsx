/**
 * 编辑队伍模态框
 * @param {Object} props
 * @param {boolean} props.isOpen - 是否显示
 * @param {Function} props.onClose - 关闭回调
 * @param {Object} props.team - 队伍信息
 * @param {Function} props.onSave - 保存回调
 */

import { motion, AnimatePresence } from 'motion/react';
import { useState } from 'react';
import { Button } from '../../../components/common';

function EditTeamModal({ isOpen, onClose, team, onSave }) {
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
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center">
          <motion.div
            className="fixed inset-0 bg-black/60 backdrop-blur-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
          />
          <motion.div
            className="relative w-full max-w-[600px] m-4 border border-neutral-300 rounded-md bg-black/80"
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.9, opacity: 0 }}
          >
            <div className="p-6 border-b border-neutral-300/30">
              <h3 className="text-lg font-mono text-neutral-50">Edit Team</h3>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-6">
              {/* 队伍名称 */}
              <div>
                <label className="block text-neutral-400 text-sm mb-2">Team Name</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData((prev) => ({ ...prev, name: e.target.value }))}
                  className="w-full p-3 bg-neutral-900 border border-neutral-300/30 rounded-md
                                        text-neutral-50 font-mono
                                        focus:outline-none focus:border-geek-400"
                />
              </div>

              {/* 队伍描述 */}
              <div>
                <label className="block text-neutral-400 text-sm mb-2">Description</label>
                <textarea
                  value={formData.description.slice(0, 100)}
                  onChange={(e) => setFormData((prev) => ({ ...prev, description: e.target.value.slice(0, 100) }))}
                  rows={4}
                  className="w-full p-3 bg-neutral-900 border border-neutral-300/30 rounded-md
                                        text-neutral-50
                                        focus:outline-none focus:border-geek-400"
                />
              </div>

              {/* 转让队长 */}
              <div>
                <label className="block text-neutral-400 text-sm mb-2">Transfer Leadership</label>
                <select
                  value={formData.newLeader}
                  onChange={(e) => setFormData((prev) => ({ ...prev, newLeader: e.target.value }))}
                  className="select-custom select-custom-md"
                >
                  <option value="">Select a new leader</option>
                  {team.members.map((member) => (
                    <option key={member.name} value={member.name}>
                      {member.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex justify-end gap-4">
                <Button variant="ghost" size="sm" onClick={onClose}>
                  CANCEL
                </Button>
                <Button variant="primary" size="sm" type="submit">
                  SAVE CHANGES
                </Button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
}

export default EditTeamModal;
