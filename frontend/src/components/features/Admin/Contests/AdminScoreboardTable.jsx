/**
 * 管理员表格式积分榜组件
 * @param {Object} props
 * @param {Array} props.teams - 队伍列表
 * @param {Array} props.challenges - 题目列表（按分类分组）
 * @param {number} props.totalCount - 总队伍数
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {Function} props.onPageChange - 页码变化处理函数
 */

import ScoreboardTable from '../../CTFGame/Scoreboard/ScoreboardTable';
import { Pagination } from '../../../../components/common';

function AdminScoreboardTable({
  teams = [],
  challenges = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  onPageChange,
}) {
  return (
    <ScoreboardTable
      teams={teams}
      challenges={challenges}
      totalCount={totalCount}
      currentPage={currentPage}
      pageSize={pageSize}
      onPageChange={onPageChange}
      isAdmin={true}
      PaginationComponent={Pagination}
    />
  );
}

export default AdminScoreboardTable;
