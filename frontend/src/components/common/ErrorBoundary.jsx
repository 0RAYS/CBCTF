import { Component } from 'react';

/**
 * React 错误边界组件
 * 捕获子组件渲染时抛出的未处理错误, 显示降级 UI
 */
class ErrorBoundary extends Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }

  componentDidCatch(error, info) {
    error;
    info;
  }

  handleReset = () => {
    this.setState({ hasError: false, error: null });
  };

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.handleReset);
      }

      return (
        <div
          role="alert"
          className="flex flex-col items-center justify-center min-h-[200px] p-8 text-center text-neutral-400 font-mono"
        >
          <p className="text-lg text-red-400 mb-2">Something went wrong.</p>
          <p className="text-sm mb-4 max-w-md truncate opacity-70">{String(this.state.error)}</p>
          <button
            onClick={this.handleReset}
            className="px-4 py-2 text-sm border border-neutral-500 rounded-md hover:border-neutral-300 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70"
          >
            Try again
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
