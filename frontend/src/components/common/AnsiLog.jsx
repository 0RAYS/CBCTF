import { useMemo, forwardRef, useRef, useEffect, useImperativeHandle } from 'react';
import DOMPurify from 'dompurify';
import { ansiToHtml } from '../../utils/ansi';

/**
 * AnsiLog — 将含 ANSI 转义码的日志文本渲染为带颜色的 HTML。
 *
 * Props:
 *   content        {string|string[]} 原始日志文本（可含 ANSI 转义码）；传数组时每项视为一行
 *   className      {string}          附加到外层容器的 class
 *   loading        {boolean}         显示加载中状态
 *   empty          {string}          无内容时的提示文案
 *   postProcess    {function}        (lineHtml: string) => string，对每行转换后的 HTML 做额外处理
 *   allowedAttr    {string[]}        额外允许的 HTML 属性（追加到默认白名单 class/style 之上）
 *   onClick        {function}        点击事件（用于事件委托）
 *   sentinel       {ReactNode}       渲染在内容末尾的哨兵节点（用于无限滚动）
 *   scrollToBottom {boolean}         内容更新后是否自动滚动到底部（默认 false）
 */
const AnsiLog = forwardRef(function AnsiLog(
  {
    content,
    className = '',
    loading = false,
    empty = '',
    postProcess,
    allowedAttr = [],
    onClick,
    sentinel,
    scrollToBottom = false,
  },
  ref
) {
  const innerRef = useRef(null);

  // 将外部 ref 指向内部滚动容器
  useImperativeHandle(ref, () => innerRef.current);

  const html = useMemo(() => {
    if (!content || (Array.isArray(content) && content.length === 0)) return '';

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

  // 内容更新后自动滚到底部
  useEffect(() => {
    if (!scrollToBottom || loading || !innerRef.current) return;
    const el = innerRef.current;
    el.scrollTop = el.scrollHeight;
  }, [html, loading, scrollToBottom]);

  const baseClass = 'bg-neutral-950 border border-neutral-700 rounded p-3';

  if (loading) {
    return (
      <div className={`flex items-center justify-center ${baseClass} min-h-[400px] ${className}`}>
        <div className="w-3 h-3 rounded-full border border-neutral-500 border-t-transparent animate-spin" />
      </div>
    );
  }

  if (!content || (Array.isArray(content) && content.length === 0)) {
    return (
      <div className={`flex items-center justify-center ${baseClass} min-h-[400px] ${className}`}>
        <p className="text-neutral-600 text-xs font-mono">{empty}</p>
      </div>
    );
  }

  return (
    <div ref={innerRef} className={`overflow-auto ${baseClass} ${className}`} onClick={onClick}>
      <div dangerouslySetInnerHTML={{ __html: html }} />
      {sentinel}
    </div>
  );
});

export default AnsiLog;
