import { useTranslation } from 'react-i18next';
import { Pagination } from '../../common';
import ScoreboardRanking from './ScoreboardRanking.jsx';

/** 管理后台积分榜组件 */
function AdminRanking({
  teams = [],
  currentPage = 1,
  pageSize = 6,
  totalCount = 0,
  onPageChange,
  onRowClick,
  loading,
}) {
  const { t, i18n } = useTranslation();

  return (
    <ScoreboardRanking
      teams={teams}
      locale={i18n.language || 'en-US'}
      labels={{
        rank: t('common.scoreboard.headers.rank'),
        team: t('common.scoreboard.headers.team'),
        score: t('common.scoreboard.headers.score'),
        challenges: t('common.scoreboard.headers.challenges'),
        lastSubmit: t('common.scoreboard.headers.lastSubmit'),
        total: t('common.scoreboard.total'),
      }}
      emptyMessage={t('common.noData')}
      loading={loading}
      onRowClick={onRowClick}
      footer={
        totalCount > pageSize ? (
          <div className="mt-6">
            <Pagination
              total={Math.ceil(totalCount / pageSize)}
              current={currentPage}
              onChange={onPageChange}
              showTotal
              totalItems={totalCount}
            />
          </div>
        ) : null
      }
    />
  );
}

export default AdminRanking;
