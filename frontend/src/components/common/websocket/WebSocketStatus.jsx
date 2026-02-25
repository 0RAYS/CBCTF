import { useEffect, useRef } from 'react';
import { useWebSocket } from '../WebSocketProvider.jsx';
import { IconWifiOff } from '@tabler/icons-react';

const HEARTBEAT_INTERVAL = 30_000;

export const WebSocketStatus = () => {
  const { send, isConnected, connectionInfo } = useWebSocket();

  // Keep a ref to the latest send so the interval closure never becomes stale,
  // even though WebSocketProvider re-renders every second (connectionInfo poll).
  const sendRef = useRef(send);
  useEffect(() => {
    sendRef.current = send;
  });

  useEffect(() => {
    if (!isConnected) return;
    const id = setInterval(() => {
      sendRef.current({ type: 'heartbeat', msg: 'ping' });
    }, HEARTBEAT_INTERVAL);
    return () => clearInterval(id);
  }, [isConnected]); // restart only when connection state changes

  if (!connectionInfo.hasToken) {
    return null;
  }

  if (!isConnected) {
    return (
      <div className="fixed bottom-4 left-4 z-50 flex items-center gap-2 bg-red-500/90 text-white px-3 py-2 rounded-md shadow-lg">
        <IconWifiOff size={16} />
        <span className="text-sm font-mono">WebSocket 断开连接</span>
      </div>
    );
  }

  return null;
};
