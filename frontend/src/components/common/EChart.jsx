import ReactEChartsCore from 'echarts-for-react/lib/core';
import * as echarts from 'echarts/core';
import { LineChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, GraphicComponent } from 'echarts/components';
import { LegacyGridContainLabel } from 'echarts/features';
import { CanvasRenderer } from 'echarts/renderers';

echarts.use([LineChart, GridComponent, TooltipComponent, GraphicComponent, LegacyGridContainLabel, CanvasRenderer]);

export default function EChart({ ref, ...props }) {
  return <ReactEChartsCore ref={ref} {...props} echarts={echarts} />;
}
