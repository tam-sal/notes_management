import React, { createContext, useState, useEffect, useContext } from 'react';
import axiosInstance from '../config/axios';

const AuthContext = createContext();

export const useAuth = () => {
  return useContext(AuthContext);
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const checkAuth = async () => {
    try {
      const response = await axiosInstance.get('/user/auth-check');
      const { authenticated, user_id } = response.data.data;

      if (authenticated) {
        setUser({ id: user_id });
      } else {
        console.error('User is not authenticated');
        setUser(null);
      }
    } catch (err) {
      const errorMsg = err.response?.data?.error
        ? Object.values(err.response?.data?.error).join(' ')
        : 'An error occurred during authentication.';
      setError(errorMsg);
      console.error('Error during auth check:', errorMsg);
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);
  useEffect(() => {
  }, [user]);

  return (
    <AuthContext.Provider value={{ user, loading, error, checkAuth }}>
      {children}
    </AuthContext.Provider>
  );
};