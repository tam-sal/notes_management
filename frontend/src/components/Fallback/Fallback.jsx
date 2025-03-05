// components/ErrorPage/ErrorPage.jsx
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const Fallback = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const timeout = setTimeout(() => navigate('/'), 5000);
    return () => clearTimeout(timeout);
  }, [navigate]);

  return (
    <div className="text-center">
      <h1 className="text-4xl font-bold mb-4">404 - Page Not Found</h1>
      <p className="text-lg">Redirecting to homepage in 5 seconds...</p>
      <button
        onClick={() => navigate('/')}
        className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        Go Home Now
      </button>
    </div>
  );
};

export default Fallback;