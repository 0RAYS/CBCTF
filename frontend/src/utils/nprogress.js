import NProgress from 'nprogress';
import 'nprogress/nprogress.css';
import './nprogress.css'; // 导入自定义样式

// NProgress 基础配置
NProgress.configure({
  showSpinner: false,
  minimum: 0.1,
  trickleSpeed: 200,
  easing: 'ease',
  speed: 500,
});

// 计数器
let requestCount = 0;

// 开始加载
export const startLoading = () => {
  requestCount++;
  if (requestCount === 1) {
    NProgress.start();
  }
};

// 结束加载
export const finishLoading = () => {
  requestCount--;
  if (requestCount === 0) {
    NProgress.done();
  }
};

// 强制结束加载
// export const forceFinishLoading = () => {
//   requestCount = 0;
//   NProgress.done(true);
// };

// 自定义样式
const style = document.createElement('style');
style.textContent = `
    #nprogress {
        pointer-events: none;
        position: fixed;
        top: 80px;  /* NavBar 的高度 */
        left: 0;
        right: 0;
        z-index: 950;
    }

    #nprogress .bar {
        position: absolute;
        left: 0;
        width: 100%;
    }

    /* 隐藏其他效果 */
    #nprogress .peg {
        display: none;
    }
`;
document.head.appendChild(style);

// export default NProgress;
