import { useWebSocket } from '../WebSocketProvider.jsx';
import { IconWifiOff } from '@tabler/icons-react';

export const WebSocketStatus = () => {
  const { send, isConnected, connectionInfo } = useWebSocket();

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
  send({ type: 'heartbeat', msg: 'ping' });
};
