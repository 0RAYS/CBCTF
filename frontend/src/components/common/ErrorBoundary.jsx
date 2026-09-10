import { Component } from 'react';
import { withTranslation } from 'react-i18next';

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

  componentDidUpdate(previousProps) {
    if (previousProps.resetKey !== this.props.resetKey && this.state.hasError) {
      this.handleReset();
    }
  }

  handleReset = () => {
    this.setState({ hasError: false, error: null });
  };

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.handleReset);
      }

      const { t } = this.props;
      const isChunkError = /dynamically imported module|module script|loading chunk|CSS chunk|Importing a module/i.test(
        String(this.state.error)
      );

      return (
        <div
          role="alert"
          className="flex flex-col items-center justify-center min-h-[200px] p-8 text-center text-neutral-400 font-mono"
        >
          <p className="text-lg text-neutral-100 mb-2">{t('common.pageError.title')}</p>
          <p className="text-sm mb-5 max-w-md leading-relaxed">
            {t(isChunkError ? 'common.pageError.resource' : 'common.pageError.description')}
          </p>
          <div className="flex flex-wrap justify-center gap-3">
            {!isChunkError && (
              <button
                type="button"
                onClick={this.handleReset}
                className="px-4 py-2 text-sm border border-neutral-500 rounded-md hover:border-neutral-300 transition-colors"
              >
                {t('common.retry')}
              </button>
            )}
            <button
              type="button"
              onClick={() => window.location.reload()}
              className="px-4 py-2 text-sm border border-geek-400 text-geek-300 rounded-md hover:bg-geek-400/10 transition-colors"
            >
              {t('common.refresh')}
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default withTranslation()(ErrorBoundary);
