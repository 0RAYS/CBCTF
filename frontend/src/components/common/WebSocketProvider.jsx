import { createContext, useContext, useEffect, useState } from 'react';
import websocketService from '../../api/websocket.js';

const WebSocketContext = createContext();

// eslint-disable-next-line react-refresh/only-export-components
export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
};

export const WebSocketProvider = ({ children }) => {
  const [isConnected, setIsConnected] = useState(false);
  const [connectionInfo, setConnectionInfo] = useState({
    connected: false,
    readyState: 'CLOSED',
    reconnectAttempts: 0,
    maxReconnectAttempts: 5,
    hasToken: false,
  });

  useEffect(() => {
    // 检查初始token状态
    const checkTokenAndConnect = () => {
      const hasToken = websocketService.hasValidToken();
      setConnectionInfo((prev) => ({ ...prev, hasToken }));

      if (hasToken) {
        websocketService.connect();
      }
    };

    // 初始检查
    checkTokenAndConnect();

    // 监听连接状态变化
    const removeConnectionHandler = websocketService.addConnectionHandler((connected) => {
      setIsConnected(connected);
      setConnectionInfo(websocketService.getConnectionInfo());
    });

    // 监听localStorage变化（用于同一页面内的token变化）
    const handleStorageChange = () => {
      checkTokenAndConnect();
    };

    // 监听localStorage变化
    window.addEventListener('storage', handleStorageChange);
    const handleOnline = () => {
      if (websocketService.hasValidToken()) {
        websocketService.connect();
      }
    };
    const handleVisibility = () => {
      if (document.visibilityState === 'visible' && websocketService.hasValidToken()) {
        websocketService.connect();
      }
    };
    window.addEventListener('online', handleOnline);
    document.addEventListener('visibilitychange', handleVisibility);

    // 定期检查token状态（用于处理程序化token变化）
    const tokenCheckInterval = setInterval(() => {
      const currentHasToken = websocketService.hasValidToken();
      if (currentHasToken !== connectionInfo.hasToken) {
        checkTokenAndConnect();
      }
    }, 1000);

    // 页面卸载时清理
    return () => {
      removeConnectionHandler();
      window.removeEventListener('storage', handleStorageChange);
      window.removeEventListener('online', handleOnline);
      document.removeEventListener('visibilitychange', handleVisibility);
      clearInterval(tokenCheckInterval);
    };
  }, [connectionInfo.hasToken]);

  // 定期更新连接信息
  useEffect(() => {
    const interval = setInterval(() => {
      setConnectionInfo(websocketService.getConnectionInfo());
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const value = {
    isConnected,
    connectionInfo,
    send: websocketService.send.bind(websocketService),
    addMessageHandler: websocketService.addMessageHandler.bind(websocketService),
    addConnectionHandler: websocketService.addConnectionHandler.bind(websocketService),
    connect: websocketService.connect.bind(websocketService),
    disconnect: websocketService.disconnect.bind(websocketService),
    getConnectionInfo: websocketService.getConnectionInfo.bind(websocketService),
    hasValidToken: websocketService.hasValidToken.bind(websocketService),
  };

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
};
