import { useEffect, useRef, useState } from 'react';

/** 滚动文本组件 */
function ScrollingText({ text, className = '', speed = 30, maxWidth = 200 }) {
  const [shouldScroll, setShouldScroll] = useState(false);
  const containerRef = useRef(null);
  const textRef = useRef(null);

  useEffect(() => {
    if (!containerRef.current || !textRef.current) return;

    const container = containerRef.current;
    const textElement = textRef.current;

    const needsScroll = textElement.scrollWidth > container.clientWidth;
    setShouldScroll(needsScroll);

    if (needsScroll) {
      let animationId;
      let startTime = Date.now();
      let isPausedState = false;

      const animate = () => {
        if (!isPausedState) {
          const elapsed = Date.now() - startTime;
          const translateX = -(elapsed / 1000) * speed;

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
      className={`min-w-0 overflow-hidden ${className}`}
      style={{ maxWidth: `${maxWidth}px` }}
      title={text}
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
