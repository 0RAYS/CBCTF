import { useState, useEffect, useCallback, useMemo } from 'react';
import { IconPlus, IconTrash } from '@tabler/icons-react';
import { toast } from '../../utils/toast';
import {
  getWebhookList,
  createWebhook,
  updateWebhook,
  deleteWebhook,
  getEvents,
  getAllWebhookHistory,
  getWebhookHistory,
} from '../../api/admin/webhook';
import AdminWebhook from '../../components/features/Admin/AdminWebhook';
import AdminWebhookHistory from '../../components/features/Admin/AdminWebhookHistory';
import { Modal } from '../../components/common';
import ModalButton from '../../components/common/ModalButton';
import { Button } from '../../components/common/index.js';
import Input from '../../components/common/Input';
import { useTranslation } from 'react-i18next';

function WebhookManagement() {
  // 状态管理
  const [webhooks, setWebhooks] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedWebhook, setSelectedWebhook] = useState(null);
  const [mode, setMode] = useState('edit'); // 'edit' | 'create' | 'delete'
  const [editForm, setEditForm] = useState({
    name: '',
    url: '',
    method: 'POST',
    headers: {},
    timeout: 0,
    retry: 0,
    events: [],
    on: false,
  });

  // 防抖状态，用于优化请求头key的更新
  const [headerKeyUpdates, setHeaderKeyUpdates] = useState({});

  // 历史记录状态
  const [webhookHistory, setWebhookHistory] = useState([]);
  const [historyTotalCount, setHistoryTotalCount] = useState(0);
  const [historyCurrentPage, setHistoryCurrentPage] = useState(1);
  const [selectedHistory, setSelectedHistory] = useState(null);
  const [isHistoryDetailModalOpen, setIsHistoryDetailModalOpen] = useState(false);
  const [activeTab, setActiveTab] = useState('webhook'); // 'webhook' | 'history'
  const [currentWebhookId, setCurrentWebhookId] = useState(null);
  const { t } = useTranslation();

  // 事件列表
  const [availableEvents, setAvailableEvents] = useState([]);

  const fetchEvents = async () => {
    try {
      const response = await getEvents();
      if (response.code === 200) {
        setAvailableEvents(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.webhook.toast.fetchEventsFailed') });
    }
  };

  const fetchWebhooks = async () => {
    try {
      const response = await getWebhookList({
        limit: 20,
        offset: (currentPage - 1) * 20,
      });

      if (response.code === 200) {
        setWebhooks(response.data.webhooks || response.data);
        setTotalCount(response.data.count || response.data.length);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.webhook.toast.fetchListFailed') });
    }
  };

  const fetchWebhookHistory = async () => {
    try {
      let response;
      if (currentWebhookId) {
        response = await getWebhookHistory(currentWebhookId, {
          limit: 20,
          offset: (historyCurrentPage - 1) * 20,
        });
      } else {
        response = await getAllWebhookHistory({
          limit: 20,
          offset: (historyCurrentPage - 1) * 20,
        });
      }
      if (response.code === 200) {
        setWebhookHistory(response.data.histories || response.data);
        setHistoryTotalCount(response.data.count || response.data.length);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.webhook.toast.fetchHistoryFailed') });
    }
  };

  // 数据获取
  useEffect(() => {
    if (activeTab === 'webhook') {
      fetchWebhooks();
    } else if (activeTab === 'history') {
      fetchWebhookHistory();
    }
  }, [currentPage, historyCurrentPage, activeTab, currentWebhookId]);

  // 获取可用事件列表
  useEffect(() => {
    fetchEvents();
  }, []);

  // 清理定时器
  useEffect(() => {
    return () => {
      Object.values(headerKeyUpdates).forEach((timeoutId) => {
        if (timeoutId) {
          clearTimeout(timeoutId);
        }
      });
    };
  }, [headerKeyUpdates]);

  const handleCreateClick = () => {
    setMode('create');
    setSelectedWebhook(null);
    setEditForm({
      name: '',
      url: '',
      method: 'POST',
      headers: {},
      timeout: 0,
      retry: 0,
      events: [],
      on: false,
    });
    setIsModalOpen(true);
  };

  const handleEditClick = (webhook) => {
    setMode('edit');
    setSelectedWebhook(webhook);
    setEditForm({
      name: webhook.name,
      url: webhook.url,
      method: webhook.method,
      headers: webhook.headers || {},
      timeout: webhook.timeout,
      retry: webhook.retry,
      events: webhook.events || [],
      on: webhook.on,
    });
    setIsModalOpen(true);
  };

  const handleDeleteClick = (webhook) => {
    setMode('delete');
    setSelectedWebhook(webhook);
    setIsModalOpen(true);
  };

  const handleViewHistory = (webhook) => {
    setCurrentWebhookId(webhook.id);
    setActiveTab('history');
    setHistoryCurrentPage(1);
  };

  const handleViewAllHistory = () => {
    setCurrentWebhookId(null);
    setActiveTab('history');
    setHistoryCurrentPage(1);
  };

  const handleWebhookClick = (webhook) => {
    handleEditClick(webhook);
  };

  const handleHistoryClick = (history) => {
    setSelectedHistory(history);
    setIsHistoryDetailModalOpen(true);
  };

  const handleCreateWebhook = async () => {
    try {
      // 清理空的请求头
      const cleanedHeaders = {};
      Object.entries(editForm.headers).forEach(([key, value]) => {
        if (key && key.trim() !== '' && value && value.trim() !== '') {
          cleanedHeaders[key.trim()] = value.trim();
        }
      });

      const formData = { ...editForm, headers: cleanedHeaders };
      const response = await createWebhook(formData);
      if (response.code === 200) {
        toast.success({ description: t('admin.webhook.toast.createSuccess') });
        setIsModalOpen(false);
        fetchWebhooks();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.webhook.toast.createFailed') });
    }
  };

  const handleUpdateWebhook = async () => {
    try {
      // 清理空的请求头
      const cleanedHeaders = {};
      Object.entries(editForm.headers).forEach(([key, value]) => {
        if (key && key.trim() !== '' && value && value.trim() !== '') {
          cleanedHeaders[key.trim()] = value.trim();
        }
      });

      // 只发送有值的字段
      const updateData = {};
      if (editForm.name !== selectedWebhook.name) updateData.name = editForm.name;
      if (editForm.url !== selectedWebhook.url) updateData.url = editForm.url;
      if (editForm.method !== selectedWebhook.method) updateData.method = editForm.method;
      if (JSON.stringify(cleanedHeaders) !== JSON.stringify(selectedWebhook.headers))
        updateData.headers = cleanedHeaders;
      if (editForm.timeout !== selectedWebhook.timeout) updateData.timeout = editForm.timeout;
      if (editForm.retry !== selectedWebhook.retry) updateData.retry = editForm.retry;
      if (JSON.stringify(editForm.events) !== JSON.stringify(selectedWebhook.events))
        updateData.events = editForm.events;
      if (editForm.on !== selectedWebhook.on) updateData.on = editForm.on;

      const response = await updateWebhook(selectedWebhook.id, updateData);
      if (response.code === 200) {
        toast.success({ description: t('admin.webhook.toast.updateSuccess') });
        setIsModalOpen(false);
        fetchWebhooks();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.webhook.toast.updateFailed') });
    }
  };

  const handleDeleteWebhook = async () => {
    try {
      const response = await deleteWebhook(selectedWebhook.id);
      if (response.code === 200) {
        toast.success({ description: t('admin.webhook.toast.deleteSuccess') });
        setIsModalOpen(false);
        fetchWebhooks();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.webhook.toast.deleteFailed') });
    }
  };

  const handleHeadersChange = useCallback((key, value) => {
    setEditForm((prev) => {
      const newHeaders = { ...prev.headers };
      if (value === '' || value.trim() === '') {
        delete newHeaders[key];
      } else {
        newHeaders[key] = value.trim();
      }
      return { ...prev, headers: newHeaders };
    });
  }, []);

  const handleAddHeader = useCallback(() => {
    setEditForm((prev) => {
      // 生成一个唯一的键名
      let counter = 1;
      let newKey = `header_${counter}`;
      while (Object.prototype.hasOwnProperty.call(prev.headers, newKey)) {
        counter++;
        newKey = `header_${counter}`;
      }
      const newHeaders = { ...prev.headers, [newKey]: '' };
      return { ...prev, headers: newHeaders };
    });
  }, []);

  // 防抖更新请求头key
  const debouncedUpdateHeaderKey = useCallback(
    (oldKey, newKey, value) => {
      // 如果新key为空或与旧key相同，不进行更新
      if (!newKey || newKey.trim() === '' || newKey === oldKey) {
        return;
      }

      const trimmedNewKey = newKey.trim();

      // 检查新key是否已存在（除了当前key）
      if (Object.keys(editForm.headers).some((key) => key !== oldKey && key === trimmedNewKey)) {
        return;
      }

      const newHeaders = { ...editForm.headers };
      delete newHeaders[oldKey];
      newHeaders[trimmedNewKey] = value;
      setEditForm((prev) => ({ ...prev, headers: newHeaders }));
    },
    [editForm.headers]
  );

  const handleUpdateHeaderKey = (oldKey, newKey, value) => {
    // 使用防抖更新
    const timeoutId = setTimeout(() => {
      debouncedUpdateHeaderKey(oldKey, newKey, value);
    }, 300);

    // 清理之前的定时器
    setHeaderKeyUpdates((prev) => {
      if (prev[oldKey]) {
        clearTimeout(prev[oldKey]);
      }
      return { ...prev, [oldKey]: timeoutId };
    });
  };

  const handleEventToggle = (event) => {
    const newEvents = editForm.events.includes(event)
      ? editForm.events.filter((e) => e !== event)
      : [...editForm.events, event];
    setEditForm({ ...editForm, events: newEvents });
  };

  // 使用useMemo缓存请求头渲染
  const headersList = useMemo(() => {
    return Object.entries(editForm.headers).map(([key, value]) => (
      <div key={key} className="flex gap-2 items-center">
        <Input
          type="text"
          defaultValue={key}
          onBlur={(e) => handleUpdateHeaderKey(key, e.target.value, value)}
          placeholder={t('admin.webhook.form.headerKeyPlaceholder')}
          className="flex-1"
        />
        <Input
          type="text"
          value={value}
          onChange={(e) => handleHeadersChange(key, e.target.value)}
          placeholder={t('admin.webhook.form.headerValuePlaceholder')}
          className="flex-1"
        />
        <Button
          variant="ghost"
          size="icon"
          className="!bg-transparent !text-red-400 hover:!text-red-300"
          onClick={() => handleHeadersChange(key, '')}
        >
          <IconTrash size={18} />
        </Button>
      </div>
    ));
  }, [editForm.headers, handleHeadersChange, handleUpdateHeaderKey, t]);

  const renderModalContent = () => {
    if (mode === 'delete') {
      return (
        <div className="text-center">
          <p className="text-neutral-300 mb-4">
            {t('admin.webhook.modal.deletePrompt')}{' '}
            <span className="font-semibold text-red-400">{selectedWebhook?.name}</span>?
          </p>
          <p className="text-neutral-400 text-sm">{t('admin.webhook.modal.deleteWarning')}</p>
        </div>
      );
    }

    return (
      <div className="space-y-4">
        {/* 基本信息 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.form.nameLabel')}
            </label>
            <Input
              type="text"
              value={editForm.name}
              onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
              placeholder={t('admin.webhook.form.namePlaceholder')}
              fullWidth
              required={mode === 'create'}
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.form.methodLabel')}
            </label>
            <select
              value={editForm.method}
              onChange={(e) => setEditForm({ ...editForm, method: e.target.value })}
              className="w-full px-3 py-1.5 bg-neutral-800 border border-neutral-700 rounded-lg text-neutral-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required={mode === 'create'}
            >
              <option value="GET">GET</option>
              <option value="POST">POST</option>
            </select>
          </div>
        </div>

        {/* URL */}
        <div>
          <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.webhook.form.urlLabel')}</label>
          <Input
            type="text"
            value={editForm.url}
            onChange={(e) => setEditForm({ ...editForm, url: e.target.value })}
            placeholder={t('admin.webhook.form.urlPlaceholder')}
            fullWidth
            required={mode === 'create'}
          />
        </div>

        {/* 超时和重试配置 */}
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.form.timeoutLabel')}
            </label>
            <Input
              type="number"
              value={editForm.timeout}
              onChange={(e) => setEditForm({ ...editForm, timeout: parseInt(e.target.value) })}
              fullWidth
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.form.retryLabel')}
            </label>
            <Input
              type="number"
              value={editForm.retry}
              onChange={(e) => setEditForm({ ...editForm, retry: parseInt(e.target.value) })}
              fullWidth
            />
          </div>
        </div>

        {/* 请求头配置 */}
        <div>
          <div className="flex justify-between items-center mb-3">
            <label className="block text-neutral-300 text-sm font-medium">{t('admin.webhook.form.headers')}</label>
            <Button
              variant="primary"
              size="sm"
              align="icon-left"
              icon={<IconPlus size={16} />}
              onClick={handleAddHeader}
            >
              {t('admin.webhook.form.addHeader')}
            </Button>
          </div>
          <div className="space-y-2">{headersList}</div>
        </div>

        {/* 事件选择 */}
        <div>
          <label className="block text-neutral-300 text-sm font-medium mb-2">{t('admin.webhook.form.events')}</label>
          <div className="bg-neutral-800 border border-neutral-700 rounded-lg p-3 max-h-40 overflow-y-auto">
            {availableEvents.length > 0 ? (
              <div className="grid grid-cols-2 gap-2">
                {availableEvents.map((event) => (
                  <label key={event} className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={editForm.events.includes(event)}
                      onChange={() => handleEventToggle(event)}
                      className="text-blue-500"
                    />
                    <span className="text-neutral-300 text-sm">{event}</span>
                  </label>
                ))}
              </div>
            ) : (
              <p className="text-neutral-500 text-sm">{t('admin.webhook.form.noEvents')}</p>
            )}
          </div>
          <p className="text-neutral-400 text-xs mt-1">{t('admin.webhook.form.eventsHint')}</p>
        </div>

        {/* 状态 */}
        <div className="flex items-center">
          <input
            type="checkbox"
            id="on"
            checked={editForm.on}
            onChange={(e) => setEditForm({ ...editForm, on: e.target.checked })}
            className="mr-2"
          />
          <label htmlFor="on" className="text-neutral-300">
            {t('admin.webhook.form.enable')}
          </label>
        </div>
      </div>
    );
  };

  const renderHistoryDetailContent = () => {
    if (!selectedHistory) return null;

    return (
      <div className="space-y-4">
        {/* 基本信息 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.history.webhook')}
            </label>
            <div className="text-neutral-50">{selectedHistory.webhook}</div>
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.history.event')}
            </label>
            <div className="text-neutral-50">{selectedHistory.event}</div>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.history.status')}
            </label>
            <div>
              {selectedHistory.success ? (
                <span className="text-green-400">{t('admin.webhook.history.statusSuccess')}</span>
              ) : (
                <span className="text-red-400">{t('admin.webhook.history.statusFailed')}</span>
              )}
            </div>
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.history.responseCode')}
            </label>
            <div className="text-neutral-50">{selectedHistory.resp || '-'}</div>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.history.duration')}
            </label>
            <div className="text-neutral-50">{selectedHistory.duration ? `${selectedHistory.duration}ms` : '-'}</div>
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.history.retry')}
            </label>
            <div className="text-neutral-50">{selectedHistory.retry}</div>
          </div>
        </div>

        {/* 错误信息 */}
        {selectedHistory.error && (
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.webhook.history.error')}
            </label>
            <div className="bg-red-900/20 border border-red-700 rounded-lg p-3">
              <pre className="text-red-300 text-sm whitespace-pre-wrap font-mono">{selectedHistory.error}</pre>
            </div>
          </div>
        )}
      </div>
    );
  };

  const renderModalFooter = () => {
    return (
      <>
        <ModalButton onClick={() => setIsModalOpen(false)}>{t('common.cancel')}</ModalButton>
        <ModalButton
          variant={mode === 'delete' ? 'danger' : 'primary'}
          onClick={
            mode === 'create' ? handleCreateWebhook : mode === 'edit' ? handleUpdateWebhook : handleDeleteWebhook
          }
        >
          {mode === 'create' ? t('common.create') : mode === 'edit' ? t('common.save') : t('common.delete')}
        </ModalButton>
      </>
    );
  };

  return (
    <>
      {/* 标签页切换 */}
      <div className="w-full mx-auto mb-6">
        <div className="flex border-b border-neutral-700">
          <button
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === 'webhook'
                ? 'text-blue-400 border-b-2 border-blue-400'
                : 'text-neutral-400 hover:text-neutral-300'
            }`}
            onClick={() => setActiveTab('webhook')}
          >
            {t('admin.webhook.tabs.webhook')}
          </button>
          <button
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === 'history'
                ? 'text-blue-400 border-b-2 border-blue-400'
                : 'text-neutral-400 hover:text-neutral-300'
            }`}
            onClick={handleViewAllHistory}
          >
            {t('admin.webhook.tabs.history')}
          </button>
        </div>
      </div>

      {/* Webhook配置管理 */}
      {activeTab === 'webhook' && (
        <AdminWebhook
          webhooks={webhooks}
          totalCount={totalCount}
          currentPage={currentPage}
          pageSize={10}
          loading={false}
          onPageChange={setCurrentPage}
          onCreateWebhook={handleCreateClick}
          onEditWebhook={handleEditClick}
          onDeleteWebhook={handleDeleteClick}
          onViewHistory={handleViewHistory}
          onWebhookClick={handleWebhookClick}
        />
      )}

      {/* Webhook历史记录 */}
      {activeTab === 'history' && (
        <AdminWebhookHistory
          webhookHistory={webhookHistory}
          totalCount={historyTotalCount}
          currentPage={historyCurrentPage}
          pageSize={20}
          loading={false}
          onPageChange={setHistoryCurrentPage}
          onViewDetail={handleHistoryClick}
          onHistoryClick={handleHistoryClick}
          webhookName={currentWebhookId ? webhooks.find((w) => w.id === currentWebhookId)?.name : null}
        />
      )}

      {/* Webhook配置模态框 */}
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={
          mode === 'create'
            ? t('admin.webhook.modal.createTitle')
            : mode === 'edit'
              ? t('admin.webhook.modal.editTitle')
              : t('admin.webhook.modal.deleteTitle')
        }
        size={mode !== 'delete' ? 'lg' : 'sm'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>

      {/* 历史记录详情模态框 */}
      <Modal
        isOpen={isHistoryDetailModalOpen}
        onClose={() => setIsHistoryDetailModalOpen(false)}
        title={t('admin.webhook.history.detailTitle')}
        size="lg"
        footer={
          <ModalButton onClick={() => setIsHistoryDetailModalOpen(false)}>
            {t('admin.webhook.actions.close')}
          </ModalButton>
        }
      >
        {renderHistoryDetailContent()}
      </Modal>
    </>
  );
}

export default WebhookManagement;
