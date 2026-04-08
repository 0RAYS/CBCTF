import { useEffect } from 'react';
import { HashRouter as Router } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import AppRoutes from './routes';
import { fetchUserInfo, fetchAccessibleRoutes, setInitialized } from './store/user';
import { fetchBranding } from './store/branding';
import BrandingHead from './components/features/Branding/BrandingHead';

function App() {
  const dispatch = useDispatch();
  const token = localStorage.getItem('token');

  useEffect(() => {
    const initializeAuth = async () => {
      const tasks = [dispatch(fetchBranding())];
      if (token) {
        tasks.push(dispatch(fetchUserInfo()), dispatch(fetchAccessibleRoutes()));
      }
      await Promise.all(tasks);
      dispatch(setInitialized());
    };
    initializeAuth();
  }, [dispatch, token]);

  return (
    <div className="relative h-screen w-screen">
      <Router>
        <BrandingHead />
        <AppRoutes />
      </Router>
    </div>
  );
}

export default App;
