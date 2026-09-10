import { motion } from 'motion/react';
import { useState } from 'react';
import { Card } from '../../../common';
import { useTranslation } from 'react-i18next';

function WriteupUpload({ onUploadWriteup, writeups = [] }) {
  const { t, i18n } = useTranslation();
  const [isDragging, setIsDragging] = useState(false);

  const formatFileSize = (sizeInBytes) => {
    if (sizeInBytes < 1024) return `${sizeInBytes} B`;
    if (sizeInBytes < 1024 * 1024) return `${(sizeInBytes / 1024).toFixed(1)} KB`;
    return `${(sizeInBytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString(i18n.language || 'en-US', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragging(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      onUploadWriteup(e.dataTransfer.files[0]);
    }
  };

  return (
    <Card variant="default" padding="lg" animate>
      <div className="space-y-4">
        <div className="text-neutral-50 font-mono text-lg">{t('game.contestEnded.writeups.title')}</div>

        <motion.div
          className={`border-2 border-dashed rounded-md p-6 text-center cursor-pointer transition-all duration-300
            ${isDragging ? 'border-yellow-400 bg-yellow-400/10' : 'border-neutral-400 hover:border-yellow-400 hover:bg-yellow-400/5'}`}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={() => document.getElementById('writeup-upload').click()}
        >
          <input
            id="writeup-upload"
            type="file"
            className="hidden"
            accept=".pdf,.doc,.docx"
            onChange={(e) => {
              if (e.target.files && e.target.files[0]) {
                onUploadWriteup(e.target.files[0]);
                e.target.value = '';
              }
            }}
          />
          <div className="flex flex-col items-center justify-center space-y-2">
            <span className="text-2xl">📄</span>
            <span className="text-neutral-300 font-mono">{t('game.contestEnded.writeups.dropzone')}</span>
            <span className="text-neutral-400 text-sm">{t('game.contestEnded.writeups.supports')}</span>
          </div>
        </motion.div>

        {writeups.length > 0 && (
          <motion.div
            className="mt-4 border border-neutral-300/30 rounded-md overflow-hidden"
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
          >
            <div className="p-3 bg-neutral-800/50 border-b border-neutral-300/30">
              <div className="text-neutral-300 font-mono text-sm">{t('game.contestEnded.writeups.uploaded')}</div>
            </div>
            <div className="divide-y divide-neutral-300/20">
              {writeups.map((file) => (
                <div key={file.id} className="p-3 flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <div className="text-xl">
                      {file.suffix === '.pdf' ? '📕' : file.suffix === '.docx' || file.suffix === '.doc' ? '📘' : '📄'}
                    </div>
                    <div>
                      <div className="text-neutral-50 font-mono text-sm">{file.filename}</div>
                      <div className="text-neutral-400 text-xs flex space-x-3">
                        <span>{formatFileSize(file.size)}</span>
                        <span>•</span>
                        <span>{formatDate(file.date)}</span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </motion.div>
        )}
      </div>
    </Card>
  );
}

export default WriteupUpload;
