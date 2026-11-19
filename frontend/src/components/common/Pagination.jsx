import { motion } from 'motion/react';
import { IconChevronsLeft, IconChevronsRight } from '@tabler/icons-react';
import Button from './Button';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

/** 统一分页组件 */
function Pagination({
  current = 1,
  total = 1,
  onChange,
  showTotal = false,
  totalItems,
  simple = false,
  showJumpTo = false,
  showEdgeButtons = true,
  animate = false,
  className = '',
}) {
  const { t } = useTranslation();
  const [jumpPage, setJumpPage] = useState('');
  const liveRef = useRef(null);

  useEffect(() => {
    if (liveRef.current && total > 0) {
      liveRef.current.textContent = t('common.pagination.liveAnnounce', {
        current,
        total,
        defaultValue: `Page ${current} of ${total}`,
      });
    }
  }, [current, total, t]);

  const calculatePageRange = () => {
    const range = [];
    const showItems = 5;
    const side = Math.floor(showItems / 2);

    if (total <= 0) return [];

    let start = Math.max(current - side, 1);
    let end = Math.min(start + showItems - 1, total);

    if (end - start + 1 < showItems) {
      start = Math.max(end - showItems + 1, 1);
    }

    if (start > 2) {
      range.push(1);
      range.push('prev-dot');
    } else if (start === 2) {
      range.push(1);
    }

    for (let i = start; i <= end; i++) {
      range.push(i);
    }

    if (end < total - 1) {
      range.push('next-dot');
      range.push(total);
    } else if (end === total - 1) {
      range.push(total);
    }

    return range;
  };

  const renderPageButton = (page) => {
    if (typeof page === 'string') {
      return (
        <span
          key={page}
          className="w-8 h-8 flex items-center justify-center text-neutral-400 font-mono"
          aria-hidden="true"
        >
          ...
        </span>
      );
    }

    const isCurrent = current === page;
    return (
      <Button
        key={page}
        variant={isCurrent ? 'primary' : 'ghost'}
        size="icon"
        className={`!w-10 !h-10 ${isCurrent ? '' : '!bg-transparent'}`}
        onClick={() => onChange?.(page)}
        disabled={page <= 0}
        animate={false}
        aria-label={t('common.pagination.goToPage', { page, defaultValue: `Go to page ${page}` })}
        aria-current={isCurrent ? 'page' : undefined}
      >
        {page}
      </Button>
    );
  };

  const renderNavigationButton = (type) => {
    const isNext = type === 'next';
    const disabled = isNext ? current >= total || total <= 0 : current <= 1;

    return (
      <Button
        variant="ghost"
        size="sm"
        disabled={disabled}
        onClick={() => !disabled && onChange?.(isNext ? current + 1 : current - 1)}
        className={`!bg-neutral-900/50 ${disabled ? '!border-neutral-300/10 !text-neutral-600' : ''}`}
        animate={false}
        aria-label={
          isNext
            ? t('common.pagination.nextPage', { defaultValue: 'Next page' })
            : t('common.pagination.prevPage', { defaultValue: 'Previous page' })
        }
      >
        {isNext ? t('common.next') : t('common.previous')}
      </Button>
    );
  };

  const renderEdgeButtons = () => {
    const disabledFirst = current <= 1;
    const disabledLast = current >= total || total <= 0;

    return (
      <>
        <Button
          variant="ghost"
          size="icon"
          disabled={disabledFirst}
          onClick={() => !disabledFirst && onChange?.(1)}
          className={`!w-10 !h-10 !bg-neutral-900/50 ${disabledFirst ? '!border-neutral-300/10 !text-neutral-600' : ''}`}
          animate={false}
          aria-label={t('common.pagination.firstPage', { defaultValue: 'First page' })}
        >
          <IconChevronsLeft size={16} />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          disabled={disabledLast}
          onClick={() => !disabledLast && onChange?.(total)}
          className={`!w-10 !h-10 !bg-neutral-900/50 ${disabledLast ? '!border-neutral-300/10 !text-neutral-600' : ''}`}
          animate={false}
          aria-label={t('common.pagination.lastPage', { defaultValue: 'Last page' })}
        >
          <IconChevronsRight size={16} />
        </Button>
      </>
    );
  };

  const handleJumpToPage = () => {
    const page = parseInt(jumpPage);
    if (page >= 1 && page <= total) {
      onChange?.(page);
      setJumpPage('');
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleJumpToPage();
    }
  };

  const PaginationContainer = animate ? motion.div : 'div';
  const containerProps = animate
    ? {
        initial: { opacity: 0, y: 20 },
        animate: { opacity: 1, y: 0 },
        transition: { delay: 0.1 },
      }
    : {};

  return (
    <PaginationContainer
      role="navigation"
      aria-label={t('common.pagination.label', { defaultValue: 'Pagination' })}
      className={`flex flex-wrap items-center justify-center gap-2 max-w-full ${className}`}
      {...containerProps}
    >
      {/* 隐藏的 aria-live 区域, 宣告翻页结果 */}
      <span ref={liveRef} aria-live="polite" aria-atomic="true" className="sr-only" />

      {showTotal && totalItems !== undefined && (
        <span className="text-sm font-mono text-neutral-400 mr-4 whitespace-nowrap">
          {t('common.pagination.totalItems', { count: totalItems })}
        </span>
      )}

      {showEdgeButtons && total > 7 && !simple && <div className="flex items-center mr-1">{renderEdgeButtons()}</div>}

      {total > 0 && renderNavigationButton('prev')}

      {!simple && total > 0 && (
        <div className="flex items-center gap-1 md:gap-2 flex-wrap justify-center">
          {calculatePageRange().map(renderPageButton)}
        </div>
      )}

      {total > 0 && renderNavigationButton('next')}

      {showJumpTo && !simple && total > 0 && (
        <div className="flex items-center gap-2 ml-4">
          <span className="text-sm text-neutral-400">{t('common.pagination.jumpTo')}</span>
          <input
            type="number"
            min="1"
            max={total}
            value={jumpPage}
            onChange={(e) => setJumpPage(e.target.value)}
            onKeyPress={handleKeyPress}
            aria-label={t('common.pagination.jumpToLabel', { defaultValue: 'Jump to page number' })}
            className="w-16 px-2 py-1 text-sm bg-neutral-800 border border-neutral-700 rounded text-neutral-300 focus:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70"
            placeholder={t('common.pagination.pagePlaceholder')}
          />
          <span className="text-sm text-neutral-400">{t('common.pagination.page')}</span>
          <Button
            variant="ghost"
            size="sm"
            onClick={handleJumpToPage}
            disabled={!jumpPage || parseInt(jumpPage) < 1 || parseInt(jumpPage) > total}
            className="!bg-black/30"
            animate={false}
          >
            {t('common.pagination.jump')}
          </Button>
        </div>
      )}

      {simple && total > 0 && (
        <span className="text-sm font-mono text-neutral-400 mx-2 whitespace-nowrap">
          {current} / {total}
        </span>
      )}
    </PaginationContainer>
  );
}

export default Pagination;
