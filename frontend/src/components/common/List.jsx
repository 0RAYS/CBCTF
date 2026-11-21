import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import Loading from '@/components/common/Loading';
import EmptyState from '@/components/common/EmptyState';
import useContextMenu from '@/components/common/useContextMenu';
import RowActions from '@/components/common/RowActions';
import TruncatedText from '@/components/common/TruncatedText';

const INTERACTIVE_SELECTOR =
  'button, a, input, select, textarea, label, [role="button"], [role="checkbox"], [contenteditable="true"], [data-row-interactive]';

/** 通用列表/表格组件；minWidth 是表格最小宽度, 窄屏时表格区域内横向滚动 */
function List({
  columns = [],
  data = [],
  renderCell,
  onRowClick,
  getRowActions,
  actionsColumn,
  loading = false,
  empty,
  emptyContent,
  footer,
  rowClassName,
  variant = 'default',
  animate = true,
  className = '',
  minWidth = 640,
}) {
  const { t } = useTranslation();
  const contextMenu = useContextMenu(data, getRowActions, loading);
  const columnCount = columns.length + (actionsColumn ? 1 : 0);
  const actionsWidth = actionsColumn?.width || 136;
  // Existing admin columns use percentage weights that often total less than 100.
  // Reserve the action width, then distribute the remaining space proportionally.
  const columnWeight =
    actionsColumn && columns.every((column) => /^\d+(\.\d+)?%$/.test(column.width))
      ? columns.reduce((total, column) => total + parseFloat(column.width), 0)
      : 0;

  const defaultRenderCell = (item, column) => {
    const value = item[column.key];
    if (value === undefined || value === null) return '-';
    return String(value);
  };

  const cellRenderer = renderCell || defaultRenderCell;

  const getRowVariantClass = (rowIndex) => {
    if (variant === 'striped') {
      return rowIndex % 2 === 0 ? 'bg-black/20' : '';
    }
    if (variant === 'bordered') {
      return 'border-b border-neutral-300/10';
    }
    return '';
  };

  const resolveEmptyContent = () => {
    if (typeof emptyContent === 'string') {
      return <EmptyState title={emptyContent} />;
    }
    if (emptyContent) {
      return emptyContent;
    }
    return <EmptyState title={t('common.noData')} />;
  };

  return (
    <div className={`min-w-0 ${className}`}>
      {contextMenu.menu}
      <div className="overflow-x-auto">
        <table className="w-full table-fixed" style={{ minWidth }}>
          <thead>
            <tr className="bg-black/40">
              {columns.map((column, index) => (
                <th
                  scope="col"
                  key={index}
                  className="p-4 text-left text-neutral-400 font-mono [overflow-wrap:anywhere]"
                  style={{
                    width: columnWeight
                      ? `${(parseFloat(column.width) / columnWeight) * 100}%`
                      : column.width || 'auto',
                  }}
                >
                  {typeof column.label === 'string' ? <TruncatedText>{column.label}</TruncatedText> : column.label}
                </th>
              ))}
              {actionsColumn && (
                <th
                  scope="col"
                  className="p-4 text-left text-neutral-400 font-mono [overflow-wrap:anywhere]"
                  style={{ width: actionsWidth }}
                >
                  {typeof actionsColumn.label === 'string' ? (
                    <TruncatedText>{actionsColumn.label}</TruncatedText>
                  ) : (
                    actionsColumn.label
                  )}
                </th>
              )}
            </tr>
          </thead>

          <tbody>
            {loading ? (
              <tr>
                <td colSpan={columnCount}>
                  <Loading compact />
                </td>
              </tr>
            ) : (empty ?? data.length === 0) ? (
              <tr>
                <td colSpan={columnCount} className="p-8 text-center text-neutral-400">
                  {resolveEmptyContent()}
                </td>
              </tr>
            ) : (
              data.map((item, rowIndex) => {
                const RowComponent = animate ? motion.tr : 'tr';
                const rowMotionProps = animate
                  ? {
                      initial: { opacity: 0, y: 10 },
                      animate: { opacity: 1, y: 0 },
                      transition: { delay: rowIndex * 0.05 },
                      whileHover: { backgroundColor: 'rgba(0, 0, 0, 0.5)' },
                    }
                  : {};

                return (
                  <RowComponent
                    key={rowIndex}
                    className={`border-t border-neutral-300/10 hover:bg-black/40 transition-colors focus-visible:outline focus-visible:outline-geek-400 focus-visible:-outline-offset-2
                              ${onRowClick ? 'cursor-pointer' : ''}
                              ${getRowVariantClass(rowIndex)}
                              ${rowClassName ? rowClassName(item, rowIndex) : ''}`}
                    {...contextMenu.getRowProps(item)}
                    tabIndex={onRowClick || getRowActions ? 0 : undefined}
                    onClick={(event) => {
                      if (event.defaultPrevented || event.target.closest(INTERACTIVE_SELECTOR)) return;
                      onRowClick?.(item, rowIndex);
                    }}
                    onKeyDown={(event) => {
                      contextMenu.getRowProps(item).onKeyDown?.(event);
                      if (event.defaultPrevented || event.target !== event.currentTarget || !onRowClick) return;
                      if (event.key === 'Enter' || event.key === ' ') {
                        event.preventDefault();
                        onRowClick(item, rowIndex);
                      }
                    }}
                    {...rowMotionProps}
                  >
                    {columns.map((column, colIndex) => {
                      const content = cellRenderer(item, column, rowIndex, colIndex);
                      return (
                        <td key={colIndex} className="p-4 text-neutral-300 font-mono break-words overflow-hidden">
                          {typeof content === 'string' || typeof content === 'number' ? (
                            <TruncatedText>{content}</TruncatedText>
                          ) : (
                            content
                          )}
                        </td>
                      );
                    })}
                    {actionsColumn && (
                      <td className="p-4" data-row-interactive>
                        <RowActions actions={getRowActions?.(item) || []} />
                      </td>
                    )}
                  </RowComponent>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {footer && <div className="p-4 border-t border-neutral-300/30 bg-black/20">{footer}</div>}
    </div>
  );
}

export default List;
