import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { memo } from 'react';

const plugins = [remarkGfm];
const components = {
  // Discard the AST node; it is not a DOM attribute.
  img: ({ node: _, ...props }) => <img {...props} loading="lazy" decoding="async" />,
};

function MarkdownContent({ children, className = '' }) {
  return (
    <div className={`prose prose-invert max-w-none ${className}`}>
      <ReactMarkdown remarkPlugins={plugins} components={components}>
        {children || ''}
      </ReactMarkdown>
    </div>
  );
}

export default memo(MarkdownContent);
