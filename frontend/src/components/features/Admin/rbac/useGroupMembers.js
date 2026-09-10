import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  assignUserToGroup,
  removeUserFromGroup,
  getGroupUsers,
  getGroupAvailableUsers,
} from '../../../../api/admin/rbac';
import { useDebounce } from '../../../../hooks';
import { toast } from '../../../../utils/toast';
import { lastAvailablePage, successfulAssignmentIds, toggleVisibleCandidates } from './payloads';

export default function useGroupMembers(group, onChanged) {
  const { t } = useTranslation();
  const [users, setUsers] = useState([]);
  const [userCount, setUserCount] = useState(0);
  const [userPage, setUserPage] = useState(1);
  const [loadingUsers, setLoadingUsers] = useState(true);
  const [candidates, setCandidates] = useState([]);
  const [candidateCount, setCandidateCount] = useState(0);
  const [candidatePage, setCandidatePage] = useState(1);
  const [loadingCandidates, setLoadingCandidates] = useState(true);
  const [queries, setQueries] = useState({ name: '', email: '', description: '' });
  const debouncedQueries = useDebounce(queries, 300);
  const [selectedIds, setSelectedIds] = useState([]);
  const [revision, setRevision] = useState(0);
  const [pending, setPending] = useState(false);
  const busy = useRef(false);
  const pageSize = 10;

  useEffect(() => {
    let cancelled = false;
    setLoadingUsers(true);
    getGroupUsers(group.id, { limit: pageSize, offset: (userPage - 1) * pageSize })
      .then((response) => {
        if (response.code !== 200) throw new Error(t('admin.rbac.groups.toast.fetchUsersFailed'));
        if (cancelled) return;
        const count = response.data.count || 0;
        setUsers(response.data.users || []);
        setUserCount(count);
        setUserPage(lastAvailablePage(count, pageSize, userPage));
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchUsersFailed') });
      })
      .finally(() => {
        if (!cancelled) setLoadingUsers(false);
      });
    return () => {
      cancelled = true;
    };
  }, [group.id, userPage, revision, t]);

  useEffect(() => {
    let cancelled = false;
    const params = { limit: pageSize, offset: (candidatePage - 1) * pageSize };
    Object.entries(debouncedQueries).forEach(([key, value]) => {
      if (value.trim()) params[key] = value.trim();
    });
    setLoadingCandidates(true);
    getGroupAvailableUsers(group.id, params)
      .then((response) => {
        if (response.code !== 200) throw new Error(t('admin.rbac.groups.toast.fetchCandidatesFailed'));
        if (cancelled) return;
        const count = response.data.count || 0;
        const page = lastAvailablePage(count, pageSize, candidatePage);
        setCandidates(page === candidatePage ? response.data.users || [] : []);
        setCandidateCount(count);
        setCandidatePage(page);
      })
      .catch((error) => {
        if (cancelled) return;
        toast.danger({ description: error.message || t('admin.rbac.groups.toast.fetchCandidatesFailed') });
        setCandidates([]);
        setCandidateCount(0);
      })
      .finally(() => {
        if (!cancelled) setLoadingCandidates(false);
      });
    return () => {
      cancelled = true;
    };
  }, [group.id, candidatePage, debouncedQueries, revision, t]);

  function changeQuery(key, value) {
    setQueries((prev) => ({ ...prev, [key]: value }));
    setCandidatePage(1);
  }

  async function removeUser(user) {
    if (busy.current) return;
    busy.current = true;
    setPending(true);
    try {
      const response = await removeUserFromGroup(group.id, { user_id: user.id });
      if (response.code !== 200) throw new Error(t('admin.rbac.groups.toast.removeUserFailed'));
      toast.success({ description: t('admin.rbac.groups.toast.removeUserSuccess') });
      setUserPage((page) => (page === userPage && users.length === 1 && page > 1 ? page - 1 : page));
      setRevision((value) => value + 1);
      onChanged();
    } catch (error) {
      toast.danger({ description: error.message || t('admin.rbac.groups.toast.removeUserFailed') });
    } finally {
      busy.current = false;
      setPending(false);
    }
  }

  async function assignSelected() {
    if (busy.current || selectedIds.length === 0) return;
    busy.current = true;
    setPending(true);
    const ids = [...selectedIds];
    try {
      const results = await Promise.allSettled(ids.map((userId) => assignUserToGroup(group.id, { user_id: userId })));
      const successfulIds = successfulAssignmentIds(ids, results);
      if (successfulIds.length > 0) {
        // Only successful assignments leave the cross-page selection; failures remain retryable.
        setSelectedIds((prev) => prev.filter((id) => !successfulIds.includes(id)));
        setRevision((value) => value + 1);
        onChanged();
      }
      if (successfulIds.length === ids.length) {
        toast.success({
          description:
            successfulIds.length === 1
              ? t('admin.rbac.groups.toast.assignUserSuccess')
              : t('admin.rbac.groups.toast.assignUsersSuccess', { count: successfulIds.length }),
        });
      } else if (successfulIds.length > 0) {
        toast.warning({
          description: t('admin.rbac.groups.toast.assignUsersPartial', {
            success: successfulIds.length,
            total: ids.length,
          }),
        });
      } else {
        toast.danger({ description: t('admin.rbac.groups.toast.assignUserFailed') });
      }
    } finally {
      busy.current = false;
      setPending(false);
    }
  }

  return {
    users,
    userCount,
    userPage,
    setUserPage,
    loadingUsers,
    candidates,
    candidateCount,
    candidatePage,
    setCandidatePage,
    loadingCandidates,
    queries,
    changeQuery,
    selectedIds,
    pageSize,
    pending,
    searching: Object.values(debouncedQueries).some((value) => value.trim()),
    allVisibleSelected: candidates.length > 0 && candidates.every((user) => selectedIds.includes(user.id)),
    toggleCandidate: (id) =>
      setSelectedIds((prev) => (prev.includes(id) ? prev.filter((value) => value !== id) : [...prev, id])),
    toggleAllCandidates: () => setSelectedIds((prev) => toggleVisibleCandidates(prev, candidates)),
    removeUser,
    assignSelected,
  };
}
