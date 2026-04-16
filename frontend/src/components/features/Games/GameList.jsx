import { useState } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { IconChevronLeft, IconChevronRight } from '@tabler/icons-react';
import Button from '../../common/Button';
import { useTranslation } from 'react-i18next';

const PAGE_SIZE = 5;

function GameCard({ game, onGameAction, user }) {
  const { t } = useTranslation();

  const getStatusStyle = (status) => {
    switch (status) {
      case 'upcoming':
        return 'text-yellow-400 border-yellow-400/40';
      case 'active':
        return 'text-green-400 border-green-400/40';
      case 'ended':
        return 'text-neutral-400 border-neutral-400/40';
      default:
        return 'text-neutral-300 border-neutral-300/40';
    }
  };

  const handleAction = () => {
    const action = game.status === 'ended' ? 'view' : 'join';
    onGameAction?.(game.id, action);
  };

  return (
    <div className="border border-neutral-300/20 rounded-md overflow-hidden bg-black/30 hover:border-neutral-300/50 transition-colors duration-200 flex flex-col sm:flex-row">
      {/* 封面图（保持原比例） */}
      {game.image && (
        <div className="sm:w-48 shrink-0 bg-neutral-900/50">
          <img src={game.image} alt={game.title} className="w-full h-full object-contain" />
        </div>
      )}

      <div className="flex-1 p-5 flex flex-col sm:flex-row sm:items-center gap-4">
        {/* 左侧: 标题 + 描述 */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3 mb-2 flex-wrap">
            <h3 className="text-base font-mono text-neutral-50 tracking-wide">{game.title}</h3>
            <span className={`px-2 py-0.5 border rounded text-xs shrink-0 ${getStatusStyle(game.status)}`}>
              {t(`game.status.${game.status}`)}
            </span>
          </div>
          <div className="text-neutral-400 text-sm line-clamp-2 [&>p]:m-0 prose prose-invert prose-sm max-w-none">
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{game.description || ''}</ReactMarkdown>
          </div>
        </div>

        {/* 右侧: 时间 + 操作按钮 */}
        <div className="flex flex-col items-start sm:items-end gap-3 shrink-0">
          <div className="text-xs font-mono text-neutral-400 space-y-1 sm:text-right">
            <div>
              <span className="text-neutral-500">{t('common.start')}: </span>
              {new Date(game.startTime).toLocaleString()}
            </div>
            <div>
              <span className="text-neutral-500">{t('common.end')}: </span>
              {new Date(game.endTime).toLocaleString()}
            </div>
          </div>
          <Button variant="primary" size="sm" onClick={handleAction}>
            {user?.user
              ? game.status === 'ended'
                ? t('game.actions.results')
                : t('game.actions.joinNow')
              : game.status === 'ended'
                ? t('game.actions.ended')
                : t('auth.login')}
          </Button>
        </div>
      </div>
    </div>
  );
}

function GameList({ games = [], onGameAction, user }) {
  const [page, setPage] = useState(1);

  const totalPages = Math.ceil(games.length / PAGE_SIZE);
  const pageGames = games.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

  return (
    <div className="w-full max-w-[1200px] mx-auto">
      <AnimatePresence mode="wait">
        <motion.div
          key={page}
          initial={{ opacity: 0, y: 6 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -6 }}
          transition={{ duration: 0.18 }}
          className="space-y-3"
        >
          {pageGames.map((game) => (
            <GameCard key={game.id} game={game} onGameAction={onGameAction} user={user} />
          ))}
        </motion.div>
      </AnimatePresence>

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2 pt-6">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className={`w-9 h-9 border rounded-md flex items-center justify-center transition-colors duration-200
              ${
                page === 1
                  ? 'border-neutral-700 text-neutral-600 cursor-not-allowed'
                  : 'border-neutral-500 text-neutral-300 hover:border-geek-400 hover:text-geek-400'
              }`}
          >
            <IconChevronLeft size={15} />
          </button>

          {Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
            <button
              key={p}
              onClick={() => setPage(p)}
              className={`w-9 h-9 border rounded-md flex items-center justify-center font-mono text-sm transition-colors duration-200
                ${
                  p === page
                    ? 'border-geek-400 text-geek-400 bg-geek-400/10'
                    : 'border-neutral-600 text-neutral-400 hover:border-neutral-400 hover:text-neutral-200'
                }`}
            >
              {p}
            </button>
          ))}

          <button
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            className={`w-9 h-9 border rounded-md flex items-center justify-center transition-colors duration-200
              ${
                page === totalPages
                  ? 'border-neutral-700 text-neutral-600 cursor-not-allowed'
                  : 'border-neutral-500 text-neutral-300 hover:border-geek-400 hover:text-geek-400'
              }`}
          >
            <IconChevronRight size={15} />
          </button>
        </div>
      )}
    </div>
  );
}

export default GameList;
