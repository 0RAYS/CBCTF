import { API_CONFIG } from './config.js';

class WebSocketService {
  constructor() {
    this.ws = null;
    this.baseReconnectDelay = 1000; // 1s
    this.maxReconnectDelay = 30000; // 30s
    this.reconnectDelay = this.baseReconnectDelay;
    this.reconnectAttempts = 0;
    this.isConnecting = false;
    this.isManualClose = false;
    this.messageHandlers = new Set();
    this.connectionHandlers = new Set();
    this.reconnectTimer = null;
  }

  // 检查是否有有效的token
  hasValidToken() {
    const token = localStorage.getItem('token');
    return token && token.trim() !== '';
  }

  getTokenValue() {
    const token = localStorage.getItem('token');
    if (!token) return '';
    if (token.includes(' ')) {
      const parts = token.split(' ');
      return parts[1] || '';
    }
    return token;
  }

  // 连接WebSocket
  connect() {
    // 检查是否有token，如果没有则不连接
    if (!this.hasValidToken()) {
      return;
    }

    if (this.isConnecting || this.ws?.readyState === WebSocket.OPEN || this.ws?.readyState === WebSocket.CONNECTING) {
      return;
    }

    this.isConnecting = true;
    this.isManualClose = false;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    try {
      const tokenValue = this.getTokenValue();
      if (!tokenValue) {
        this.isConnecting = false;
        return;
      }

      let wsUrl;
      // 获取当前域名和协议
      if (API_CONFIG.BASE_URL.startsWith('http')) {
        const protocol = API_CONFIG.BASE_URL.startsWith('https') ? 'wss:' : 'ws:';
        const host = API_CONFIG.BASE_URL.replace(/^https?:\/\//, '').replace(/\/$/, '');
        wsUrl = `${protocol}//${host}`;
      } else {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        wsUrl = `${protocol}//${window.location.host}`;
      }
      const magic = localStorage.getItem('LXM') || '';
      // JWT 含有 '.' 字符，不是合法的 Sec-WebSocket-Protocol token。
      // 将整个 token 做 base64url 编码后传输，后端再解码。
      const encodedToken = btoa(tokenValue).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
      const encodedMagic = magic ? btoa(magic).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '') : '';
      const protocols = encodedMagic ? ['token-' + encodedToken, 'magic-' + encodedMagic] : ['token-' + encodedToken];
      this.ws = new WebSocket(wsUrl + '/ws', protocols);

      this.ws.onopen = () => {
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        this.reconnectDelay = this.baseReconnectDelay;
        this.notifyConnectionChange(true);
      };

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          this.notifyMessageHandlers(data);
          // eslint-disable-next-line no-unused-vars
        } catch (error) {
          /* empty */
        }
      };

      this.ws.onclose = () => {
        this.isConnecting = false;
        this.notifyConnectionChange(false);

        if (!this.isManualClose && this.hasValidToken()) {
          this.scheduleReconnect();
        }
      };

      this.ws.onerror = () => {
        this.isConnecting = false;
        if (!this.isManualClose && this.hasValidToken()) {
          this.scheduleReconnect();
        }
        if (this.ws && this.ws.readyState !== WebSocket.CLOSING && this.ws.readyState !== WebSocket.CLOSED) {
          this.ws.close();
        }
      };
      // eslint-disable-next-line no-unused-vars
    } catch (error) {
      this.isConnecting = false;
      if (this.hasValidToken()) {
        this.scheduleReconnect();
      }
    }
  }

  scheduleReconnect() {
    if (!this.hasValidToken()) {
      return;
    }
    if (this.reconnectTimer) {
      return;
    }
    const jitter = Math.floor(Math.random() * 500);
    const delay = Math.min(this.reconnectDelay, this.maxReconnectDelay) + jitter;
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      if (!this.isManualClose && this.hasValidToken()) {
        this.connect();
      }
    }, delay);
    this.reconnectAttempts += 1;
    this.reconnectDelay = Math.min(this.baseReconnectDelay * 2 ** this.reconnectAttempts, this.maxReconnectDelay);
  }

  send(message) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    }
  }

  disconnect() {
    this.isManualClose = true;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      this.ws.onopen = null;
      this.ws.onclose = null;
      this.ws.onerror = null;
      this.ws.onmessage = null;
      this.ws.close();
      this.ws = null;
    }
  }

  isConnected() {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  addMessageHandler(handler) {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  addConnectionHandler(handler) {
    this.connectionHandlers.add(handler);
    return () => this.connectionHandlers.delete(handler);
  }

  notifyMessageHandlers(data) {
    this.messageHandlers.forEach((handler) => {
      handler(data);
    });
  }

  notifyConnectionChange(isConnected) {
    this.connectionHandlers.forEach((handler) => {
      handler(isConnected);
    });
  }

  getConnectionInfo() {
    if (!this.ws) {
      return {
        connected: false,
        readyState: 'CLOSED',
        hasToken: this.hasValidToken(),
        reconnectAttempts: this.reconnectAttempts,
        reconnectDelay: this.reconnectDelay,
      };
    }

    const states = {
      [WebSocket.CONNECTING]: 'CONNECTING',
      [WebSocket.OPEN]: 'OPEN',
      [WebSocket.CLOSING]: 'CLOSING',
      [WebSocket.CLOSED]: 'CLOSED',
    };

    return {
      connected: this.ws.readyState === WebSocket.OPEN,
      readyState: states[this.ws.readyState],
      hasToken: this.hasValidToken(),
      reconnectAttempts: this.reconnectAttempts,
      reconnectDelay: this.reconnectDelay,
    };
  }

  handleTokenChange() {
    if (this.hasValidToken()) {
      // 如果有token且当前未连接，则尝试连接
      if (!this.isConnected() && !this.isConnecting) {
        this.connect();
      }
    } else {
      // 如果没有token，则断开连接
      if (this.isConnected() || this.isConnecting) {
        this.disconnect();
      }
    }
  }
}

const websocketService = new WebSocketService();

window.addEventListener('storage', (event) => {
  if (event.key === 'token') {
    websocketService.handleTokenChange();
  }
});

export default websocketService;
