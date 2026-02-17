import { useState, useEffect } from 'react';
import ReactECharts from 'echarts-for-react';
import { getSystemStatus } from '../../api/admin/system';
import AdminDashboard from '../../components/features/Admin/AdminDashboard';
import { toast } from '../../utils/toast.js';
import { useTranslation } from 'react-i18next';

function Dashboard() {
  const [status, setStatus] = useState(null);
  const { t } = useTranslation();

  const fetchSystemStatus = async () => {
    try {
      const response = await getSystemStatus(true);
      if (response.code === 200) {
        setStatus(response.data);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.dashboard.toast.fetchFailed') });
    }
  };

  useEffect(() => {
    fetchSystemStatus().then();
    const interval = setInterval(() => {
      fetchSystemStatus().then();
    }, 3000);
    return () => {
      clearInterval(interval);
    };
  }, []);

  // 生成 ECharts 配置
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
            color: '#ccc',
          },
        },
        axisLabel: {
          color: '#666',
        },
        splitLine: {
          show: true,
          lineStyle: {
            color: '#e0e0e0',
            type: 'dashed',
            width: 1,
          },
        },
      },
      yAxis: {
        type: 'value',
        axisLine: {
          lineStyle: {
            color: '#ccc',
          },
        },
        axisLabel: {
          color: '#666',
        },
        splitLine: {
          show: true,
          lineStyle: {
            color: '#e0e0e0',
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
            color: '#8884d8',
          },
          itemStyle: {
            color: '#8884d8',
          },
          symbol: 'none',
        },
        {
          name: t('admin.dashboard.chart.memory'),
          type: 'line',
          data: status.metrics.map((item) => item.mem),
          smooth: true,
          lineStyle: {
            color: '#82ca9d',
          },
          itemStyle: {
            color: '#82ca9d',
          },
          symbol: 'none',
        },
        {
          name: t('admin.dashboard.chart.disk'),
          type: 'line',
          data: status.metrics.map((item) => item.disk),
          smooth: true,
          lineStyle: {
            color: '#ffc658',
          },
          itemStyle: {
            color: '#ffc658',
          },
          symbol: 'none',
        },
      ],
    };
  };

  return (
    <div className="p-4">
      <AdminDashboard
        status={status}
        chartContent={
          <div style={{ width: '100%', height: '450px' }}>
            <ReactECharts option={getChartOption()} style={{ height: '100%', width: '100%' }} />
          </div>
        }
      />
    </div>
  );
}

export default Dashboard;
