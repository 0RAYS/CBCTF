import { useEffect } from 'react';
import { HashRouter as Router } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import AppRoutes from './routes';
import { fetchUserInfo, fetchAccessibleRoutes, setInitialized } from './store/user';
import { WebSocketStatus } from './components/common/websocket/WebSocketStatus.jsx';
import { WebSocketNotice } from './components/common/websocket/WebSocketNotice.jsx';

function App() {
  const dispatch = useDispatch();
  const token = localStorage.getItem('token');

  useEffect(() => {
    const initializeAuth = async () => {
      if (token) {
        await Promise.all([dispatch(fetchUserInfo()), dispatch(fetchAccessibleRoutes())]);
      }
      dispatch(setInitialized());
    };
    initializeAuth();
  }, [dispatch, token]);

  return (
    <div className="relative h-screen w-screen">
      <Router>
        <AppRoutes />
      </Router>
      <WebSocketNotice />
      <WebSocketStatus />
    </div>
  );
}

export default App;
