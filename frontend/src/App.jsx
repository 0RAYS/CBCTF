import { useEffect } from 'react';
import { HashRouter as Router } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import AppRoutes from '@/routes';
import { fetchAccessibleRoutes, fetchUserInfo, setInitialized } from '@/store/user';
import { fetchBranding } from '@/store/branding';
import { fetchPublicConfig } from '@/store/settings.js';
import BrandingHead from '@/components/features/Branding/BrandingHead';

function App() {
  const dispatch = useDispatch();

  useEffect(() => {
    const initializeAuth = async () => {
      await Promise.all([
        dispatch(fetchBranding()),
        dispatch(fetchPublicConfig()),
        dispatch(fetchUserInfo()),
        dispatch(fetchAccessibleRoutes()),
      ]);
      dispatch(setInitialized());
    };
    initializeAuth();
  }, [dispatch]);

  return (
    <div className="relative h-screen w-screen">
      {/* useTransitions={false}: the router updates the location synchronously instead of
          inside startTransition, so a navigation commits immediately and the route's
          Suspense fallback is shown while the new page loads. */}
      <Router useTransitions={false}>
        <BrandingHead />
        <AppRoutes />
      </Router>
    </div>
  );
}

export default App;
