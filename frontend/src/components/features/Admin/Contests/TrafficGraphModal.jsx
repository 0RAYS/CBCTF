import { useEffect, useRef, useState } from 'react';
import ReactECharts from 'echarts-for-react';
import { motion } from 'motion/react';
import { IconX, IconRefresh, IconPlayerPlay, IconPlayerPause } from '@tabler/icons-react';
import { Button } from '../../../../components/common';
import { toast } from '../../../../utils/toast';
import { getContestTeamTraffic } from '../../../../api/admin/contest.js';
import { useTranslation } from 'react-i18next';

/**
 * 流量关系图弹窗组件
 * @param {Object} props
 * @param {boolean} props.isOpen - 弹窗是否打开
 * @param {Function} props.onClose - 关闭弹窗回调
 * @param {Object} props.container - 容器信息
 * @param {number} props.contestId - 比赛ID
 * @param {number} props.teamId - 队伍ID
 */
function TrafficGraphModal({ isOpen, onClose, container, contestId, teamId }) {
  const { t, i18n } = useTranslation();
  const chartRef = useRef(null);
  const [connections, setConnections] = useState([]);
  const [ipL, setIpL] = useState([]);
  const [timeShift, setTimeShift] = useState(0);
  const [requestDuration, setRequestDuration] = useState(1);
  const [maxDuration, setMaxDuration] = useState(60);

  // 播放相关状态
  const [isPlaying, setIsPlaying] = useState(false);
  const [playInterval, setPlayInterval] = useState(null);
  const [useDemoMode, setUseDemoMode] = useState(false);

  // 生成样例数据 - 模拟20秒时长的网络流量
  const generateDemoData = (shift) => {
    const demoIPs = [
      '10.0.1.1',
      '10.0.1.2',
      '10.0.1.3',
      '192.168.1.100',
      '192.168.1.101',
      '172.16.0.5',
      '172.16.0.10',
      '8.8.8.8',
    ];

    // 基础连接定义，根据时间偏移产生变化
    const baseConnections = [
      { src: 0, dst: 3, type: 'TCP', subtype: 'HTTP' },
      { src: 0, dst: 4, type: 'TCP', subtype: 'HTTPS' },
      { src: 1, dst: 5, type: 'UDP', subtype: 'DNS' },
      { src: 1, dst: 7, type: 'UDP', subtype: 'DNS' },
      { src: 2, dst: 3, type: 'TCP', subtype: 'SSH' },
      { src: 2, dst: 6, type: 'TCP', subtype: 'HTTPS' },
      { src: 3, dst: 7, type: 'TCP', subtype: 'HTTP' },
      { src: 4, dst: 5, type: 'UDP', subtype: 'QUIC' },
      { src: 5, dst: 6, type: 'TCP', subtype: 'HTTPS' },
      { src: 6, dst: 7, type: 'TCP', subtype: 'HTTP' },
      { src: 0, dst: 7, type: 'TCP', subtype: 'HTTPS' },
      { src: 1, dst: 3, type: 'UDP', subtype: 'DNS' },
    ];

    // 根据 shift 决定哪些连接可见、流量大小变化
    const seed = shift % 20;
    const demoConnections = baseConnections
      .filter((_, i) => {
        // 随时间逐步显示更多连接，shift=0 至少6条，shift=19 全部12条
        const threshold = 6 + Math.floor((seed / 19) * 6);
        return i < threshold;
      })
      .map((c, i) => ({
        src_ip: demoIPs[c.src],
        dst_ip: demoIPs[c.dst],
        type: c.type,
        subtype: c.subtype,
        // 流量大小和次数随时间波动
        size: 200 + Math.floor(Math.sin(seed * 0.8 + i) * 150 + 150) * (i + 1),
        count: 5 + Math.floor(Math.abs(Math.cos(seed * 0.5 + i * 1.3)) * 40),
      }));

    return {
      ip: demoIPs,
      connections: demoConnections,
      duration: 20,
    };
  };

  // 从后端获取流量数据
  const fetchTrafficData = async (shift = 0, dur = 1, forceLive = false) => {
    if (!container?.id) return;

    // 已处于demo模式且非强制刷新时，直接使用样例数据
    if (useDemoMode && !forceLive) {
      const demo = generateDemoData(shift);
      setConnections(demo.connections);
      setIpL(demo.ip);
      setMaxDuration(demo.duration);
      return;
    }

    try {
      const response = await getContestTeamTraffic(contestId, teamId, container.id, {
        time_shift: shift,
        duration: dur,
      });

      if (response.code === 200) {
        setUseDemoMode(false);
        setConnections(response.data.connections || []);
        setIpL(response.data.ip || []);
        setMaxDuration(response.data.duration || 60);
      } else {
        throw new Error(t('admin.contests.trafficGraph.toast.fetchFailed'));
      }
    } catch {
      // 请求失败时使用样例数据
      toast.warning({ description: t('admin.contests.trafficGraph.toast.demoFallback') });
      setUseDemoMode(true);
      const demo = generateDemoData(shift);
      setConnections(demo.connections);
      setIpL(demo.ip);
      setMaxDuration(demo.duration);
    }
  };

  // 开始播放
  const startPlayback = () => {
    if (isPlaying) return;

    setIsPlaying(true);
    let currentShift = timeShift;

    const interval = setInterval(() => {
      if (currentShift >= maxDuration) {
        // 播放完成，重置到开始
        currentShift = 0;
        setTimeShift(0);
      } else {
        currentShift += 1;
        setTimeShift(currentShift);
      }
    }, 1000); // 每秒更新一次

    setPlayInterval(interval);
  };

  // 停止播放
  const stopPlayback = () => {
    if (!isPlaying) return;

    setIsPlaying(false);
    if (playInterval) {
      clearInterval(playInterval);
      setPlayInterval(null);
    }
  };

  // 生成图表数据
  const generateChartData = () => {
    if (!connections || !ipL) return { nodes: [], links: [], effectScatterData: [], linesData: [] };

    const nodes = [];
    const links = [];
    const effectScatterData = [];
    const linesData = [];
    const nodeMap = new Map(); // ip -> { id, x, y }

    // 预计算每个IP的连接数
    const connectionCount = new Map();
    connections.forEach((conn) => {
      connectionCount.set(conn.src_ip, (connectionCount.get(conn.src_ip) || 0) + 1);
      connectionCount.set(conn.dst_ip, (connectionCount.get(conn.dst_ip) || 0) + 1);
    });
    const maxConn = Math.max(1, ...connectionCount.values());

    // 计算节点位置 - 将IP节点排列在一个圆形上
    const radius = 200;

    // 处理所有IP节点 - 在圆形上均匀分布
    ipL.forEach((ip, index) => {
      const nodeId = `ip_${index}`;
      const angle = (index / ipL.length) * 2 * Math.PI;
      const x = radius * Math.cos(angle);
      const y = radius * Math.sin(angle);
      const count = connectionCount.get(ip) || 0;
      const sizeRatio = count / maxConn;
      const symbolSize = 20 + sizeRatio * 40; // 20–60

      nodes.push({
        id: nodeId,
        name: ip,
        symbolSize,
        x,
        y,
        connectionCount: count,
        itemStyle: {
          color: {
            type: 'radial',
            x: 0.5,
            y: 0.5,
            r: 0.5,
            colorStops: [
              { offset: 0, color: '#93c5fd' },
              { offset: 0.7, color: '#3b82f6' },
              { offset: 1, color: '#1e3a5f' },
            ],
          },
          shadowBlur: 15,
          shadowColor: 'rgba(59, 130, 246, 0.6)',
        },
      });
      nodeMap.set(ip, { id: nodeId, x, y });

      // effectScatter 涟漪数据
      effectScatterData.push({
        value: [x, y],
        symbolSize: symbolSize * 0.6,
      });
    });

    // 处理连接关系
    connections.forEach((connection) => {
      const sourceInfo = nodeMap.get(connection.src_ip);
      const targetInfo = nodeMap.get(connection.dst_ip);

      if (sourceInfo && targetInfo) {
        const lineWidth = Math.min(connection.count / 10, 5) + 1;
        const color = getConnectionColor(connection.type);

        links.push({
          source: sourceInfo.id,
          target: targetInfo.id,
          value: connection.size,
          count: connection.count,
          type: connection.type,
          subtype: connection.subtype,
          lineStyle: {
            color,
            width: lineWidth,
          },
        });

        // lines 系列数据 - 粒子流动
        const speed = 0.3 + (connection.count / maxConn) * 0.7; // 0.3–1
        const trailLength = 0.1 + (connection.count / maxConn) * 0.2; // 0.1–0.3
        const particleSize = 1.5 + (connection.count / maxConn) * 1.5; // 1.5–3

        linesData.push({
          coords: [
            [sourceInfo.x, sourceInfo.y],
            [targetInfo.x, targetInfo.y],
          ],
          lineStyle: {
            color,
            width: 0,
          },
          effect: {
            period: 12 / speed,
            trailLength,
            symbolSize: particleSize,
          },
        });
      }
    });

    return { nodes, links, effectScatterData, linesData };
  };

  // 根据连接类型获取颜色
  const getConnectionColor = (type) => {
    switch (type?.toUpperCase()) {
      case 'TCP':
        return '#3b82f6';
      case 'UDP':
        return '#ef4444';
      default:
        return '#6b7280';
    }
  };

  const getChartOption = () => {
    const { nodes, links, effectScatterData, linesData } = generateChartData();

    return {
      tooltip: {
        trigger: 'item',
        formatter: function (params) {
          if (params.dataType === 'node') {
            const connCount = params.data.connectionCount ?? 0;
            return `<div style="color: #f3f4f6;">
              <strong>${params.data.name}</strong><br/>
              ${t('admin.contests.trafficGraph.tooltip.connections', { count: connCount })}
            </div>`;
          } else if (params.dataType === 'edge') {
            const connectionInfo = params.data;
            const countValue = connectionInfo.count ?? t('common.notAvailable');
            const typeValue = connectionInfo.type || t('common.notAvailable');
            const subtypeValue = connectionInfo.subtype || t('common.notAvailable');
            return `<div style="color: #f3f4f6;">
              <strong>${t('admin.contests.trafficGraph.tooltip.connectionTitle')}</strong><br/>
              ${t('admin.contests.trafficGraph.tooltip.size', { value: connectionInfo.value })}<br/>
              ${t('admin.contests.trafficGraph.tooltip.count', { count: countValue })}<br/>
              ${t('admin.contests.trafficGraph.tooltip.protocol', { type: typeValue, subtype: subtypeValue })}
            </div>`;
          }
          return '';
        },
        backgroundColor: 'rgba(0, 0, 0, 0.9)',
        borderColor: '#374151',
        borderWidth: 1,
        textStyle: {
          color: '#f3f4f6',
        },
      },
      xAxis: {
        show: false,
        type: 'value',
        min: -280,
        max: 280,
      },
      yAxis: {
        show: false,
        type: 'value',
        min: -280,
        max: 280,
      },
      dataZoom: [
        {
          type: 'inside',
          xAxisIndex: 0,
          filterMode: 'none',
        },
        {
          type: 'inside',
          yAxisIndex: 0,
          filterMode: 'none',
        },
      ],
      graphic: [
        {
          type: 'group',
          right: 20,
          bottom: 20,
          children: [
            {
              type: 'rect',
              shape: { width: 120, height: 80, r: 6 },
              style: { fill: 'rgba(0,0,0,0.6)', stroke: '#374151', lineWidth: 1 },
            },
            { type: 'circle', shape: { cx: 18, cy: 18, r: 5 }, style: { fill: '#3b82f6' } },
            { type: 'text', style: { text: 'TCP', fill: '#d1d5db', fontSize: 11, x: 30, y: 13 } },
            { type: 'circle', shape: { cx: 18, cy: 42, r: 5 }, style: { fill: '#ef4444' } },
            { type: 'text', style: { text: 'UDP', fill: '#d1d5db', fontSize: 11, x: 30, y: 37 } },
            { type: 'circle', shape: { cx: 18, cy: 66, r: 5 }, style: { fill: '#6b7280' } },
            { type: 'text', style: { text: 'Other', fill: '#d1d5db', fontSize: 11, x: 30, y: 61 } },
          ],
        },
      ],
      animation: true,
      animationDuration: 1000,
      animationDurationUpdate: 1500,
      animationEasingUpdate: 'quinticInOut',
      series: [
        // Series 0: effectScatter 涟漪脉冲
        {
          type: 'effectScatter',
          coordinateSystem: 'cartesian2d',
          data: effectScatterData,
          symbolSize: (val, params) => params.data.symbolSize || 12,
          showEffectOn: 'render',
          rippleEffect: {
            brushType: 'stroke',
            period: 4,
            scale: 3,
            number: 2,
          },
          itemStyle: {
            color: 'rgba(59, 130, 246, 0.4)',
          },
          silent: true,
          z: 1,
        },
        // Series 1: graph 网络拓扑
        {
          type: 'graph',
          layout: 'none',
          coordinateSystem: 'cartesian2d',
          data: nodes,
          links: links,
          edgeSymbol: ['circle', 'arrow'],
          edgeSymbolSize: [4, 10],
          label: {
            show: true,
            position: 'inside',
            formatter: '{b}',
            color: '#f3f4f6',
            fontSize: 9,
          },
          emphasis: {
            focus: 'adjacency',
            itemStyle: {
              shadowBlur: 20,
              shadowColor: 'rgba(147, 197, 253, 0.8)',
            },
            lineStyle: {
              width: 4,
            },
          },
          labelLayout: {
            hideOverlap: true,
          },
          lineStyle: {
            curveness: 0.2,
            opacity: 0.7,
          },
          z: 2,
        },
        // Series 2: lines 粒子流动
        {
          type: 'lines',
          coordinateSystem: 'cartesian2d',
          polyline: false,
          data: linesData,
          effect: {
            show: true,
            period: 16,
            trailLength: 0.15,
            symbol: 'arrow',
            symbolSize: 2,
            color: '#93c5fd',
          },
          lineStyle: {
            width: 0,
            opacity: 0,
          },
          silent: true,
          z: 3,
        },
      ],
    };
  };

  useEffect(() => {
    if (isOpen && container?.id) {
      fetchTrafficData(timeShift, requestDuration);
    }
  }, [isOpen, container?.id, timeShift, requestDuration]);

  useEffect(() => {
    if (isOpen && chartRef.current) {
      const chart = chartRef.current.getEchartsInstance();
      chart.resize();
    }
  }, [isOpen]);

  // 清理播放定时器
  useEffect(() => {
    return () => {
      if (playInterval) {
        clearInterval(playInterval);
      }
    };
  }, [playInterval]);

  const handleRefresh = () => {
    fetchTrafficData(timeShift, requestDuration, true);
  };

  const handleTimeShiftChange = (newShift) => {
    setTimeShift(newShift);
  };

  const handleRequestDurationChange = (newDuration) => {
    setRequestDuration(newDuration);
  };

  const handlePlayPause = () => {
    if (isPlaying) {
      stopPlayback();
    } else {
      startPlayback();
    }
  };

  if (!isOpen) return null;

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        className="relative w-full max-w-7xl h-[90vh] bg-neutral-900 border border-neutral-700 rounded-lg shadow-2xl overflow-hidden flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        {/* 弹窗头部 */}
        <div className="flex items-center justify-between p-4 border-b border-neutral-700 bg-neutral-800 flex-shrink-0">
          <div className="flex items-center gap-4">
            <h2 className="text-xl font-mono text-neutral-50">{t('admin.contests.trafficGraph.title')}</h2>
            <div className="flex items-center gap-2">
              <span className="text-sm text-neutral-400">{t('admin.contests.trafficGraph.controls.timeShift')}</span>
              <input
                type="number"
                min="0"
                max={maxDuration}
                value={timeShift}
                onChange={(e) => handleTimeShiftChange(parseInt(e.target.value) || 0)}
                className="w-20 px-2 py-1 bg-neutral-700 border border-neutral-600 rounded text-white text-sm"
              />
              <span className="text-sm text-neutral-400">{t('admin.contests.trafficGraph.controls.seconds')}</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm text-neutral-400">{t('admin.contests.trafficGraph.controls.timeSlice')}</span>
              <input
                type="number"
                min="1"
                max={maxDuration}
                value={requestDuration}
                onChange={(e) => handleRequestDurationChange(parseInt(e.target.value) || 60)}
                className="w-20 px-2 py-1 bg-neutral-700 border border-neutral-600 rounded text-white text-sm"
              />
              <span className="text-sm text-neutral-400">{t('admin.contests.trafficGraph.controls.seconds')}</span>
            </div>
            <Button variant="ghost" size="icon" className="!text-geek-400 hover:!text-geek-300" onClick={handleRefresh}>
              <IconRefresh size={16} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className={`${isPlaying ? '!text-red-400 hover:!text-red-300' : '!text-green-400 hover:!text-green-300'}`}
              onClick={handlePlayPause}
            >
              {isPlaying ? <IconPlayerPause size={16} /> : <IconPlayerPlay size={16} />}
            </Button>
          </div>
          <Button variant="ghost" size="icon" className="!text-neutral-400 hover:!text-neutral-200" onClick={onClose}>
            <IconX size={20} />
          </Button>
        </div>

        {/* 图表容器 */}
        <div className="flex-1 p-4 min-h-0 flex flex-col">
          {useDemoMode && (
            <div className="mb-2 px-4 py-2 bg-yellow-900/60 border border-yellow-600/50 rounded text-yellow-300 text-sm text-center flex-shrink-0">
              {t('admin.contests.trafficGraph.demoNotice')}
            </div>
          )}
          <div className="flex-1 min-h-0">
            <ReactECharts
              ref={chartRef}
              option={getChartOption()}
              notMerge={false}
              lazyUpdate={true}
              style={{ height: '100%', width: '100%' }}
              opts={{
                renderer: 'canvas',
              }}
            />
          </div>
        </div>

        {/* 底部信息 */}
        <div className="p-4 border-t border-neutral-700 bg-neutral-800 flex-shrink-0">
          <div className="flex items-center justify-between text-sm text-neutral-400">
            <div className="flex items-center gap-4">
              {
                <>
                  <span className="font-mono">
                    {t('admin.contests.trafficGraph.footer.ipCount', { count: ipL.length || 0 })}
                  </span>
                  <span className="font-mono">
                    {t('admin.contests.trafficGraph.footer.connectionCount', { count: connections.length || 0 })}
                  </span>
                  <span className="font-mono">
                    {t('admin.contests.trafficGraph.footer.maxDuration', { count: maxDuration || 0 })}
                  </span>
                  <span className="font-mono">
                    {t('admin.contests.trafficGraph.footer.timeSlice', { count: requestDuration || 0 })}
                  </span>
                  {isPlaying && (
                    <span className="font-mono text-green-400">
                      {t('admin.contests.trafficGraph.footer.playing', {
                        current: timeShift,
                        total: maxDuration,
                      })}
                    </span>
                  )}
                  {useDemoMode && (
                    <span className="font-mono text-yellow-400">
                      {t('admin.contests.trafficGraph.footer.demoMode')}
                    </span>
                  )}
                </>
              }
              <span className="font-mono">
                {t('admin.contests.trafficGraph.footer.updatedAt', {
                  time: new Date().toLocaleString(i18n.language || 'en-US'),
                })}
              </span>
            </div>
          </div>
        </div>
      </motion.div>
    </motion.div>
  );
}

export default TrafficGraphModal;
