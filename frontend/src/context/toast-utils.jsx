import toast from 'react-hot-toast';

const showToast = (isSuccessful, message) => {
  if (isSuccessful) {
    toast.success(message, {
      duration: 3000,
      position: 'top-right',
      style: {
        borderRadius: '10px',
        background: '#4CAF50',
        color: '#fff',
      },
    });
  } else {
    toast.error(message, {
      duration: 3000,
      position: 'top-right',
      style: {
        borderRadius: '10px',
        background: '#f44336',
        color: '#fff',
      },
    });
  }
};


export default showToast;