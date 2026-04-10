import { useEffect } from 'react';
import { HashRouter as Router } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import AppRoutes from './routes';
import { fetchUserInfo, fetchAccessibleRoutes, setInitialized } from './store/user';
import { fetchBranding } from './store/branding';
import BrandingHead from './components/features/Branding/BrandingHead';

function App() {
  const dispatch = useDispatch();

  useEffect(() => {
    const initializeAuth = async () => {
      await Promise.all([dispatch(fetchBranding()), dispatch(fetchUserInfo()), dispatch(fetchAccessibleRoutes())]);
      dispatch(setInitialized());
    };
    initializeAuth();
  }, [dispatch]);

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
