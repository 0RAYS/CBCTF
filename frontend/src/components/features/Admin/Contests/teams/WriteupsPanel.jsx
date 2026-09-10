import { useTranslation } from 'react-i18next';
import { IconDownload } from '@tabler/icons-react';
import { Button, Card, List, Pagination, Spinner } from '../../../../common';
import { formatFileSize, TEAM_DETAIL_PAGE_SIZE } from './teamDetailData.js';

export default function WriteupsPanel({
  writeups = [],
  count = 0,
  page = 1,
  loading,
  onPageChange,
  onUserClick,
  onDownload,
}) {
  const { t, i18n } = useTranslation();
  const columns = [
    ['date', 'submittedAt'],
    ['filename', 'filename'],
    ['size', 'size'],
    ['hash', 'hash'],
    ['user_id', 'uploader'],
    ['actions', 'actions'],
  ].map(([key, label]) => ({ key, label: t(`admin.contests.teamDetail.writeups.columns.${label}`) }));
  const renderCell = (writeup, { key }) => {
    if (key === 'date') return new Date(writeup.date).toLocaleString(i18n.language);
    if (key === 'size') return formatFileSize(writeup.size);
    if (key === 'hash')
      return (
        <span className="block max-w-xs truncate" title={writeup.hash}>
          {writeup.hash}
        </span>
      );
    if (key === 'user_id' && onUserClick)
      return (
        <Button variant="ghost" size="sm" onClick={() => onUserClick(writeup.user_id)}>
          {writeup.user_id}
        </Button>
      );
    if (key === 'actions')
      return (
        <Button
          variant="ghost"
          size="icon"
          onClick={() => onDownload(writeup)}
          title={t('admin.contests.teamDetail.writeups.actions.download')}
        >
          <IconDownload size={18} />
        </Button>
      );
    return writeup[key] ?? '-';
  };
  return (
    <section>
      <h2 className="text-xl font-mono text-neutral-50 mb-4">{t('admin.contests.teamDetail.sections.writeups')}</h2>
      {loading ? (
        <Card className="flex justify-center p-8">
          <Spinner />
        </Card>
      ) : (
        <List
          className="[&_table]:min-w-[800px]"
          columns={columns}
          data={writeups}
          renderCell={renderCell}
          empty={!writeups.length}
          emptyContent={t('admin.contests.teamDetail.empty.writeups')}
          footer={
            <Pagination
              total={Math.ceil(count / TEAM_DETAIL_PAGE_SIZE)}
              current={page}
              pageSize={TEAM_DETAIL_PAGE_SIZE}
              onChange={onPageChange}
              showTotal
              totalItems={count}
            />
          }
        />
      )}
    </section>
  );
}
