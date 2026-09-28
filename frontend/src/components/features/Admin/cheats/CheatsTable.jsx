import { useTranslation } from 'react-i18next';
import { IconEdit, IconEye, IconTrash } from '@tabler/icons-react';
import { EmptyState, List, Pagination } from '../../../common';
import { CheatIp, CheatModels, CheatStatus, CheatVerdict } from './CheatEvidence';

export default function CheatsTable({
  cheats,
  total,
  page,
  pageSize,
  loading,
  onPageChange,
  onAction,
  openUserDetail,
  openTeamDetail,
}) {
  const { t, i18n } = useTranslation();
  const columns = [
    ['id', 'id', '5%'],
    ['model', 'model', '15%'],
    ['type', 'type', '10%'],
    ['reason_type', 'reasonType', '12%'],
    ['reason', 'reason', '14%'],
    ['ip', 'ip', '10%'],
    ['checked', 'status', '10%'],
    ['time', 'time', '14%'],
  ].map(([key, label, width]) => ({
    key,
    label: t(`admin.contests.cheats.columns.${label}`),
    width,
  }));

  function renderCell(item, column) {
    switch (column.key) {
      case 'model':
        return <CheatModels model={item.model} openUserDetail={openUserDetail} openTeamDetail={openTeamDetail} />;
      case 'type':
        return <CheatVerdict type={item.type} />;
      case 'reason_type':
        return (
          <span className="text-neutral-300 font-mono text-xs">
            {item.reason_type ? t(`admin.contests.cheats.reasonTypes.${item.reason_type}`, item.reason_type) : '-'}
          </span>
        );
      case 'reason':
        return (
          <span className="max-w-50 truncate block" title={item.reason}>
            {item.reason || '-'}
          </span>
        );
      case 'ip':
        return <CheatIp ip={item.ip} />;
      case 'checked':
        return <CheatStatus checked={item.checked} />;
      case 'time':
        return item.time ? new Date(item.time).toLocaleString(i18n.language || 'en-US') : '-';
      default:
        return item[column.key] || '-';
    }
  }
  return (
    <List
      columns={columns}
      data={cheats}
      renderCell={renderCell}
      onRowClick={(item) => onAction('detail', item)}
      getRowActions={(item) => [
        {
          key: 'detail',
          label: t('admin.contests.cheats.modals.detailTitle'),
          icon: <IconEye size={18} />,
          onClick: () => onAction('detail', item),
        },
        {
          key: 'edit',
          label: t('common.edit'),
          icon: <IconEdit size={18} />,
          onClick: () => onAction('edit', item),
        },
        {
          key: 'delete',
          label: t('common.delete'),
          icon: <IconTrash size={18} />,
          danger: true,
          onClick: () => onAction('delete', item),
        },
      ]}
      loading={loading}
      empty={!loading && cheats.length === 0}
      emptyContent={<EmptyState title={t('admin.contests.cheats.empty')} />}
      footer={
        total > pageSize && (
          <Pagination
            total={Math.ceil(total / pageSize)}
            current={page}
            onChange={onPageChange}
            showTotal
            totalItems={total}
          />
        )
      }
    />
  );
}
