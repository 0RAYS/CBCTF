import './lib/monacoSetup';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import 'nprogress/nprogress.css';
import App from './App.jsx';
import './i18n';
import { store } from './store';
import { Provider } from 'react-redux';
import ToastProvider from './components/common/ToastProvider';
import ModalProvider from './components/common/ModalProvider';
import ErrorBoundary from './components/common/ErrorBoundary';

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <ErrorBoundary>
      <ModalProvider>
        <Provider store={store}>
            <ToastProvider position="bottom-right" maxToasts={5}>
              <App />
            </ToastProvider>
          </Provider>
      </ModalProvider>
    </ErrorBoundary>
  </StrictMode>
);
