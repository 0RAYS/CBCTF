import { useState, useEffect } from 'react';
import { toast } from '../../utils/toast';
import { downloadBlobResponse } from '../../utils/fileDownload';
import AdminFiles from '../../components/features/Admin/AdminFiles';
import { Modal } from '../../components/common';
import ModalButton from '../../components/common/ModalButton';
import DeleteConfirmation from '../../components/common/DeleteConfirmation';
import { getFileList, batchDeleteFiles, getFileUrl, downloadFile } from '../../api/admin/file.js';
import { useTranslation } from 'react-i18next';

function FilesManagement() {
  const [files, setFiles] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedFiles, setSelectedFiles] = useState([]);
  const [fileType, setFileType] = useState('file');
  const pageSize = 20;
  const { t } = useTranslation();

  const fetchFiles = async () => {
    try {
      const response = await getFileList({
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
        type: fileType,
      });
      if (response.code === 200) {
        // 处理文件数据
        const processedFiles = response.data.files.map((file) => ({
          id: file.id,
          url: getFileUrl(file.id, file.type),
          filename: file.filename,
          size: file.size,
          type: file.type,
          suffix: file.suffix,
          hash: file.hash,
          uploadTime: file.date,
          model: file.model,
          modelId: file.model_id,
        }));
        setFiles(processedFiles);
        setTotalCount(response.data.count);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.files.toast.fetchFailed') });
    }
  };

  useEffect(() => {
    fetchFiles();
  }, [currentPage, fileType]);

  // 处理删除单个文件
  const handleDelete = (file) => {
    setSelectedFiles([file.id]);
    setIsModalOpen(true);
  };

  // 处理批量删除
  const handleBatchDelete = (selectedIds) => {
    setSelectedFiles(selectedIds);
    setIsModalOpen(true);
  };

  // 处理下载文件
  const handleDownload = async (file) => {
    try {
      const response = await downloadFile(file.id);
      if (response.headers?.['file'] === 'true') {
        downloadBlobResponse(response, file.filename, 'application/octet-stream');
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.files.toast.downloadFailed') });
    }
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    try {
      const response = await batchDeleteFiles({
        file_ids: selectedFiles,
      });
      if (response.code === 200) {
        toast.success({ description: t('admin.files.toast.deleteSuccess') });
        setIsModalOpen(false);
        fetchFiles();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.files.toast.deleteFailed') });
    }
  };

  return (
    <>
      <AdminFiles
        files={files}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        fileType={fileType}
        onPageChange={setCurrentPage}
        onTypeChange={setFileType}
        onDelete={handleDelete}
        onBatchDelete={handleBatchDelete}
        onDownload={handleDownload}
      />

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={t('admin.files.modal.confirmTitle')}
        size="sm"
        footer={
          <>
            <ModalButton variant="default" onClick={() => setIsModalOpen(false)}>
              {t('common.cancel')}
            </ModalButton>
            <ModalButton variant="danger" onClick={handleConfirmDelete}>
              {t('admin.files.actions.confirm')}
            </ModalButton>
          </>
        }
      >
        <DeleteConfirmation
          message={t('admin.files.modal.confirmPrompt', { count: selectedFiles.length })}
          warning={t('admin.files.modal.confirmWarning')}
        />
      </Modal>
    </>
  );
}

export default FilesManagement;
