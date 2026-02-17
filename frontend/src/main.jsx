import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import 'nprogress/nprogress.css';
import App from './App.jsx';
import './i18n';
import { store } from './store';
import { Provider } from 'react-redux';
import ToastProvider from './components/common/ToastProvider';
import { WebSocketProvider } from './components/common/WebSocketProvider';
import ModalProvider from './components/common/ModalProvider';

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <ModalProvider>
      <Provider store={store}>
        <WebSocketProvider>
          <ToastProvider position="bottom-right" maxToasts={5}>
            <App />
          </ToastProvider>
        </WebSocketProvider>
      </Provider>
    </ModalProvider>
  </StrictMode>
);
