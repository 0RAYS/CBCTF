import { useState, useEffect, useRef } from 'react';

/**
 * 滚动文本组件
 * @param {Object} props
 * @param {string} props.text - 要显示的文本
 * @param {string} props.className - 自定义样式类
 * @param {number} props.speed - 滚动速度（像素/秒）
 * @param {number} props.maxWidth - 最大宽度
 */
function ScrollingText({ text, className = '', speed = 30, maxWidth = 200 }) {
  const [shouldScroll, setShouldScroll] = useState(false);
  const containerRef = useRef(null);
  const textRef = useRef(null);

  useEffect(() => {
    if (!containerRef.current || !textRef.current) return;

    const container = containerRef.current;
    const textElement = textRef.current;

    // 检查是否需要滚动
    const needsScroll = textElement.scrollWidth > container.clientWidth;
    setShouldScroll(needsScroll);

    if (needsScroll) {
      // 开始滚动动画
      let animationId;
      let startTime = Date.now();
      let isPausedState = false;

      const animate = () => {
        if (!isPausedState) {
          const elapsed = Date.now() - startTime;
          const translateX = -(elapsed / 1000) * speed;

          // 当文本完全移出容器时，重置位置
          if (translateX < -(textElement.scrollWidth - container.clientWidth)) {
            startTime = Date.now();
            textElement.style.transform = 'translateX(0)';
          } else {
            textElement.style.transform = `translateX(${translateX}px)`;
          }
        }

        animationId = requestAnimationFrame(animate);
      };

      animate();

      // 鼠标悬停时暂停
      const handleMouseEnter = () => {
        isPausedState = true;
      };

      const handleMouseLeave = () => {
        isPausedState = false;
        startTime = Date.now() - (Date.now() - startTime);
      };

      container.addEventListener('mouseenter', handleMouseEnter);
      container.addEventListener('mouseleave', handleMouseLeave);

      return () => {
        cancelAnimationFrame(animationId);
        container.removeEventListener('mouseenter', handleMouseEnter);
        container.removeEventListener('mouseleave', handleMouseLeave);
      };
    }
  }, [text, speed]);

  return (
    <div
      ref={containerRef}
      className={`overflow-hidden ${className}`}
      style={{ maxWidth: `${maxWidth}px` }}
      title={shouldScroll ? text : undefined}
    >
      <div
        ref={textRef}
        className={`whitespace-nowrap ${shouldScroll ? '' : 'truncate'}`}
        style={shouldScroll ? { transition: 'transform 0.1s ease-out' } : {}}
      >
        {text}
      </div>
    </div>
  );
}

export default ScrollingText;
