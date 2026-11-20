import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getUserRanking } from '../../api/user';
import { Avatar, List, Pagination } from '../../components/common';
import Checkbox from '../../components/common/Checkbox';
import { useHiddenVisibility } from '../../hooks/useHiddenVisibility.js';
import { toast } from '../../utils/toast.js';

const PAGE_SIZE = 20;

function UserLeaderboardPage() {
  const { t, i18n } = useTranslation();
  const [currentPage, setCurrentPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showHidden, setShowHidden] = useHiddenVisibility('cbctf.leaderboard.showHiddenUsers');

  useEffect(() => {
    let ignore = false;
    setLoading(true);
    getUserRanking({
      limit: PAGE_SIZE,
      offset: (currentPage - 1) * PAGE_SIZE,
      show_hidden: showHidden ? 1 : 0,
    })
      .then((response) => {
        if (ignore || response.code !== 200) return;
        setUsers(response.data.users || []);
        setTotalCount(response.data.count || 0);
      })
      .catch((error) => {
        if (!ignore) toast.danger({ description: error.message || t('game.leaderboard.toast.fetchFailed') });
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });
    return () => {
      ignore = true;
    };
  }, [currentPage, showHidden, t]);

  const columns = [
    { key: 'rank', label: t('common.scoreboard.headers.rank'), width: '12%' },
    { key: 'name', label: t('common.scoreboard.headers.team'), width: '48%' },
    { key: 'score', label: t('common.scoreboard.headers.score'), width: '20%' },
    { key: 'solved', label: t('common.solved'), width: '20%' },
  ];

  const renderCell = (user, column) => {
    switch (column.key) {
      case 'rank':
        return <span className="tabular-nums">{user.rank}</span>;
      case 'name':
        return (
          <div className="flex min-w-0 items-center gap-2">
            <Avatar src={user.picture} name={user.name} size="xs" className="shrink-0" />
            <span className="min-w-0 truncate text-neutral-50">{user.name}</span>
            {user.hidden && (
              <span className="shrink-0 rounded bg-amber-400/15 px-1.5 py-0.5 text-[10px] text-amber-300">
                {t('common.hidden')}
              </span>
            )}
          </div>
        );
      case 'score':
        return <span className="tabular-nums">{(user.score ?? 0).toLocaleString(i18n.language || 'en-US')}</span>;
      case 'solved':
        return <span className="tabular-nums">{user.solved ?? 0}</span>;
      default:
        return null;
    }
  };

  return (
    <div className="contest-container mx-auto min-w-0 space-y-4">
      <div className="flex flex-wrap justify-between items-center gap-3">
        <h1 className="text-xl sm:text-3xl font-mono text-neutral-50 tracking-wider">{t('nav.leaderboard')}</h1>
        <Checkbox
          checked={showHidden}
          onChange={(event) => {
            setShowHidden(event.target.checked);
            setCurrentPage(1);
          }}
          label={t('game.leaderboard.showHiddenUsers')}
          className="text-xs text-neutral-400"
        />
      </div>
      <List
        columns={columns}
        data={users}
        renderCell={renderCell}
        loading={loading}
        minWidth={640}
        footer={
          totalCount > PAGE_SIZE ? (
            <div className="flex justify-center">
              <Pagination
                current={currentPage}
                total={Math.ceil(totalCount / PAGE_SIZE)}
                onChange={setCurrentPage}
                showTotal={true}
                totalItems={totalCount}
                animate
              />
            </div>
          ) : null
        }
      />
    </div>
  );
}

export default UserLeaderboardPage;
