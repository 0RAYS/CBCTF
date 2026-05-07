import { Suspense } from 'react';
import Loading from '../components/common/Loading';

export const withSuspense = (Component) => (
  <Suspense fallback={<Loading />}>
    <Component />
  </Suspense>
);

export const withGuard = (element, Guard, guardProps = {}) => <Guard {...guardProps}>{element}</Guard>;
