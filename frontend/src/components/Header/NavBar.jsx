import { useNavigate } from 'react-router-dom';
import axiosInstance from '../../config/axios';
import { useAuth } from '../../context/AuthContext';
import showToast from './../../context/toast-utils';

const NavBar = ({ setIsMenuOpen }) => {
  const { user, checkAuth } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      const response = await axiosInstance.post('/user/logout');
      if (response.data?.status === 200) {
        showToast(true, response.data.data.message);
        await checkAuth();
        setIsMenuOpen(false);
        navigate('/');
      }
    } catch (err) {
      showToast(false, err.response?.data?.error || 'Logout failed');
    }
  };

  const handleNavigation = (path) => {
    setIsMenuOpen(false);
    navigate(path);
  };

  return (
    <div className="flex flex-col lg:flex-row space-y-4 lg:space-y-0 lg:space-x-4 text-white">
      <button onClick={() => handleNavigation('/')} className="hover:text-primary transition">
        Home
      </button>

      {!user ? (
        <>
          <button onClick={() => handleNavigation('/register')} className="hover:text-primary transition">
            Register
          </button>
          <button onClick={() => handleNavigation('/login')} className="hover:text-primary transition">
            Login
          </button>
        </>
      ) : (
        <>
          <button onClick={() => handleNavigation('/create-note')} className="hover:text-primary transition">
            Create
          </button>
          <button onClick={() => handleNavigation('/notes')} className="hover:text-primary transition">
            Notes
          </button>
          <button onClick={handleLogout} className="hover:text-primary transition">
            Logout
          </button>
        </>
      )}
    </div>
  );
};

export default NavBar;