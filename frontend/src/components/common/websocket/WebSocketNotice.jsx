import { useWebSocket } from '../WebSocketProvider.jsx';
import { toast } from '../../../utils/toast.js';
import { useEffect } from 'react';

export const WebSocketNotice = () => {
  const { addMessageHandler } = useWebSocket();

  useEffect(() => {
    return addMessageHandler((data) => {
      if (data.type === 'contest_notice') {
        switch (data.level) {
          case 'error':
            toast.danger({ title: data.title, description: data.msg });
            break;
          case 'warning':
            toast.warning({ title: data.title, description: data.msg });
            break;
          case 'success':
            toast.success({ title: data.title, description: data.msg });
            break;
          case 'info':
            toast.info({ title: data.title, description: data.msg });
            break;
          case 'notice':
            toast.notice({ title: data.title, description: data.msg });
            break;
          default:
            toast.default({ title: data.title, description: data.msg });
            break;
        }
      }
    });
  }, [addMessageHandler]);
};
