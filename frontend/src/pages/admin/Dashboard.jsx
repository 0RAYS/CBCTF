import { lazy, Suspense, useEffect, useEffectEvent, useState } from 'react';
import { getSystemStatus } from '@/api/admin/system';
import AdminDashboard from '@/components/features/Admin/AdminDashboard';
import Loading from '@/components/common/Loading';
import { toast } from '@/utils/toast.js';
import { useTranslation } from 'react-i18next';
import { startVisiblePolling } from '@/utils/visiblePolling.js';

const ReactECharts = lazy(() => import('@/components/common/EChart'));

// ECharts 需要字面量色值，取自 tailwind.config.js 的设计 token
const AXIS_LINE = '#2a2a2a';
const AXIS_LABEL = '#8a8a8a';
const SPLIT_LINE = '#1a1a1a';
const SERIES_COLORS = {
  cpu: '#597ef7',
  memory: '#4ade80',
  disk: '#facc15',
};

function Dashboard() {
  const [status, setStatus] = useState(null);
  // loading 仅覆盖首次请求，轮询只更新 status
  const [loading, setLoading] = useState(true);
  const { t } = useTranslation();

  const reportError = useEffectEvent((error) => {
    toast.danger({ description: error.message || t('admin.dashboard.toast.fetchFailed') });
  });

  useEffect(() => {
    let active = true;
    let pending = false;
    const fetchSystemStatus = async () => {
      if (pending) return;
      pending = true;
      try {
        const response = await getSystemStatus(true);
        if (active && response.code === 200) setStatus(response.data);
      } catch (error) {
        if (active) reportError(error);
      } finally {
        pending = false;
        if (active) setLoading(false);
      }
    };
    fetchSystemStatus();
    const stopPolling = startVisiblePolling(fetchSystemStatus, 10000);
    return () => {
      active = false;
      stopPolling();
    };
  }, []);

  const getChartOption = () => {
    if (!status?.metrics) return {};

    return {
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
        },
      },
      xAxis: {
        type: 'category',
        data: status.metrics.map((item) => new Date(item.timestamp).toLocaleTimeString()),
        axisLine: {
          lineStyle: {
            color: AXIS_LINE,
          },
        },
        axisLabel: {
          color: AXIS_LABEL,
        },
        splitLine: {
          show: true,
          lineStyle: {
            color: SPLIT_LINE,
            type: 'dashed',
            width: 1,
          },
        },
      },
      yAxis: {
        type: 'value',
        axisLine: {
          lineStyle: {
            color: AXIS_LINE,
          },
        },
        axisLabel: {
          color: AXIS_LABEL,
        },
        splitLine: {
          show: true,
          lineStyle: {
            color: SPLIT_LINE,
            type: 'dashed',
            width: 1,
          },
        },
      },
      series: [
        {
          name: t('admin.dashboard.chart.cpu'),
          type: 'line',
          data: status.metrics.map((item) => item.cpu),
          smooth: true,
          lineStyle: {
            color: SERIES_COLORS.cpu,
          },
          itemStyle: {
            color: SERIES_COLORS.cpu,
          },
          symbol: 'none',
        },
        {
          name: t('admin.dashboard.chart.memory'),
          type: 'line',
          data: status.metrics.map((item) => item.mem),
          smooth: true,
          lineStyle: {
            color: SERIES_COLORS.memory,
          },
          itemStyle: {
            color: SERIES_COLORS.memory,
          },
          symbol: 'none',
        },
        {
          name: t('admin.dashboard.chart.disk'),
          type: 'line',
          data: status.metrics.map((item) => item.disk),
          smooth: true,
          lineStyle: {
            color: SERIES_COLORS.disk,
          },
          itemStyle: {
            color: SERIES_COLORS.disk,
          },
          symbol: 'none',
        },
      ],
    };
  };

  if (loading) {
    return <Loading />;
  }

  return (
    <AdminDashboard
      status={status}
      chartContent={
        <div className="h-[450px] w-full">
          <Suspense
            fallback={
              <div className="flex h-full items-center justify-center text-neutral-400">{t('common.loading')}</div>
            }
          >
            <ReactECharts option={getChartOption()} style={{ height: '100%', width: '100%' }} />
          </Suspense>
        </div>
      }
    />
  );
}

export default Dashboard;
