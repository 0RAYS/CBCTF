import { useState, useCallback } from 'react';
import { toast } from '../utils/toast';

/**
 * Hook for managing CRUD modal state and operations.
 *
 * @param {Object} options
 * @param {Object} options.defaultForm - Default form values for create mode
 * @param {Function} options.createApi - API function for create (receives form data)
 * @param {Function} options.updateApi - API function for update (receives id, form data)
 * @param {Function} options.deleteApi - API function for delete (receives id)
 * @param {Function} options.onSuccess - Callback after successful operation (e.g., refetch data)
 * @param {Function} [options.beforeCreate] - Transform form data before create API call
 * @param {Function} [options.beforeUpdate] - Transform form data before update API call
 * @param {Function} [options.itemToForm] - Convert item to edit form values
 * @param {Object} [options.messages] - Toast messages { createSuccess, createFailed, updateSuccess, updateFailed, deleteSuccess, deleteFailed }
 */
export function useCRUDModal({
  defaultForm = {},
  createApi,
  updateApi,
  deleteApi,
  onSuccess,
  beforeCreate,
  beforeUpdate,
  itemToForm,
  messages = {},
}) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [mode, setMode] = useState('create'); // 'create' | 'edit' | 'delete'
  const [selectedItem, setSelectedItem] = useState(null);
  const [editForm, setEditForm] = useState({ ...defaultForm });

  const openCreate = useCallback(() => {
    setEditForm({ ...defaultForm });
    setSelectedItem(null);
    setMode('create');
    setIsModalOpen(true);
  }, [defaultForm]);

  const openEdit = useCallback(
    (item) => {
      setSelectedItem(item);
      setEditForm(itemToForm ? itemToForm(item) : { ...item });
      setMode('edit');
      setIsModalOpen(true);
    },
    [itemToForm]
  );

  const openDelete = useCallback((item) => {
    setSelectedItem(item);
    setMode('delete');
    setIsModalOpen(true);
  }, []);

  const closeModal = useCallback(() => {
    setIsModalOpen(false);
  }, []);

  const handleSubmit = useCallback(async () => {
    try {
      let response;
      if (mode === 'create' && createApi) {
        const data = beforeCreate ? beforeCreate(editForm) : editForm;
        response = await createApi(data);
        if (response.code === 200) {
          toast.success({ description: messages.createSuccess });
          setIsModalOpen(false);
          onSuccess?.();
        }
      } else if (mode === 'edit' && updateApi) {
        const data = beforeUpdate ? beforeUpdate(editForm, selectedItem) : editForm;
        response = await updateApi(selectedItem.id, data);
        if (response.code === 200) {
          toast.success({ description: messages.updateSuccess });
          setIsModalOpen(false);
          onSuccess?.();
        }
      } else if (mode === 'delete' && deleteApi) {
        response = await deleteApi(selectedItem.id);
        if (response.code === 200) {
          toast.success({ description: messages.deleteSuccess });
          setIsModalOpen(false);
          onSuccess?.();
        }
      }
    } catch (error) {
      const fallback =
        mode === 'create' ? messages.createFailed : mode === 'edit' ? messages.updateFailed : messages.deleteFailed;
      toast.danger({ description: error.message || fallback });
    }
  }, [mode, editForm, selectedItem, createApi, updateApi, deleteApi, beforeCreate, beforeUpdate, onSuccess, messages]);

  return {
    isModalOpen,
    mode,
    selectedItem,
    editForm,
    setEditForm,
    openCreate,
    openEdit,
    openDelete,
    closeModal,
    handleSubmit,
  };
}
