import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

const plugins = [remarkGfm];

export default function MarkdownContent({ children, className = '' }) {
  return (
    <div className={`prose prose-invert max-w-none ${className}`}>
      <ReactMarkdown remarkPlugins={plugins}>{children || ''}</ReactMarkdown>
    </div>
  );
}
