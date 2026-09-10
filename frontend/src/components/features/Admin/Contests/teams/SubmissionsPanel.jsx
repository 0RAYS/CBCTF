import { useTranslation } from 'react-i18next';
import { Button, Card, List, Pagination, Spinner, StatusTag } from '../../../../common';
import IpLookup from '../../network/IpLookup';
import { TEAM_DETAIL_PAGE_SIZE } from './teamDetailData.js';

export default function SubmissionsPanel({
  submissions = [],
  count = 0,
  page = 1,
  loading,
  onPageChange,
  onUserClick,
}) {
  const { t, i18n } = useTranslation();
  const columns = [
    ['challenge_id', 'challengeId'],
    ['team_id', 'teamId'],
    ['user_id', 'userId'],
    ['value', 'submittedFlag'],
    ['score', 'score'],
    ['ip', 'ip'],
    ['solved', 'status'],
    ['created_at', 'time'],
  ].map(([key, label]) => ({ key, label: t(`admin.contests.teamDetail.submissions.columns.${label}`) }));
  const renderCell = (submission, { key }) => {
    if (key === 'ip') return <IpLookup ip={submission.ip} />;
    if (key === 'solved')
      return (
        <StatusTag
          type={submission.solved ? 'success' : 'error'}
          text={t(`admin.contests.teamDetail.submissions.status.${submission.solved ? 'correct' : 'incorrect'}`)}
        />
      );
    if (key === 'created_at')
      return submission.created_at ? new Date(submission.created_at).toLocaleString(i18n.language) : '-';
    if (key === 'user_id' && onUserClick)
      return (
        <Button variant="ghost" size="sm" onClick={() => onUserClick(submission.user_id)}>
          {submission.user_id}
        </Button>
      );
    if (key === 'value')
      return (
        <span className="block max-w-xs truncate" title={submission.value}>
          {submission.value}
        </span>
      );
    return submission[key] ?? '-';
  };
  return (
    <section>
      <h2 className="text-xl font-mono text-neutral-50 mb-4">{t('admin.contests.teamDetail.sections.submissions')}</h2>
      {loading ? (
        <Card className="flex justify-center p-8">
          <Spinner />
        </Card>
      ) : (
        <List
          className="[&_table]:min-w-[960px]"
          columns={columns}
          data={submissions}
          renderCell={renderCell}
          empty={!submissions.length}
          emptyContent={t('admin.contests.teamDetail.empty.submissions')}
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
