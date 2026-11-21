import { forwardRef, useEffect, useImperativeHandle, useMemo, useRef } from 'react';
import DOMPurify from 'dompurify';
import { ansiToHtml } from '@/utils/ansi';
import Spinner from '@/components/common/Spinner';

const EMPTY_ATTRIBUTES = [];

/** AnsiLog — 将含 ANSI 转义码的日志文本渲染为带颜色的 HTML */
const AnsiLog = forwardRef(function AnsiLog(
  {
    content,
    className = '',
    loading = false,
    empty = '',
    postProcess,
    allowedAttr = EMPTY_ATTRIBUTES,
    onClick,
    sentinel,
    scrollToBottom = false,
  },
  ref
) {
  const innerRef = useRef(null);

  useImperativeHandle(ref, () => innerRef.current);

  const isEmpty = !content || (Array.isArray(content) && content.length === 0);

  const html = useMemo(() => {
    if (isEmpty) return '';

    const lines = Array.isArray(content) ? content : content.split('\n');

    const raw = lines
      .map((line) => {
        const converted = ansiToHtml(line).replace(/\n$/, '');
        const processed = postProcess ? postProcess(converted) : converted;
        return `<div class="whitespace-pre-wrap break-words leading-5 font-mono text-xs">${processed}</div>`;
      })
      .join('');

    return DOMPurify.sanitize(raw, {
      ALLOWED_TAGS: ['div', 'span'],
      ALLOWED_ATTR: ['class', 'style', ...allowedAttr],
    });
  }, [content, postProcess, allowedAttr]);

  useEffect(() => {
    if (!scrollToBottom || loading || !innerRef.current) return;
    const el = innerRef.current;
    el.scrollTop = el.scrollHeight;
  }, [html, loading, scrollToBottom]);

  const baseClass = 'bg-neutral-900 border border-neutral-700 rounded-md p-3 overflow-auto';

  return (
    <div ref={innerRef} className={`${baseClass} ${className}`} onClick={onClick}>
      {loading ? (
        <div className="flex items-center justify-center min-h-[400px]">
          <Spinner size="xs" colorClassName="border-neutral-500" />
        </div>
      ) : isEmpty ? (
        <div className="flex items-center justify-center min-h-[400px]">
          <p className="text-neutral-600 text-xs font-mono">{empty}</p>
        </div>
      ) : (
        <div dangerouslySetInnerHTML={{ __html: html }} />
      )}
      {/* sentinel 始终渲染，保证 IntersectionObserver 能在 mount 时就拿到 ref */}
      {sentinel}
    </div>
  );
});

export default AnsiLog;
