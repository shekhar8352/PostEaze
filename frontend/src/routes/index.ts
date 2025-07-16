import { useRoutes } from 'react-router-dom';
import { authRoutes } from '@/features/auth';

const AppRoutes = () => {
  const routes = [...authRoutes]; // Later add dashboard, posts, etc.
  return useRoutes(routes);
};

export default AppRoutes;