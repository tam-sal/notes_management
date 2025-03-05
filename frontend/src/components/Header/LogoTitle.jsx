import logo from "./../../assets/logo_aux.png"
import { Link } from 'react-router-dom';

const LogoTitle = () => {
  return (
    <div className="flex items-center space-x-4 pl-8">
      <Link
        to="/"
        className="flex items-center space-x-2 hover:opacity-80 transition"
      >
        <img
          src={logo}
          alt="Logo"
          className="w-10 h-12 rounded-md"
        />
        <h1 className="text-2xl font-bold">
          Note IT
        </h1>
      </Link>
    </div>
  );
};

export default LogoTitle;