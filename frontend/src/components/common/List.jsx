import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import Loading from './Loading';
import EmptyState from './EmptyState';

/**
 * 通用列表/表格组件（扩展自AdminList）
 * @param {Object} props
 * @param {Array} props.columns - 列定义, 格式: [{key: 'id', label: 'ID', width: '10%'}, ...]
 * @param {Array} props.data - 列表数据
 * @param {Function} props.renderCell - 自定义单元格渲染函数 (item, column, rowIndex, colIndex) => ReactNode
 * @param {Function} props.onRowClick - 行点击事件
 * @param {boolean} props.loading - 加载状态
 * @param {boolean} props.empty - 是否为空
 * @param {ReactNode|string} props.emptyContent - 空状态内容
 * @param {ReactNode} props.footer - 底部内容, 通常是分页
 * @param {Function} props.rowClassName - 自定义行类名函数 (item, index) => string
 * @param {'default'|'striped'|'bordered'} props.variant - 表格变体
 * @param {boolean} props.animate - 是否启用行动画
 * @param {string} props.className - 额外的自定义类名
 * @param {number|string} props.minWidth - Table minimum width; narrow screens scroll inside the table region
 */
function List({
  columns = [],
  data = [],
  renderCell,
  onRowClick,
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

  // 默认的单元格渲染逻辑
  const defaultRenderCell = (item, column) => {
    const value = item[column.key];
    if (value === undefined || value === null) return '-';
    return String(value);
  };

  // 使用传入的自定义渲染函数或默认渲染函数
  const cellRenderer = renderCell || defaultRenderCell;

  // 表格变体样式
  const getRowVariantClass = (rowIndex) => {
    if (variant === 'striped') {
      return rowIndex % 2 === 0 ? 'bg-black/20' : '';
    }
    if (variant === 'bordered') {
      return 'border-b border-neutral-300/10';
    }
    return '';
  };

  // 默认空状态内容
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
    <div className={className}>
      {/* 表格区域 */}
      <div className="overflow-x-auto">
        <table className="w-full table-fixed" style={{ minWidth }}>
          {/* 表头 */}
          <thead>
            <tr className="bg-black/40">
              {columns.map((column, index) => (
                <th
                  scope="col"
                  key={index}
                  className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap"
                  style={{ width: column.width || 'auto' }}
                >
                  {column.label}
                </th>
              ))}
            </tr>
          </thead>

          {/* 表格内容 */}
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={columns.length}>
                  <Loading compact />
                </td>
              </tr>
            ) : (empty ?? data.length === 0) ? (
              <tr>
                <td colSpan={columns.length} className="p-8 text-center text-neutral-400">
                  {resolveEmptyContent()}
                </td>
              </tr>
            ) : (
              // 数据行
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
                    className={`border-t border-neutral-300/10 hover:bg-black/40 transition-colors
                              ${onRowClick ? 'cursor-pointer' : ''}
                              ${getRowVariantClass(rowIndex)}
                              ${rowClassName ? rowClassName(item, rowIndex) : ''}`}
                    onClick={() => onRowClick && onRowClick(item, rowIndex)}
                    {...rowMotionProps}
                  >
                    {columns.map((column, colIndex) => (
                      <td key={colIndex} className="p-4 text-neutral-300 font-mono break-words overflow-hidden">
                        {cellRenderer(item, column, rowIndex, colIndex)}
                      </td>
                    ))}
                  </RowComponent>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {/* 底部区域 */}
      {footer && <div className="p-4 border-t border-neutral-300/30 bg-black/20">{footer}</div>}
    </div>
  );
}

export default List;
